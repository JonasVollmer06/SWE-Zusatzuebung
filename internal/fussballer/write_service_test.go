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
}

func (r *fakeWriteRepository) Create(ctx context.Context, request CreateFussballerRequest) (*Fussballer, error) {
	if r.createFunc == nil {
		return nil, errors.New("unexpected Create call")
	}

	return r.createFunc(ctx, request)
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
