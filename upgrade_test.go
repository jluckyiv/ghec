package ghuc_test

import (
	"testing"

	"github.com/jluckyiv/ghuc"
)

func TestEnhanceFromExample(t *testing.T) {
	input := ghuc.Enhancement{
		BaseEnhancement: ghuc.EnhanceAttack,
		MultipleTarget:  3,
		Level:           ghuc.Level3,
	}
	actual := input.Cost()
	if actual != ghuc.Cost(150) {
		t.Fatalf("Expected 150, got %d", actual)
	}

	input = ghuc.Enhancement{
		BaseEnhancement:      ghuc.EnhanceAddAttackHex,
		MultipleTarget:       3,
		Level:                ghuc.Level3,
		PreviousEnhancements: ghuc.PreviousEnhancements1,
	}
	actual = input.Cost()
	if actual != ghuc.Cost(191) {
		t.Fatalf("Expected 191, got %d", actual)
	}

	// ...
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
	//
}
