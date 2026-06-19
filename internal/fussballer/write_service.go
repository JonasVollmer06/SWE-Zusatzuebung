package fussballer

import (
	"context"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var ErrValidation = errors.New("validation failed")

type WriteRepository interface {
	Create(ctx context.Context, request CreateFussballerRequest) (*Fussballer, error)
}

type WriteService struct {
	repository WriteRepository
	validator  *validator.Validate
}

func NewWriteService(repository WriteRepository) *WriteService {
	return &WriteService{
		repository: repository,
		validator:  validator.New(),
	}
}

func (s *WriteService) Create(ctx context.Context, request CreateFussballerRequest) (*Fussballer, error) {
	request = normalizeCreateRequest(request)

	if err := s.validator.Struct(request); err != nil {
		return nil, ErrValidation
	}

	if !request.Position.IsValid() {
		return nil, ErrValidation
	}

	return s.repository.Create(ctx, request)
}

func normalizeCreateRequest(request CreateFussballerRequest) CreateFussballerRequest {
	request.Nachname = strings.TrimSpace(request.Nachname)
	request.Nationalitaet = strings.TrimSpace(request.Nationalitaet)
	request.Username = strings.TrimSpace(request.Username)

	if request.Adresse != nil {
		request.Adresse.PLZ = strings.TrimSpace(request.Adresse.PLZ)
		request.Adresse.Ort = strings.TrimSpace(request.Adresse.Ort)
		request.Adresse.Bundesland = strings.TrimSpace(request.Adresse.Bundesland)
	}

	return request
}
