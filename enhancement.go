package ghec

import "fmt"

// enhancement is a struct that holds the information needed to calculate its
// cost. It is not exported to limit the API surface area. Its only methods are
// With* methods to set its fields and Cost to calculate its cost.
type enhancement struct {
	// baseEnhancement is the base enhancement to calculate the cost.
	// Each base enhancement has a fixed cost.
	baseEnhancement BaseEnhancement
	// level is the level of the ability card to calculate the cost.
	// It must be between 1 and 9.
	level Level
	// multipleTarget serves two purposes:
	// 1. It triggers the multiplier for multiple-target enhancements.
	// 2. It sets the number of current hexes for Add Attack Hex enhancements.
	multipleTarget int
	// previousEnhancements is the number of previous enhancements on the ability
	// card. It must be between 0 and 3.
	previousEnhancements PreviousEnhancements
}

// NewEnhancement creates a new enhancement to calculate its cost.
func NewEnhancement(baseEnhancement BaseEnhancement) enhancement {
	return enhancement{
		baseEnhancement:      baseEnhancement,
		level:                Level1,
		multipleTarget:       0,
		previousEnhancements: PreviousEnhancements0,
	}
}

// WithMultipleTarget sets the number of targets for the enhancement.
// It also sets the number of current hexes for Add Attack Hex enhancements.
func (e enhancement) WithMultipleTarget(multipleTarget int) enhancement {
	e.multipleTarget = multipleTarget
	return e
}

// WithLevel sets the level of the ability card for the enhancement.
func (e enhancement) WithLevel(level Level) enhancement {
	e.level = level
	return e
}

// WithPreviousEnhancements sets the number of previous enhancements on the
// card.
func (e enhancement) WithPreviousEnhancements(previousEnhancements PreviousEnhancements) enhancement {
	e.previousEnhancements = previousEnhancements
	return e
}

// Cost calculates the cost of the enhancement.
// It returns an error if the level or previous enhancements are out of bounds,
// since the With* methods do not validate their inputs.
func (e enhancement) Cost() (Cost, error) {
	if e.level < 1 || e.level > 9 {
		return 0, fmt.Errorf("level must be between 1 and 9, not %d", e.level)
	}
	if e.previousEnhancements < 0 || e.previousEnhancements > 3 {
		return 0, fmt.Errorf("previous enhancements must be between 0 and 3, not %d", e.previousEnhancements)
	}
	cost := e.costForBaseEnhancement()
	cost += costForLevel(e.level)
	cost += costForPreviousEnhancements(e.previousEnhancements)
	return cost, nil
}

// BaseEnhancement is an enum of all the base enhancements.
type BaseEnhancement int

const (
	EnhanceMove BaseEnhancement = iota
	EnhanceAttack
	EnhanceRange
	EnhanceShield
	EnhancePush
	EnhancePull
	EnhancePierce
	EnhanceRetaliate
	EnhanceHeal
	EnhanceTarget
	EnhancePoison
	EnhanceWound
	EnhanceMuddle
	EnhanceImmobilize
	EnhanceDisarm
	EnhanceCurse
	EnhanceStrengthen
	EnhanceBless
	EnhanceJump
	EnhanceSpecificElement
	EnhanceAnyElement
	EnhanceSummonsMove
	EnhanceSummonsAttack
	EnhanceSummonsRange
	EnhanceSummonsHP
	EnhanceAddAttackHex
)

// Cost is the cost of an enhancement.
// Probably overkill to have a type for this.
type Cost int

// BaseCostEnhance* are the base costs for each base enhancement.
const (
	BaseCostEnhanceMove            Cost = 30
	BaseCostEnhanceAttack          Cost = 50
	BaseCostEnhanceRange           Cost = 30
	BaseCostEnhanceShield          Cost = 100
	BaseCostEnhancePush            Cost = 30
	BaseCostEnhancePull            Cost = 30
	BaseCostEnhancePierce          Cost = 30
	BaseCostEnhanceRetaliate       Cost = 100
	BaseCostEnhanceHeal            Cost = 30
	BaseCostEnhanceTarget          Cost = 50
	BaseCostEnhancePoison          Cost = 75
	BaseCostEnhanceWound           Cost = 75
	BaseCostEnhanceMuddle          Cost = 50
	BaseCostEnhanceImmobilize      Cost = 100
	BaseCostEnhanceDisarm          Cost = 150
	BaseCostEnhanceCurse           Cost = 75
	BaseCostEnhanceStrengthen      Cost = 50
	BaseCostEnhanceBless           Cost = 50
	BaseCostEnhanceJump            Cost = 50
	BaseCostEnhanceSpecificElement Cost = 100
	BaseCostEnhanceAnyElement      Cost = 150
	BaseCostEnhanceSummonsMove     Cost = 100
	BaseCostEnhanceSummonsAttack   Cost = 100
	BaseCostEnhanceSummonsRange    Cost = 50
	BaseCostEnhanceSummonsHP       Cost = 50
)

