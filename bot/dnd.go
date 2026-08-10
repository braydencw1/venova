package bot

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/braydencw1/venova/db"
)

type DiceGroup struct {
	NumDice int
	DieSize int
}

type Roll struct {
	Groups       []DiceGroup
	Modifier     int
	Advantage    bool
	Disadvantage bool
}

type GroupResult struct {
	Group DiceGroup
	Rolls []int
}

type AttemptResult struct {
	Groups []GroupResult
	Total  int
}

type RollResult struct {
	Attempts     []AttemptResult
	ChosenIdx    int
	Modifier     int
	Advantage    bool
	Disadvantage bool
}

func playDndCmd(ctx CommandCtx) error {
	msg := ctx.Message
	args := ctx.Args
	// Set dnd play date
	if !ctx.IDChecker.IsAdmin(msg.Author.ID) {
		return nil
	}
	layout := "01-02-2006"
	t, err := time.Parse(layout, args[0])
	if err != nil {
		return ctx.Reply("Invalid date. Use MM-DD-YYYY (e.g., 08-15-2026).")
	}
	currRoleId := getMemberDNDRole(msg.Member)
	if currRoleId == "" {
		return ctx.Reply("Your dnd role is not found in the db.")
	}
	if err := db.InsertPlayDate(t, currRoleId); err != nil {
		return fmt.Errorf("error inserting into table %w", err)
	}
	if err = ctx.Reply("The Date has been updated."); err != nil {
		return err
	}
	return nil
}

func whenIsDndCmd(ctx CommandCtx) error {
	msg := ctx.Message.Message
	// Play dates are stored at midnight, so compare against the start of
	// today or the bot claims there is no session on game day itself.
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	currRoleId := getMemberDNDRole(msg.Member)
	if currRoleId == "" {
		return ctx.Reply("Could not find DND role")
	}

	dateOfPlay, _, err := db.GetLatestPlayDate(currRoleId)
	if err != nil {
		return ctx.Reply("Could not find play date information. Perhaps wrong server.")
	}

	fmtDate := fmt.Sprint(dateOfPlay.Format("01-02-2006"))
	if dateOfPlay.Before(today) {
		return ctx.Reply(fmt.Sprintf("There is no date currently set. Your last session was: %s", fmtDate))
	}

	if err = ctx.Reply(fmtDate); err != nil {
		return err
	}
	return nil
}

func rollCmd(ctx CommandCtx) error {
	if strings.EqualFold(ctx.Args[0], "stats") {
		return ctx.Reply(RollStats())
	}

	roll, err := ParseRoll(ctx.Args)
	if err != nil {
		return ctx.Reply(upperFirst(err.Error()))
	}

	return ctx.Reply(roll.Execute().FormatMessage())
}

// upperFirst capitalizes lint-friendly lowercase error strings for display as
// Discord replies.
func upperFirst(s string) string {
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return s
	}
	return string(unicode.ToUpper(r)) + s[size:]
}

var (
	diceTermPattern = regexp.MustCompile(`^(\d*)d(\d+)$`)
	errInvalidRoll  = errors.New("invalid format: use dice like 2d6+1d8+3 or d20, optionally with adv/dis")
)

// ParseRoll parses full !roll arguments: a dice expression, possibly split
// across arguments, with advantage/disadvantage flags anywhere among them.
func ParseRoll(args []string) (Roll, error) {
	advantage := false
	disadvantage := false

	rollArgs := []string{}
	for _, arg := range args {
		switch strings.ToLower(arg) {
		case "a", "adv", "advantage":
			advantage = true
		case "d", "dis", "disadvantage":
			disadvantage = true
		default:
			rollArgs = append(rollArgs, arg)
		}
	}

	if advantage && disadvantage {
		return Roll{}, errors.New("pick either advantage or disadvantage, not both")
	}

	roll, err := ParseRollExpression(strings.Join(rollArgs, ""))
	if err != nil {
		return Roll{}, err
	}
	roll.Advantage = advantage
	roll.Disadvantage = disadvantage
	return roll, nil
}

// ParseRollExpression parses a dice expression like "2d6+1d8+5-2" into dice
// groups and a flat modifier.
func ParseRollExpression(expr string) (Roll, error) {
	expr = strings.ToLower(strings.ReplaceAll(expr, " ", ""))
	if expr == "" {
		return Roll{}, errors.New("roll something, e.g. !roll 2d6+3 or !roll d20 adv")
	}

	// Normalize subtraction into "+-" so the expression splits into signed terms.
	expr = strings.ReplaceAll(expr, "-", "+-")
	terms := strings.Split(expr, "+")

	var roll Roll
	totalDice := 0
	for i, term := range terms {
		if term == "" {
			// A leading sign produces one empty first term; anything else
			// (e.g. "2d6++3") is malformed.
			if i == 0 {
				continue
			}
			return Roll{}, errInvalidRoll
		}

		if m := diceTermPattern.FindStringSubmatch(term); m != nil {
			numDice := 1
			if m[1] != "" {
				n, err := strconv.Atoi(m[1])
				if err != nil || n <= 0 || n > 100 {
					return Roll{}, errors.New("invalid number of dice (must be 1–100)")
				}
				numDice = n
			}
			totalDice += numDice
			if totalDice > 100 {
				return Roll{}, errors.New("too many dice (100 max across the whole roll)")
			}

			dieSize, err := strconv.Atoi(m[2])
			if err != nil || dieSize <= 0 || dieSize > 1000 {
				return Roll{}, errors.New("invalid die size (must be 1–1000)")
			}

			roll.Groups = append(roll.Groups, DiceGroup{NumDice: numDice, DieSize: dieSize})
			continue
		}

		if after, ok := strings.CutPrefix(term, "-"); ok && diceTermPattern.MatchString(after) {
			return Roll{}, errors.New("cannot subtract dice")
		}

		mod, err := strconv.Atoi(term)
		if err != nil {
			return Roll{}, errInvalidRoll
		}
		roll.Modifier += mod
	}

	if len(roll.Groups) == 0 {
		return Roll{}, errInvalidRoll
	}
	return roll, nil
}

