package bot

import (
	"strings"
	"testing"
)

func TestParseRollFlags(t *testing.T) {
	cases := []struct {
		name string
		args []string
		adv  bool
		dis  bool
	}{
		{"no flags", []string{"d20"}, false, false},
		{"adv after dice", []string{"d20", "adv"}, true, false},
		{"adv before dice", []string{"adv", "d20"}, true, false},
		{"short a", []string{"a", "d20"}, true, false},
		{"long advantage", []string{"d20", "advantage"}, true, false},
		{"dis after dice", []string{"d20", "dis"}, false, true},
		{"short d", []string{"d20", "d"}, false, true},
		{"long disadvantage", []string{"disadvantage", "d20"}, false, true},
		{"mixed case", []string{"d20", "ADV"}, true, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			roll, err := ParseRoll(c.args)
			if err != nil {
				t.Fatalf("ParseRoll(%v) error: %v", c.args, err)
			}
			if roll.Advantage != c.adv || roll.Disadvantage != c.dis {
				t.Errorf("ParseRoll(%v) adv=%v dis=%v, want adv=%v dis=%v",
					c.args, roll.Advantage, roll.Disadvantage, c.adv, c.dis)
			}
		})
	}
}

func TestParseRollErrors(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{"both flags", []string{"d20", "adv", "dis"}},
		{"no dice", []string{"adv"}},
		{"garbage", []string{"abc"}},
		{"modifier only", []string{"5"}},
		{"missing die size", []string{"2d"}},
		{"double plus", []string{"2d6++3"}},
		{"too many dice", []string{"101d6"}},
		{"too many dice across groups", []string{"60d6+41d6"}},
		{"zero sided die", []string{"2d0"}},
		{"die too large", []string{"d1001"}},
		{"dice count overflow", []string{"99999999999999999999d6"}},
		{"subtracting dice", []string{"2d6-1d4"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := ParseRoll(c.args); err == nil {
				t.Errorf("ParseRoll(%v) succeeded, want error", c.args)
			}
		})
	}
}

func TestParseRollExpression(t *testing.T) {
	cases := []struct {
		expr     string
		groups   []DiceGroup
		modifier int
	}{
		{"d20", []DiceGroup{{1, 20}}, 0},
		{"1d20+3", []DiceGroup{{1, 20}}, 3},
		{"2d6-2", []DiceGroup{{2, 6}}, -2},
		{"2d6+1d8+5", []DiceGroup{{2, 6}, {1, 8}}, 5},
		{"2d6+1d8+5-2", []DiceGroup{{2, 6}, {1, 8}}, 3},
		{"D20+1", []DiceGroup{{1, 20}}, 1},
		{"-2+d20", []DiceGroup{{1, 20}}, -2},
		{"100d6", []DiceGroup{{100, 6}}, 0},
		{"d1000", []DiceGroup{{1, 1000}}, 0},
	}
	for _, c := range cases {
		t.Run(c.expr, func(t *testing.T) {
			roll, err := ParseRollExpression(c.expr)
			if err != nil {
				t.Fatalf("ParseRollExpression(%q) error: %v", c.expr, err)
			}
			if roll.Modifier != c.modifier {
				t.Errorf("modifier = %d, want %d", roll.Modifier, c.modifier)
			}
			if len(roll.Groups) != len(c.groups) {
				t.Fatalf("groups = %v, want %v", roll.Groups, c.groups)
			}
			for i, g := range roll.Groups {
				if g != c.groups[i] {
					t.Errorf("group %d = %v, want %v", i, g, c.groups[i])
				}
			}
		})
	}
}

func TestExecuteTotals(t *testing.T) {
	// d1 dice make results deterministic: every die rolls 1.
	roll := Roll{Groups: []DiceGroup{{3, 1}, {2, 1}}, Modifier: 2}
	res := roll.Execute()
	if got := res.FinalTotal(); got != 7 {
		t.Errorf("FinalTotal = %d, want 7", got)
	}
	if len(res.Attempts) != 1 {
		t.Errorf("attempts = %d, want 1", len(res.Attempts))
	}
}

func TestExecuteBounds(t *testing.T) {
	roll := Roll{Groups: []DiceGroup{{2, 6}, {1, 8}}, Modifier: 3}
	for range 100 {
		res := roll.Execute()
		total := res.FinalTotal()
		if total < 3+3 || total > 20+3 {
			t.Fatalf("total %d out of range [6, 23]", total)
		}
	}
}

func TestExecuteAdvantagePicksHigher(t *testing.T) {
	roll := Roll{Groups: []DiceGroup{{1, 20}}, Advantage: true}
	for range 200 {
		res := roll.Execute()
		if len(res.Attempts) != 2 {
			t.Fatalf("attempts = %d, want 2", len(res.Attempts))
		}
		want := max(res.Attempts[0].Total, res.Attempts[1].Total)
		if res.FinalTotal() != want {
			t.Fatalf("advantage picked %d, want %d (attempts %d, %d)",
				res.FinalTotal(), want, res.Attempts[0].Total, res.Attempts[1].Total)
		}
	}
}

