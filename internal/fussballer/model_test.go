package fussballer

import "testing"

func TestPositionIsValid(t *testing.T) {
	validPositions := []Position{
		PositionTorwart,
		PositionVerteidiger,
		PositionMittelfeldspieler,
		PositionStuermer,
	}

	for _, position := range validPositions {
		if !position.IsValid() {
			t.Fatalf("expected position %q to be valid", position)
		}
	}

	if Position("TRAINER").IsValid() {
		t.Fatal("expected unknown position to be invalid")
	}
}
