package fussballer

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

func TestApplySearchCriteriaEmptyCriteria(t *testing.T) {
	query := &gorm.DB{}

	result, err := applySearchCriteria(query, SearchCriteria{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != query {
		t.Fatal("expected unchanged query for empty criteria")
	}
}

func TestApplySearchCriteriaRejectsInvalidPosition(t *testing.T) {
	position := Position("TRAINER")

	_, err := applySearchCriteria(&gorm.DB{}, SearchCriteria{Position: &position})

	if !errors.Is(err, ErrInvalidSearchParameter) {
		t.Fatalf("expected ErrInvalidSearchParameter, got %v", err)
	}
}
