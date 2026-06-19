package fussballer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound               = errors.New("fussballer not found")
	ErrInvalidSearchParameter = errors.New("invalid search parameter")
)

type Repository struct {
	db *pgxpool.Pool
}

type SearchCriteria struct {
	Nachname      string
	Nationalitaet string
	Position      *Position
	Limit         int
	Offset        int
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByID(ctx context.Context, id int) (*Fussballer, error) {
	const query = `
		SELECT
			f.id,
			f.version,
			f.nachname,
			f.nationalitaet,
			f.position::text,
			f.geburtsdatum,
			f.username,
			f.erzeugt,
			f.aktualisiert,
			a.id,
			a.plz,
			a.ort,
			a.bundesland,
			a.fussballer_id
		FROM fussballer.fussballer f
		LEFT JOIN fussballer.adresse a ON a.fussballer_id = f.id
		WHERE f.id = $1`

	player, err := scanFussballer(r.db.QueryRow(ctx, query, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	awards, err := r.findAuszeichnungen(ctx, id)
	if err != nil {
		return nil, err
	}
	player.Auszeichnungen = awards

	return player, nil
}

func (r *Repository) Find(ctx context.Context, criteria SearchCriteria) ([]Fussballer, error) {
	where, args, err := buildWhereClause(criteria)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
			f.id,
			f.version,
			f.nachname,
			f.nationalitaet,
			f.position::text,
			f.geburtsdatum,
			f.username,
			f.erzeugt,
			f.aktualisiert,
			a.id,
			a.plz,
			a.ort,
			a.bundesland,
			a.fussballer_id
		FROM fussballer.fussballer f
		LEFT JOIN fussballer.adresse a ON a.fussballer_id = f.id` + where + `
		ORDER BY f.id`

	if criteria.Limit > 0 {
		args = append(args, criteria.Limit)
		query += fmt.Sprintf(" LIMIT $%d", len(args))
	}

	if criteria.Offset > 0 {
		args = append(args, criteria.Offset)
		query += fmt.Sprintf(" OFFSET $%d", len(args))
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := make([]Fussballer, 0)
	for rows.Next() {
		player, err := scanFussballer(rows)
		if err != nil {
			return nil, err
		}
		players = append(players, *player)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(players) == 0 {
		return nil, ErrNotFound
	}

	return players, nil
}

func (r *Repository) Count(ctx context.Context, criteria SearchCriteria) (int, error) {
	where, args, err := buildWhereClause(criteria)
	if err != nil {
		return 0, err
	}

	query := "SELECT count(*) FROM fussballer.fussballer f" + where

	var count int
	if err := r.db.QueryRow(ctx, query, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) findAuszeichnungen(ctx context.Context, fussballerID int) ([]Auszeichnung, error) {
	const query = `
		SELECT id, bezeichnung, saison, fussballer_id
		FROM fussballer.auszeichnung
		WHERE fussballer_id = $1
		ORDER BY id`

	rows, err := r.db.Query(ctx, query, fussballerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	awards := make([]Auszeichnung, 0)
	for rows.Next() {
		var award Auszeichnung
		if err := rows.Scan(&award.ID, &award.Bezeichnung, &award.Saison, &award.FussballerID); err != nil {
			return nil, err
		}
		awards = append(awards, award)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return awards, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanFussballer(row rowScanner) (*Fussballer, error) {
	var player Fussballer
	var position sql.NullString
	var addressID sql.NullInt64
	var addressPLZ sql.NullString
	var addressOrt sql.NullString
	var addressBundesland sql.NullString
	var addressFussballerID sql.NullInt64

	err := row.Scan(
		&player.ID,
		&player.Version,
		&player.Nachname,
		&player.Nationalitaet,
		&position,
		&player.Geburtsdatum,
		&player.Username,
		&player.Erzeugt,
		&player.Aktualisiert,
		&addressID,
		&addressPLZ,
		&addressOrt,
		&addressBundesland,
		&addressFussballerID,
	)
	if err != nil {
		return nil, err
	}

	if position.Valid {
		value := Position(position.String)
		player.Position = &value
	}

	if addressID.Valid {
		player.Adresse = &Adresse{
			ID:           int(addressID.Int64),
			PLZ:          addressPLZ.String,
			Ort:          addressOrt.String,
			Bundesland:   addressBundesland.String,
			FussballerID: int(addressFussballerID.Int64),
		}
	}

	return &player, nil
}

func buildWhereClause(criteria SearchCriteria) (string, []any, error) {
	clauses := make([]string, 0, 3)
	args := make([]any, 0, 3)

	if criteria.Nachname != "" {
		args = append(args, criteria.Nachname)
		clauses = append(clauses, fmt.Sprintf("f.nachname = $%d", len(args)))
	}

	if criteria.Nationalitaet != "" {
		args = append(args, criteria.Nationalitaet)
		clauses = append(clauses, fmt.Sprintf("f.nationalitaet = $%d", len(args)))
	}

	if criteria.Position != nil {
		if !criteria.Position.IsValid() {
			return "", nil, ErrInvalidSearchParameter
		}
		args = append(args, string(*criteria.Position))
		clauses = append(clauses, fmt.Sprintf("f.position::text = $%d", len(args)))
	}

	if len(clauses) == 0 {
		return "", args, nil
	}

	return " WHERE " + strings.Join(clauses, " AND "), args, nil
}
