package fussballer

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
)

type fakeWriteRepository struct {
	createFunc func(ctx context.Context, request CreateFussballerRequest) (*Fussballer, error)
	updateFunc func(ctx context.Context, id int, request UpdateFussballerRequest) (*Fussballer, error)
	deleteFunc func(ctx context.Context, id int) error
	resetFunc  func(ctx context.Context) error
}

func (r *fakeWriteRepository) Create(ctx context.Context, request CreateFussballerRequest) (*Fussballer, error) {
	if r.createFunc == nil {
		return nil, errors.New("unexpected Create call")
	}

	return r.createFunc(ctx, request)
}

func (r *fakeWriteRepository) Update(
	ctx context.Context,
	id int,
	request UpdateFussballerRequest,
) (*Fussballer, error) {
	if r.updateFunc == nil {
		return nil, errors.New("unexpected Update call")
	}

	return r.updateFunc(ctx, id, request)
}

func (r *fakeWriteRepository) Delete(ctx context.Context, id int) error {
	if r.deleteFunc == nil {
		return errors.New("unexpected Delete call")
	}

	return r.deleteFunc(ctx, id)
}

func (r *fakeWriteRepository) Reset(ctx context.Context) error {
	if r.resetFunc == nil {
		return errors.New("unexpected Reset call")
	}

	return r.resetFunc(ctx)
}

