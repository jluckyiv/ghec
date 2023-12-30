package ghec

func NewEnhancement(baseEnhancement BaseEnhancement) enhancement {
	return enhancement{
		baseEnhancement: baseEnhancement,
	}
}

type enhancement struct {
	baseEnhancement      BaseEnhancement
	multipleTarget       int
	level                Level
	previousEnhancements PreviousEnhancements
}

func (e enhancement) WithMultipleTarget(multipleTarget int) enhancement {
	e.multipleTarget = multipleTarget
	return e
}

func (e enhancement) WithLevel(level Level) enhancement {
	e.level = level
	return e
}

func (e enhancement) WithPreviousEnhancements(previousEnhancements PreviousEnhancements) enhancement {
	e.previousEnhancements = previousEnhancements
	return e
}

func (e enhancement) Cost() Cost {
	cost := e.costForBaseEnhancement()
	cost += costForLevel(e.level)
	cost += costForPreviousEnhancements(e.previousEnhancements)
	return cost
}

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

type Cost int

const (
	CostEnhanceMove            Cost = 30
	CostEnhanceAttack          Cost = 50
	CostEnhanceRange           Cost = 30
	CostEnhanceShield          Cost = 100
	CostEnhancePush            Cost = 30
	CostEnhancePull            Cost = 30
	CostEnhancePierce          Cost = 30
	CostEnhanceRetaliate       Cost = 100
	CostEnhanceHeal            Cost = 30
	CostEnhanceTarget          Cost = 50
	CostEnhancePoison          Cost = 75
	CostEnhanceWound           Cost = 75
	CostEnhanceMuddle          Cost = 50
	CostEnhanceImmobilize      Cost = 100
	CostEnhanceDisarm          Cost = 150
	CostEnhanceCurse           Cost = 75
	CostEnhanceStrengthen      Cost = 50
	CostEnhanceBless           Cost = 50
	CostEnhanceJump            Cost = 50
	CostEnhanceSpecificElement Cost = 100
	CostEnhanceAnyElement      Cost = 150
	CostEnhanceSummonsMove     Cost = 100
	CostEnhanceSummonsAttack   Cost = 100
	CostEnhanceSummonsRange    Cost = 50
	CostEnhanceSummonsHP       Cost = 50
)

func (e enhancement) costForBaseEnhancement() Cost {
	var cost Cost
	switch e.baseEnhancement {
	case EnhanceAddAttackHex:
		return Cost(200 / e.multipleTarget)
	case EnhanceMove:
		cost = CostEnhanceMove
	case EnhanceAttack:
		cost = CostEnhanceAttack
	case EnhanceRange:
		cost = CostEnhanceRange
	case EnhanceShield:
		cost = CostEnhanceShield
	case EnhancePush:
		cost = CostEnhancePush
	case EnhancePull:
		cost = CostEnhancePull
	case EnhancePierce:
		cost = CostEnhancePierce
	case EnhanceRetaliate:
		cost = CostEnhanceRetaliate
	case EnhanceHeal:
		cost = CostEnhanceHeal
	case EnhanceTarget:
		cost = CostEnhanceTarget
	case EnhancePoison:
		cost = CostEnhancePoison
	case EnhanceWound:
		cost = CostEnhanceWound
	case EnhanceMuddle:
		cost = CostEnhanceMuddle
	case EnhanceImmobilize:
		cost = CostEnhanceImmobilize
	case EnhanceDisarm:
		cost = CostEnhanceDisarm
	case EnhanceCurse:
		cost = CostEnhanceCurse
	case EnhanceStrengthen:
		cost = CostEnhanceStrengthen
	case EnhanceBless:
		cost = CostEnhanceBless
	case EnhanceJump:
		cost = CostEnhanceJump
	case EnhanceSpecificElement:
		cost = CostEnhanceSpecificElement
	case EnhanceAnyElement:
		cost = CostEnhanceAnyElement
	case EnhanceSummonsMove:
		cost = CostEnhanceSummonsMove
	case EnhanceSummonsAttack:
		cost = CostEnhanceSummonsAttack
	case EnhanceSummonsRange:
		cost = CostEnhanceSummonsRange
	case EnhanceSummonsHP:
		cost = CostEnhanceSummonsHP
	default:
		cost = 0
	}
	if e.multipleTarget > 1 {
		cost *= 2
	}
	return cost
}

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
	case 1:
		return CostEnhanceLevel1
	case 2:
		return CostEnhanceLevel2
	case 3:
		return CostEnhanceLevel3
	case 4:
		return CostEnhanceLevel4
	case 5:
		return CostEnhanceLevel5
	case 6:
		return CostEnhanceLevel6
	case 7:
		return CostEnhanceLevel7
	case 8:
		return CostEnhanceLevel8
	case 9:
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
	case 3:
		return CostPreviousEnhancements3
	case 2:
		return CostPreviousEnhancements2
	case 1:
		return CostPreviousEnhancements1
	default:
		return CostPreviousEnhancements0
	}
}
