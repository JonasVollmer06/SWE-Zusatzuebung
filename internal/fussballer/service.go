package fussballer

import (
	"context"
	"errors"
)

const (
	DefaultPageSize   = 5
	MaxPageSize       = 100
	DefaultPageNumber = 0
)

var ErrInvalidID = errors.New("invalid id")

type ReadRepository interface {
	FindByID(ctx context.Context, id int) (*Fussballer, error)
	Find(ctx context.Context, criteria SearchCriteria) ([]Fussballer, error)
	Count(ctx context.Context, criteria SearchCriteria) (int, error)
}

type ReadService struct {
	repository ReadRepository
}

type Pageable struct {
	Number int `json:"number"`
	Size   int `json:"size"`
}

type Slice struct {
	Content       []Fussballer `json:"content"`
	TotalElements int          `json:"totalElements"`
}

func NewReadService(repository ReadRepository) *ReadService {
	return &ReadService{repository: repository}
}

func NewPageable(number int, size int) Pageable {
	if number < DefaultPageNumber {
		number = DefaultPageNumber
	}

	if size < 1 || size > MaxPageSize {
		size = DefaultPageSize
	}

	return Pageable{
		Number: number,
		Size:   size,
	}
}

func (s *ReadService) FindByID(ctx context.Context, id int) (*Fussballer, error) {
	if id < 1 {
		return nil, ErrInvalidID
	}

	return s.repository.FindByID(ctx, id)
}

func (s *ReadService) Find(ctx context.Context, criteria SearchCriteria, pageable Pageable) (*Slice, error) {
	if err := validateSearchCriteria(criteria); err != nil {
		return nil, err
	}

	pageable = NewPageable(pageable.Number, pageable.Size)
	criteria.Limit = pageable.Size
	criteria.Offset = pageable.Number * pageable.Size

	players, err := s.repository.Find(ctx, criteria)
	if err != nil {
		return nil, err
	}

	countCriteria := criteria
	countCriteria.Limit = 0
	countCriteria.Offset = 0

	totalElements, err := s.repository.Count(ctx, countCriteria)
	if err != nil {
		return nil, err
	}

	return &Slice{
		Content:       players,
		TotalElements: totalElements,
	}, nil
}

func (s *ReadService) Count(ctx context.Context, criteria SearchCriteria) (int, error) {
	if err := validateSearchCriteria(criteria); err != nil {
		return 0, err
	}

	return s.repository.Count(ctx, criteria)
}

func validateSearchCriteria(criteria SearchCriteria) error {
	if criteria.Position != nil && !criteria.Position.IsValid() {
		return ErrInvalidSearchParameter
	}

	return nil
}