func TestWriteServiceCreate(t *testing.T) {
	birthDate := time.Date(2003, time.February, 26, 0, 0, 0, 0, time.UTC)
	request := CreateFussballerRequest{
		Nachname:      " Musiala ",
		Nationalitaet: " Deutschland ",
		Position:      PositionMittelfeldspieler,
		Geburtsdatum:  birthDate,
		Username:      " jamal ",
		Adresse: &CreateAdresseRequest{
			PLZ:        " 80331 ",
			Ort:        " Muenchen ",
			Bundesland: " Bayern ",
		},
	}
	expectedRequest := CreateFussballerRequest{
		Nachname:      "Musiala",
		Nationalitaet: "Deutschland",
		Position:      PositionMittelfeldspieler,
		Geburtsdatum:  birthDate,
		Username:      "jamal",
		Adresse: &CreateAdresseRequest{
			PLZ:        "80331",
			Ort:        "Muenchen",
			Bundesland: "Bayern",
		},
	}
	expected := &Fussballer{ID: 1008, Nachname: "Musiala"}

	service := NewWriteService(&fakeWriteRepository{
		createFunc: func(_ context.Context, got CreateFussballerRequest) (*Fussballer, error) {
			if !reflect.DeepEqual(got, expectedRequest) {
				t.Fatalf("expected request %+v, got %+v", expectedRequest, got)
			}

			return expected, nil
		},
	})

	player, err := service.Create(context.Background(), request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if player != expected {
		t.Fatalf("expected player %v, got %v", expected, player)
	}
}

func TestWriteServiceCreateRejectsMissingRequiredField(t *testing.T) {
	service := NewWriteService(&fakeWriteRepository{})

	_, err := service.Create(context.Background(), CreateFussballerRequest{
		Nationalitaet: "Deutschland",
		Position:      PositionMittelfeldspieler,
		Geburtsdatum:  time.Date(2003, time.February, 26, 0, 0, 0, 0, time.UTC),
		Username:      "jamal",
	})

	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestWriteServiceCreateRejectsInvalidPosition(t *testing.T) {
	service := NewWriteService(&fakeWriteRepository{})

	_, err := service.Create(context.Background(), CreateFussballerRequest{
		Nachname:      "Musiala",
		Nationalitaet: "Deutschland",
		Position:      Position("TRAINER"),
		Geburtsdatum:  time.Date(2003, time.February, 26, 0, 0, 0, 0, time.UTC),
		Username:      "jamal",
	})

	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestWriteServiceCreateRejectsInvalidAddress(t *testing.T) {
	service := NewWriteService(&fakeWriteRepository{})

	_, err := service.Create(context.Background(), CreateFussballerRequest{
		Nachname:      "Musiala",
		Nationalitaet: "Deutschland",
		Position:      PositionMittelfeldspieler,
		Geburtsdatum:  time.Date(2003, time.February, 26, 0, 0, 0, 0, time.UTC),
		Username:      "jamal",
		Adresse: &CreateAdresseRequest{
			PLZ: "80331",
		},
	})

	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestWriteServiceCreateReturnsRepositoryError(t *testing.T) {
	expectedErr := errors.New("database failed")
	service := NewWriteService(&fakeWriteRepository{
		createFunc: func(_ context.Context, _ CreateFussballerRequest) (*Fussballer, error) {
			return nil, expectedErr
		},
	})

	_, err := service.Create(context.Background(), CreateFussballerRequest{
		Nachname:      "Musiala",
		Nationalitaet: "Deutschland",
		Position:      PositionMittelfeldspieler,
		Geburtsdatum:  time.Date(2003, time.February, 26, 0, 0, 0, 0, time.UTC),
		Username:      "jamal",
	})

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
}

func TestWriteServiceUpdate(t *testing.T) {
	birthDate := time.Date(2003, time.February, 26, 0, 0, 0, 0, time.UTC)
	request := UpdateFussballerRequest{
		Nachname:      " Musiala ",
		Nationalitaet: " Deutschland ",
		Position:      PositionMittelfeldspieler,
		Geburtsdatum:  birthDate,
		Username:      " jamal ",
		Adresse: &CreateAdresseRequest{
			PLZ: " 80331 ",
			Ort: " Muenchen ",
		},
	}
	expectedRequest := UpdateFussballerRequest{
		Nachname:      "Musiala",
		Nationalitaet: "Deutschland",
		Position:      PositionMittelfeldspieler,
		Geburtsdatum:  birthDate,
		Username:      "jamal",
		Adresse: &CreateAdresseRequest{
			PLZ: "80331",
			Ort: "Muenchen",
		},
	}
	expected := &Fussballer{ID: 20, Nachname: "Musiala", Version: 1}

	service := NewWriteService(&fakeWriteRepository{
		updateFunc: func(_ context.Context, id int, got UpdateFussballerRequest) (*Fussballer, error) {
			if id != 20 {
				t.Fatalf("expected id 20, got %d", id)
			}
			if !reflect.DeepEqual(got, expectedRequest) {
				t.Fatalf("expected request %+v, got %+v", expectedRequest, got)
			}

			return expected, nil
		},
	})

	player, err := service.Update(context.Background(), 20, request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if player != expected {
		t.Fatalf("expected player %v, got %v", expected, player)
	}
}

func TestWriteServiceUpdateRejectsInvalidID(t *testing.T) {
	service := NewWriteService(&fakeWriteRepository{})

	_, err := service.Update(context.Background(), 0, UpdateFussballerRequest{})

	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("expected ErrInvalidID, got %v", err)
	}
}

func TestWriteServiceUpdateRejectsInvalidPosition(t *testing.T) {
	service := NewWriteService(&fakeWriteRepository{})

	_, err := service.Update(context.Background(), 20, UpdateFussballerRequest{
		Nachname:      "Musiala",
		Nationalitaet: "Deutschland",
		Position:      Position("TRAINER"),
		Geburtsdatum:  time.Date(2003, time.February, 26, 0, 0, 0, 0, time.UTC),
		Username:      "jamal",
	})

	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected ErrValidation, got %v", err)
	}
}

func TestWriteServiceDelete(t *testing.T) {
	service := NewWriteService(&fakeWriteRepository{
		deleteFunc: func(_ context.Context, id int) error {
			if id != 1000 {
				t.Fatalf("expected id 1000, got %d", id)
			}

			return nil
		},
	})

	if err := service.Delete(context.Background(), 1000); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWriteServiceDeleteRejectsInvalidID(t *testing.T) {
	service := NewWriteService(&fakeWriteRepository{})

	err := service.Delete(context.Background(), 0)

	if !errors.Is(err, ErrInvalidID) {
		t.Fatalf("expected ErrInvalidID, got %v", err)
	}
}

func TestWriteServiceReset(t *testing.T) {
	called := false
	service := NewWriteService(&fakeWriteRepository{
		resetFunc: func(_ context.Context) error {
			called = true
			return nil
		},
	})

	if err := service.Reset(context.Background()); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Fatal("expected repository reset to be called")
	}
}