func (e enhancement) costForBaseEnhancement() Cost {
	var cost Cost
	switch e.baseEnhancement {
	case EnhanceAddAttackHex:
		return Cost(200 / e.multipleTarget)
	case EnhanceMove:
		cost = BaseCostEnhanceMove
	case EnhanceAttack:
		cost = BaseCostEnhanceAttack
	case EnhanceRange:
		cost = BaseCostEnhanceRange
	case EnhanceShield:
		cost = BaseCostEnhanceShield
	case EnhancePush:
		cost = BaseCostEnhancePush
	case EnhancePull:
		cost = BaseCostEnhancePull
	case EnhancePierce:
		cost = BaseCostEnhancePierce
	case EnhanceRetaliate:
		cost = BaseCostEnhanceRetaliate
	case EnhanceHeal:
		cost = BaseCostEnhanceHeal
	case EnhanceTarget:
		cost = BaseCostEnhanceTarget
	case EnhancePoison:
		cost = BaseCostEnhancePoison
	case EnhanceWound:
		cost = BaseCostEnhanceWound
	case EnhanceMuddle:
		cost = BaseCostEnhanceMuddle
	case EnhanceImmobilize:
		cost = BaseCostEnhanceImmobilize
	case EnhanceDisarm:
		cost = BaseCostEnhanceDisarm
	case EnhanceCurse:
		cost = BaseCostEnhanceCurse
	case EnhanceStrengthen:
		cost = BaseCostEnhanceStrengthen
	case EnhanceBless:
		cost = BaseCostEnhanceBless
	case EnhanceJump:
		cost = BaseCostEnhanceJump
	case EnhanceSpecificElement:
		cost = BaseCostEnhanceSpecificElement
	case EnhanceAnyElement:
		cost = BaseCostEnhanceAnyElement
	case EnhanceSummonsMove:
		cost = BaseCostEnhanceSummonsMove
	case EnhanceSummonsAttack:
		cost = BaseCostEnhanceSummonsAttack
	case EnhanceSummonsRange:
		cost = BaseCostEnhanceSummonsRange
	case EnhanceSummonsHP:
		cost = BaseCostEnhanceSummonsHP
	default:
		cost = 0
	}
	if e.multipleTarget > 1 {
		cost *= 2
	}
	return cost
}

// Level is an enum of all the levels.
// Probably overkill to have an enum for this.
type Level int

const (
	Level1 Level = 1
	Level2 Level = 2
	Level3 Level = 3
	Level4 Level = 4
	Level5 Level = 5
	Level6 Level = 6
	Level7 Level = 7
	Level8 Level = 8
	Level9 Level = 9
)

// CostEnhanceLevel* are the cost premiums for each level.
const (
	CostEnhanceLevel1 Cost = 0
	CostEnhanceLevel2 Cost = 25
	CostEnhanceLevel3 Cost = 50
	CostEnhanceLevel4 Cost = 75
	CostEnhanceLevel5 Cost = 100
	CostEnhanceLevel6 Cost = 125
	CostEnhanceLevel7 Cost = 150
	CostEnhanceLevel8 Cost = 175
	CostEnhanceLevel9 Cost = 200
)

func costForLevel(level Level) Cost {
	switch level {
	case Level1:
		return CostEnhanceLevel1
	case Level2:
		return CostEnhanceLevel2
	case Level3:
		return CostEnhanceLevel3
	case Level4:
		return CostEnhanceLevel4
	case Level5:
		return CostEnhanceLevel5
	case Level6:
		return CostEnhanceLevel6
	case Level7:
		return CostEnhanceLevel7
	case Level8:
		return CostEnhanceLevel8
	case Level9:
		return CostEnhanceLevel9
	default:
		return 0
	}
}

type PreviousEnhancements int

const (
	PreviousEnhancements0 PreviousEnhancements = iota
	PreviousEnhancements1
	PreviousEnhancements2
	PreviousEnhancements3
)

const (
	CostPreviousEnhancements0 Cost = 0
	CostPreviousEnhancements1 Cost = 75
	CostPreviousEnhancements2 Cost = 150
	CostPreviousEnhancements3 Cost = 225
)

func costForPreviousEnhancements(previousEnhancements PreviousEnhancements) Cost {
	switch previousEnhancements {
	case PreviousEnhancements0:
		return CostPreviousEnhancements0
	case PreviousEnhancements1:
		return CostPreviousEnhancements1
	case PreviousEnhancements2:
		return CostPreviousEnhancements2
	case PreviousEnhancements3:
		return CostPreviousEnhancements3
	default:
		return CostPreviousEnhancements0
	}
}
