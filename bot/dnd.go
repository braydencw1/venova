package bot

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/braydencw1/venova/db"
)

type Roll struct {
	NumDice      int
	DieSize      int
	Modifier     int
	Advantage    bool
	Disadvantage bool
}

type RollResult struct {
	Rolls        [][]int
	Totals       []int
	FinalTotal   int
	Modifier     int
	Advantage    bool
	Disadvantage bool
}

var rollPattern = regexp.MustCompile(`^(?i)(\d*)d\d+([+-]\d+)?$`)

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
		return fmt.Errorf("error parsing date: %w", err)
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
	now := time.Now()
	currRoleId := getMemberDNDRole(msg.Member)
	if currRoleId == "" {
		return ctx.Reply("Could not find DND role")
	}
	dateOfPlay, _, err := db.GetLatestPlayDate(currRoleId)
	if err != nil {
		return ctx.Reply("Could not find play date information. Perhaps wrong server.")
	}
	fmtDate := fmt.Sprint(dateOfPlay.Format("01-02-2006"))
	if dateOfPlay.Before(now) {
		if err := ctx.Reply(fmt.Sprintf("There is no date currently set. Your last session was: %s", fmtDate)); err != nil {
			return err
		}
	}
	if err = ctx.Reply(fmtDate); err != nil {
		return err
	}
	return nil
}

func rollCmd(ctx CommandCtx) error {
	advantage := false
	disadvantage := false

	if len(ctx.Args) > 1 {

		swt := ctx.Args[1]
		switch strings.ToLower(swt) {
		case "a", "adv", "advantage":
			advantage = true
		case "d", "dis", "disadvantage":
			disadvantage = true
		}
	}

	rollArgs := []string{}
	for _, arg := range ctx.Args {
		lower := strings.ToLower(arg)
		if lower == "a" || lower == "adv" || lower == "advantage" ||
			lower == "d" || lower == "dis" || lower == "disadvantage" {
			continue
		}
		rollArgs = append(rollArgs, arg)
	}

	roll := strings.Join(rollArgs, "")
	roll = strings.ReplaceAll(roll, " ", "")

	if !strings.Contains(strings.ToLower(roll), "d") {
		return ctx.Reply("Must be NdM±K (e.g., 1d20+3 or d20).")
	}

	if !isValidRoll(roll) {
		return ctx.Reply("Invalid format. Use NdM±K (e.g., 1d20+3 or d20).")
	}

	roll = strings.ToLower(roll)

	parts := strings.SplitN(roll, "d", 2)

	if len(parts) != 2 || parts[1] == "" {
		return ctx.Reply("Invalid format. Use NdM±K (e.g., 1d20+3 or d20).")
	}

	numDiceStr := parts[0]

	if numDiceStr == "" {
		numDiceStr = "1"
	}

	numDice, err := strconv.Atoi(numDiceStr)
	if err != nil {
		return err
	}

	if numDice <= 0 || numDice > 100 {
		return ctx.Reply("Invalid number of dice (must be 1–100).")
	}

	modifier := 0
	diePart := parts[1]
	if diePart == "" {
		return ctx.Reply("Missing die size.")
	}

	var dieSizeStr, modStr string
	if strings.Contains(diePart, "+") {
		split := strings.SplitN(diePart, "+", 2)
		dieSizeStr, modStr = split[0], split[1]
		modifier, err = strconv.Atoi(modStr)
		if err != nil {
			return ctx.Reply("Invalid modifier; must be a number.")
		}
	} else if strings.Contains(diePart, "-") {
		split := strings.SplitN(diePart, "-", 2)
		dieSizeStr, modStr = split[0], split[1]
		modifier, err = strconv.Atoi(modStr)
		if err != nil {
			return ctx.Reply("Invalid modifier; must be a number.")
		}
		modifier = -modifier
	} else {
		dieSizeStr = diePart
	}

	dieSize, err := strconv.Atoi(dieSizeStr)
	if err != nil {
		return ctx.Reply("Invalid die size.")
	}

	if dieSize <= 0 || dieSize > 1000 {
		return ctx.Reply("Invalid die size (must be 1–1000).")
	}

	r := Roll{
		NumDice:      numDice,
		DieSize:      dieSize,
		Modifier:     modifier,
		Advantage:    advantage,
		Disadvantage: disadvantage,
	}
	rollResult := r.Execute()

	msg := rollResult.FormatMessage()

	return ctx.Reply(msg)
}

func isValidRoll(arg string) bool {
	return rollPattern.MatchString(arg)
}

func (r Roll) Execute() RollResult {
	doRoll := func() ([]int, int) {
		rolls := make([]int, r.NumDice)
		total := 0
		for i := 0; i < r.NumDice; i++ {
			roll := rand.Intn(r.DieSize) + 1
			rolls[i] = roll
			total += roll
		}
		total += r.Modifier
		return rolls, total
	}

	rolls1, total1 := doRoll()
	result := RollResult{
		Rolls:        [][]int{rolls1},
		Totals:       []int{total1},
		Modifier:     r.Modifier,
		Advantage:    r.Advantage,
		Disadvantage: r.Disadvantage,
		FinalTotal:   total1,
	}

	if r.Advantage || r.Disadvantage {
		rolls2, total2 := doRoll()
		result.Rolls = append(result.Rolls, rolls2)
		result.Totals = append(result.Totals, total2)

		if r.Advantage && total2 > total1 {
			result.FinalTotal = total2
		} else if r.Disadvantage && total2 < total1 {
			result.FinalTotal = total2
		}
	}

	return result
}

func (r RollResult) FormatMessage() string {
	if !r.Advantage && !r.Disadvantage {
		return fmt.Sprintf(
			"🎲 Roll: (%v) %+d = %d\n→ Result: **%d**",
			r.Rolls[0], r.Modifier, r.FinalTotal, r.FinalTotal,
		)
	}

	return fmt.Sprintf(
		":game_die: Rolls: (%v) %+d = %d | (%v) %+d = %d\n→ Result: **%d**",
		r.Rolls[0], r.Modifier, r.Totals[0],
		r.Rolls[1], r.Modifier, r.Totals[1],
		r.FinalTotal,
	)
}
