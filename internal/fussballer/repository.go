package fussballer

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNotFound               = errors.New("fussballer not found")
	ErrInvalidSearchParameter = errors.New("invalid search parameter")
)

type Repository struct {
	db *gorm.DB
}

type SearchCriteria struct {
	Nachname      string
	Nationalitaet string
	Position      *Position
	Limit         int
	Offset        int
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByID(ctx context.Context, id int) (*Fussballer, error) {
	var player Fussballer

	err := r.withRelations(ctx).
		First(&player, "id = ?", id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &player, nil
}

func (r *Repository) Find(ctx context.Context, criteria SearchCriteria) ([]Fussballer, error) {
	query, err := applySearchCriteria(r.withRelations(ctx), criteria)
	if err != nil {
		return nil, err
	}

	if criteria.Limit > 0 {
		query = query.Limit(criteria.Limit)
	}

	if criteria.Offset > 0 {
		query = query.Offset(criteria.Offset)
	}

	var players []Fussballer
	if err := query.Order("id").Find(&players).Error; err != nil {
		return nil, err
	}

	if len(players) == 0 {
		return nil, ErrNotFound
	}

	return players, nil
}

func (r *Repository) Count(ctx context.Context, criteria SearchCriteria) (int, error) {
	query, err := applySearchCriteria(r.db.WithContext(ctx).Model(&Fussballer{}), criteria)
	if err != nil {
		return 0, err
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *Repository) Create(ctx context.Context, request CreateFussballerRequest) (*Fussballer, error) {
	position := request.Position
	player := Fussballer{
		Nachname:      request.Nachname,
		Nationalitaet: request.Nationalitaet,
		Position:      &position,
		Geburtsdatum:  request.Geburtsdatum,
		Username:      request.Username,
	}

	if request.Adresse != nil {
		player.Adresse = &Adresse{
			PLZ:        request.Adresse.PLZ,
			Ort:        request.Adresse.Ort,
			Bundesland: request.Adresse.Bundesland,
		}
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Create(&player).Error
	})
	if err != nil {
		return nil, err
	}

	return &player, nil
}

func (r *Repository) withRelations(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).
		Preload("Adresse").
		Preload("Auszeichnungen")
}

func applySearchCriteria(query *gorm.DB, criteria SearchCriteria) (*gorm.DB, error) {
	if criteria.Nachname != "" {
		query = query.Where("nachname = ?", criteria.Nachname)
	}

	if criteria.Nationalitaet != "" {
		query = query.Where("nationalitaet = ?", criteria.Nationalitaet)
	}

	if criteria.Position != nil {
		if !criteria.Position.IsValid() {
			return nil, ErrInvalidSearchParameter
		}

		query = query.Where("position = ?", string(*criteria.Position))
	}

	return query, nil
}