func (r Roll) Execute() RollResult {
	attempt := func() AttemptResult {
		total := r.Modifier
		groups := make([]GroupResult, 0, len(r.Groups))
		for _, g := range r.Groups {
			rolls := make([]int, g.NumDice)
			for i := range rolls {
				rolls[i] = rand.Intn(g.DieSize) + 1
				total += rolls[i]
			}
			groups = append(groups, GroupResult{Group: g, Rolls: rolls})
		}
		return AttemptResult{Groups: groups, Total: total}
	}

	result := RollResult{
		Attempts:     []AttemptResult{attempt()},
		Modifier:     r.Modifier,
		Advantage:    r.Advantage,
		Disadvantage: r.Disadvantage,
	}

	if r.Advantage || r.Disadvantage {
		result.Attempts = append(result.Attempts, attempt())
		second := result.Attempts[1].Total
		first := result.Attempts[0].Total
		if (r.Advantage && second > first) || (r.Disadvantage && second < first) {
			result.ChosenIdx = 1
		}
	}

	return result
}

func (r RollResult) FinalTotal() int {
	return r.Attempts[r.ChosenIdx].Total
}

func (a AttemptResult) format(modifier int) string {
	parts := make([]string, 0, len(a.Groups))
	for _, g := range a.Groups {
		nums := make([]string, len(g.Rolls))
		for i, roll := range g.Rolls {
			nums[i] = strconv.Itoa(roll)
		}
		parts = append(parts, fmt.Sprintf("%dd%d [%s]", g.Group.NumDice, g.Group.DieSize, strings.Join(nums, " ")))
	}
	s := strings.Join(parts, " + ")
	if modifier != 0 {
		s += fmt.Sprintf(" %+d", modifier)
	}
	return fmt.Sprintf("%s = %d", s, a.Total)
}

func (r RollResult) FormatMessage() string {
	chosen := r.Attempts[r.ChosenIdx]

	var b strings.Builder
	if !r.Advantage && !r.Disadvantage {
		fmt.Fprintf(&b, "🎲 Roll: %s\n→ Result: **%d**", chosen.format(r.Modifier), chosen.Total)
	} else {
		label := "advantage"
		if r.Disadvantage {
			label = "disadvantage"
		}
		fmt.Fprintf(&b, ":game_die: Rolls: %s | %s\n→ Result: **%d** (%s)",
			r.Attempts[0].format(r.Modifier), r.Attempts[1].format(r.Modifier), chosen.Total, label)
	}

	nat20, nat1 := false, false
	for _, g := range chosen.Groups {
		if g.Group.DieSize != 20 {
			continue
		}
		for _, roll := range g.Rolls {
			nat20 = nat20 || roll == 20
			nat1 = nat1 || roll == 1
		}
	}
	if nat20 {
		b.WriteString("\n💥 **Natural 20!**")
	}
	if nat1 {
		b.WriteString("\n💀 **Natural 1!**")
	}

	return b.String()
}

// RollStats rolls ability scores the classic way: 4d6 drop the lowest, six times.
func RollStats() string {
	var b strings.Builder
	b.WriteString("🎲 Ability scores (4d6, drop lowest):\n")

	totals := make([]int, 0, 6)
	for range 6 {
		rolls := make([]int, 4)
		lowest := 0
		sum := 0
		for i := range rolls {
			rolls[i] = rand.Intn(6) + 1
			sum += rolls[i]
			if rolls[i] < rolls[lowest] {
				lowest = i
			}
		}
		total := sum - rolls[lowest]
		totals = append(totals, total)

		parts := make([]string, len(rolls))
		for i, roll := range rolls {
			if i == lowest {
				parts[i] = fmt.Sprintf("~~%d~~", roll)
			} else {
				parts[i] = strconv.Itoa(roll)
			}
		}
		fmt.Fprintf(&b, "[%s] = **%d**\n", strings.Join(parts, " "), total)
	}

	slices.SortFunc(totals, func(a, b int) int { return b - a })
	sum := 0
	for _, t := range totals {
		sum += t
	}
	strTotals := make([]string, len(totals))
	for i, t := range totals {
		strTotals[i] = strconv.Itoa(t)
	}
	fmt.Fprintf(&b, "→ Scores: %s (sum %d)", strings.Join(strTotals, ", "), sum)
	return b.String()
}