func TestExecuteDisadvantagePicksLower(t *testing.T) {
	roll := Roll{Groups: []DiceGroup{{1, 20}}, Disadvantage: true}
	for range 200 {
		res := roll.Execute()
		want := min(res.Attempts[0].Total, res.Attempts[1].Total)
		if res.FinalTotal() != want {
			t.Fatalf("disadvantage picked %d, want %d (attempts %d, %d)",
				res.FinalTotal(), want, res.Attempts[0].Total, res.Attempts[1].Total)
		}
	}
}

func attemptOf(mod int, groups ...GroupResult) AttemptResult {
	total := mod
	for _, g := range groups {
		for _, r := range g.Rolls {
			total += r
		}
	}
	return AttemptResult{Groups: groups, Total: total}
}

func TestFormatMessage(t *testing.T) {
	t.Run("no modifier omits +0", func(t *testing.T) {
		res := RollResult{Attempts: []AttemptResult{
			attemptOf(0, GroupResult{Group: DiceGroup{1, 6}, Rolls: []int{4}}),
		}}
		msg := res.FormatMessage()
		if strings.Contains(msg, "+0") {
			t.Errorf("message contains +0: %q", msg)
		}
		if !strings.Contains(msg, "**4**") {
			t.Errorf("message missing total: %q", msg)
		}
	})

	t.Run("modifier shown", func(t *testing.T) {
		res := RollResult{Modifier: 3, Attempts: []AttemptResult{
			attemptOf(3, GroupResult{Group: DiceGroup{1, 6}, Rolls: []int{4}}),
		}}
		if msg := res.FormatMessage(); !strings.Contains(msg, "+3") {
			t.Errorf("message missing modifier: %q", msg)
		}
	})

	t.Run("natural 20 callout", func(t *testing.T) {
		res := RollResult{Attempts: []AttemptResult{
			attemptOf(0, GroupResult{Group: DiceGroup{1, 20}, Rolls: []int{20}}),
		}}
		if msg := res.FormatMessage(); !strings.Contains(msg, "Natural 20") {
			t.Errorf("message missing nat 20 callout: %q", msg)
		}
	})

	t.Run("natural 1 callout", func(t *testing.T) {
		res := RollResult{Attempts: []AttemptResult{
			attemptOf(0, GroupResult{Group: DiceGroup{1, 20}, Rolls: []int{1}}),
		}}
		if msg := res.FormatMessage(); !strings.Contains(msg, "Natural 1") {
			t.Errorf("message missing nat 1 callout: %q", msg)
		}
	})

	t.Run("no callout for non-d20", func(t *testing.T) {
		res := RollResult{Attempts: []AttemptResult{
			attemptOf(0, GroupResult{Group: DiceGroup{1, 10}, Rolls: []int{10}}),
		}}
		if msg := res.FormatMessage(); strings.Contains(msg, "Natural") {
			t.Errorf("unexpected callout for d10: %q", msg)
		}
	})

	t.Run("callout only from chosen attempt", func(t *testing.T) {
		res := RollResult{
			Advantage: true,
			ChosenIdx: 1,
			Attempts: []AttemptResult{
				attemptOf(0, GroupResult{Group: DiceGroup{1, 20}, Rolls: []int{1}}),
				attemptOf(0, GroupResult{Group: DiceGroup{1, 20}, Rolls: []int{15}}),
			},
		}
		if msg := res.FormatMessage(); strings.Contains(msg, "Natural 1") {
			t.Errorf("callout leaked from discarded attempt: %q", msg)
		}
	})

	t.Run("advantage shows both attempts", func(t *testing.T) {
		res := RollResult{
			Advantage: true,
			ChosenIdx: 0,
			Attempts: []AttemptResult{
				attemptOf(0, GroupResult{Group: DiceGroup{1, 20}, Rolls: []int{18}}),
				attemptOf(0, GroupResult{Group: DiceGroup{1, 20}, Rolls: []int{7}}),
			},
		}
		msg := res.FormatMessage()
		for _, want := range []string{"[18]", "[7]", "**18**", "advantage"} {
			if !strings.Contains(msg, want) {
				t.Errorf("message missing %q: %q", want, msg)
			}
		}
	})
}

func TestRollStats(t *testing.T) {
	msg := RollStats()
	lines := strings.Split(strings.TrimSpace(msg), "\n")
	// Header + six score lines + summary line.
	if len(lines) != 8 {
		t.Fatalf("got %d lines, want 8: %q", len(lines), msg)
	}
	if !strings.Contains(msg, "~~") {
		t.Errorf("no dropped die marked: %q", msg)
	}
	if !strings.Contains(lines[7], "Scores:") {
		t.Errorf("missing summary line: %q", lines[7])
	}
}
