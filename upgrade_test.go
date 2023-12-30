package ghuc_test

import (
	"testing"

	"github.com/jluckyiv/ghuc"
)

func TestEnhanceFromExample(t *testing.T) {
	input := ghuc.NewEnhancement(
		ghuc.EnhanceAttack).
		WithMultipleTarget(3).
		WithLevel(ghuc.Level3)

	actual := input.Cost()
	if actual != ghuc.Cost(150) {
		t.Fatalf("Expected 150, got %d", actual)
	}

	input = ghuc.NewEnhancement(
		ghuc.EnhanceAddAttackHex).
		WithMultipleTarget(3).
		WithLevel(ghuc.Level3).
		WithPreviousEnhancements(ghuc.PreviousEnhancements1)
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
