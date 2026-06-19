package fussballer

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type fakeReadRepository struct {
	findByIDFunc func(ctx context.Context, id int) (*Fussballer, error)
	findFunc     func(ctx context.Context, criteria SearchCriteria) ([]Fussballer, error)
	countFunc    func(ctx context.Context, criteria SearchCriteria) (int, error)
}

func (r *fakeReadRepository) FindByID(ctx context.Context, id int) (*Fussballer, error) {
	if r.findByIDFunc == nil {
		return nil, errors.New("unexpected FindByID call")
	}

	return r.findByIDFunc(ctx, id)
}

func (r *fakeReadRepository) Find(ctx context.Context, criteria SearchCriteria) ([]Fussballer, error) {
	if r.findFunc == nil {
		return nil, errors.New("unexpected Find call")
	}

	return r.findFunc(ctx, criteria)
}

func (r *fakeReadRepository) Count(ctx context.Context, criteria SearchCriteria) (int, error) {
	if r.countFunc == nil {
		return 0, errors.New("unexpected Count call")
	}

	return r.countFunc(ctx, criteria)
}

func TestReadServiceFindByID(t *testing.T) {
	expected := &Fussballer{ID: 1000, Nachname: "Neuer"}

	service := NewReadService(&fakeReadRepository{
		findByIDFunc: func(_ context.Context, id int) (*Fussballer, error) {
			if id != expected.ID {
				t.Fatalf("expected id %d, got %d", expected.ID, id)
			}

			return expected, nil
		},
	})

	player, err := service.FindByID(context.Background(), expected.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if player != expected {
		t.Fatalf("expected player %v, got %v", expected, player)
	}
}

func TestReadServiceFindByIDRejectsInvalidID(t *testing.T) {
	service := NewReadService(&fakeReadRepository{})

	_, err := service.FindByID(context.Background(), 0)

	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("expected ErrInvalidID, got %v", err)
	}
}

func TestReadServiceFindAppliesPageable(t *testing.T) {
	position := PositionTorwart
	expectedCriteria := SearchCriteria{
		Nachname: "Neuer",
		Position: &position,
		Limit:    10,
		Offset:   20,
	}
	expectedCountCriteria := SearchCriteria{
		Nachname: "Neuer",
		Position: &position,
	}
	expectedPlayers := []Fussballer{{ID: 1000, Nachname: "Neuer"}}

	service := NewReadService(&fakeReadRepository{
		findFunc: func(_ context.Context, criteria SearchCriteria) ([]Fussballer, error) {
			if !reflect.DeepEqual(criteria, expectedCriteria) {
				t.Fatalf("expected criteria %+v, got %+v", expectedCriteria, criteria)
			}

			return expectedPlayers, nil
		},
		countFunc: func(_ context.Context, criteria SearchCriteria) (int, error) {
			if !reflect.DeepEqual(criteria, expectedCountCriteria) {
				t.Fatalf("expected count criteria %+v, got %+v", expectedCountCriteria, criteria)
			}

			return 42, nil
		},
	})

	result, err := service.Find(
		context.Background(),
		SearchCriteria{Nachname: "Neuer", Position: &position},
		Pageable{Number: 2, Size: 10},
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !reflect.DeepEqual(result.Content, expectedPlayers) {
		t.Fatalf("expected content %v, got %v", expectedPlayers, result.Content)
	}
	if result.TotalElements != 42 {
		t.Fatalf("expected totalElements 42, got %d", result.TotalElements)
	}
}

func TestReadServiceFindRejectsInvalidPosition(t *testing.T) {
	position := Position("TRAINER")
	service := NewReadService(&fakeReadRepository{})

	_, err := service.Find(
		context.Background(),
		SearchCriteria{Position: &position},
		Pageable{Number: 0, Size: 5},
	)

	if !errors.Is(err, ErrInvalidSearchParameter) {
		t.Fatalf("expected ErrInvalidSearchParameter, got %v", err)
	}
}
