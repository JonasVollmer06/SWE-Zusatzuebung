package fussballer

import (
	"errors"
	"reflect"
	"testing"
)

func TestBuildWhereClauseEmptyCriteria(t *testing.T) {
	where, args, err := buildWhereClause(SearchCriteria{})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if where != "" {
		t.Fatalf("expected empty where clause, got %q", where)
	}
	if len(args) != 0 {
		t.Fatalf("expected no args, got %v", args)
	}
}

func TestBuildWhereClauseWithCriteria(t *testing.T) {
	position := PositionTorwart

	where, args, err := buildWhereClause(SearchCriteria{
		Nachname:      "Neuer",
		Nationalitaet: "Deutschland",
		Position:      &position,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedWhere := " WHERE f.nachname = $1 AND f.nationalitaet = $2 AND f.position::text = $3"
	if where != expectedWhere {
		t.Fatalf("expected where clause %q, got %q", expectedWhere, where)
	}

	expectedArgs := []any{"Neuer", "Deutschland", "TORWART"}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Fatalf("expected args %v, got %v", expectedArgs, args)
	}
}

func TestBuildWhereClauseWithInvalidPosition(t *testing.T) {
	position := Position("TRAINER")

	_, _, err := buildWhereClause(SearchCriteria{Position: &position})

	if !errors.Is(err, ErrInvalidSearchParameter) {
		t.Fatalf("expected ErrInvalidSearchParameter, got %v", err)
	}
}
