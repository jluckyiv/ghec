package ghec_test

import (
	"testing"

	"github.com/jluckyiv/ghec"
)

type testCase struct {
	name     string
	base     ghec.BaseEnhancement
	targets  int
	level    ghec.Level
	prev     ghec.PreviousEnhancements
	expected ghec.Cost
}

var testCases []testCase = []testCase{
	{
		name:     "example 1 from README",
		base:     ghec.EnhanceAttack,
		targets:  3,
		level:    ghec.Level3,
		prev:     ghec.PreviousEnhancements0,
		expected: ghec.Cost(150),
	},
	{
		name:     "example 2 from README",
		base:     ghec.EnhanceAddAttackHex,
		targets:  3,
		level:    ghec.Level3,
		prev:     ghec.PreviousEnhancements1,
		expected: ghec.Cost(191),
	},
	{
		name:     "add range, level 1, previous 1",
		base:     ghec.EnhanceRange,
		targets:  2,
		level:    ghec.Level1,
		prev:     ghec.PreviousEnhancements1,
		expected: ghec.Cost(135),
	},
	{
		name:     "add target, level 2, previous 0",
		base:     ghec.EnhanceTarget,
		targets:  3,
		level:    ghec.Level2,
		prev:     ghec.PreviousEnhancements0,
		expected: ghec.Cost(125),
	},
	{
		name:     "add summons move, level 1, previous 1",
		base:     ghec.EnhanceSummonsMove,
		targets:  1,
		level:    ghec.Level1,
		prev:     ghec.PreviousEnhancements1,
		expected: ghec.Cost(175),
	},
	{
		name:     "add summons range, level 1, previous 1",
		base:     ghec.EnhanceSummonsRange,
		targets:  1,
		level:    ghec.Level1,
		prev:     ghec.PreviousEnhancements1,
		expected: ghec.Cost(125),
	},
}

func TestFunctionalOptionsAPI(t *testing.T) {
	for _, tc := range testCases {
		input := ghec.NewEnhancement(tc.base,
			ghec.OptionWithLevel(tc.level),
			ghec.OptionWithMultipleTarget(tc.targets),
			ghec.OptionWithPreviousEnhancements(tc.prev),
		)

		actual, _ := input.Cost()
		if actual != tc.expected {
			t.Fatalf("%s; expected %d, got %d", tc.name, tc.expected, actual)
		}
	}
}
