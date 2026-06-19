package fussballer

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

func (r *Repository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&Fussballer{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *Repository) Reset(ctx context.Context) error {
	seedDir, err := findSeedCSVDir()
	if err != nil {
		return err
	}

	players, err := readSeedCSV(filepath.Join(seedDir, "fussballer.csv"))
	if err != nil {
		return err
	}
	addresses, err := readSeedCSV(filepath.Join(seedDir, "adresse.csv"))
	if err != nil {
		return err
	}
	awards, err := readSeedCSV(filepath.Join(seedDir, "auszeichnung.csv"))
	if err != nil {
		return err
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(
			"TRUNCATE TABLE fussballer.fussballer_file, fussballer.adresse, fussballer.auszeichnung, fussballer.fussballer RESTART IDENTITY CASCADE",
		).Error; err != nil {
			return err
		}

		for _, player := range players {
			if len(player) != 9 {
				return fmt.Errorf("expected 9 fussballer columns, got %d", len(player))
			}

			if err := tx.Exec(`
				INSERT INTO fussballer.fussballer
					(id, version, nachname, nationalitaet, position, geburtsdatum, username, erzeugt, aktualisiert)
					OVERRIDING SYSTEM VALUE
				VALUES (?, ?, ?, ?, ?::fussballer.position_enum, ?::date, ?, ?::timestamp, ?::timestamp)
			`, player[0], player[1], player[2], player[3], player[4], player[5], player[6], player[7], player[8]).Error; err != nil {
				return err
			}
		}

		for _, address := range addresses {
			if len(address) != 5 {
				return fmt.Errorf("expected 5 address columns, got %d", len(address))
			}

			if err := tx.Exec(`
				INSERT INTO fussballer.adresse
					(id, plz, ort, bundesland, fussballer_id)
					OVERRIDING SYSTEM VALUE
				VALUES (?, ?, ?, ?, ?)
			`, address[0], address[1], address[2], address[3], address[4]).Error; err != nil {
				return err
			}
		}

		for _, award := range awards {
			if len(award) != 4 {
				return fmt.Errorf("expected 4 award columns, got %d", len(award))
			}

			if err := tx.Exec(`
				INSERT INTO fussballer.auszeichnung
					(id, bezeichnung, saison, fussballer_id)
					OVERRIDING SYSTEM VALUE
				VALUES (?, ?, ?, ?)
			`, award[0], award[1], award[2], award[3]).Error; err != nil {
				return err
			}
		}

		return nil
	})
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

func findSeedCSVDir() (string, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		seedDir := filepath.Join(workingDir, "extras", "compose", "postgres", "init", "fussballer", "csv")
		if _, err := os.Stat(filepath.Join(seedDir, "fussballer.csv")); err == nil {
			return seedDir, nil
		}

		parent := filepath.Dir(workingDir)
		if parent == workingDir {
			return "", fmt.Errorf("seed CSV directory not found")
		}
		workingDir = parent
	}
}

func readSeedCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) < 1 {
		return nil, fmt.Errorf("empty seed CSV: %s", path)
	}

	return records[1:], nil
}
