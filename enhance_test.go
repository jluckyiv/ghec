package ghec_test

import (
	"testing"

	"github.com/jluckyiv/ghec"
)

// TestEnhanceFromExample tests the example from the README
// and the Gloomhaven rules. Haven't had time to write a
// proper test suite.
func TestEnhanceFromExample(t *testing.T) {
	input := ghec.NewEnhancement(
		ghec.EnhanceAttack).
		WithMultipleTarget(3).
		WithLevel(ghec.Level3)

	actual := input.Cost()
	if actual != ghec.Cost(150) {
		t.Fatalf("Expected 150, got %d", actual)
	}

	input = ghec.NewEnhancement(
		ghec.EnhanceAddAttackHex).
		WithMultipleTarget(3).
		WithLevel(ghec.Level3).
		WithPreviousEnhancements(ghec.PreviousEnhancements1)
	actual = input.Cost()
	if actual != ghec.Cost(191) {
		t.Fatalf("Expected 191, got %d", actual)
	}
}
