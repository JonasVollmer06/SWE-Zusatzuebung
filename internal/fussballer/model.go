package fussballer

import "time"

type Position string

const (
	PositionTorwart           Position = "TORWART"
	PositionVerteidiger       Position = "VERTEIDIGER"
	PositionMittelfeldspieler Position = "MITTELFELDSPIELER"
	PositionStuermer          Position = "STUERMER"
)

func (p Position) IsValid() bool {
	switch p {
	case PositionTorwart, PositionVerteidiger, PositionMittelfeldspieler, PositionStuermer:
		return true
	default:
		return false
	}
}

type Fussballer struct {
	ID             int            `json:"id"`
	Version        int            `json:"version"`
	Nachname       string         `json:"nachname"`
	Nationalitaet  string         `json:"nationalitaet"`
	Position       *Position      `json:"position,omitempty"`
	Geburtsdatum   time.Time      `json:"geburtsdatum"`
	Username       string         `json:"username"`
	Erzeugt        time.Time      `json:"erzeugt"`
	Aktualisiert   time.Time      `json:"aktualisiert"`
	Adresse        *Adresse       `json:"adresse,omitempty"`
	Auszeichnungen []Auszeichnung `json:"auszeichnungen,omitempty"`
}

type Adresse struct {
	ID           int    `json:"id"`
	PLZ          string `json:"plz"`
	Ort          string `json:"ort"`
	Bundesland   string `json:"bundesland,omitempty"`
	FussballerID int    `json:"fussballerId"`
}

type Auszeichnung struct {
	ID           int    `json:"id"`
	Bezeichnung  string `json:"bezeichnung"`
	Saison       string `json:"saison"`
	FussballerID int    `json:"fussballerId"`
}

type CreateFussballerRequest struct {
	Nachname      string                `json:"nachname" validate:"required"`
	Nationalitaet string                `json:"nationalitaet" validate:"required"`
	Position      Position              `json:"position" validate:"required"`
	Geburtsdatum  time.Time             `json:"geburtsdatum" validate:"required"`
	Username      string                `json:"username" validate:"required"`
	Adresse       *CreateAdresseRequest `json:"adresse,omitempty"`
}

type CreateAdresseRequest struct {
	PLZ        string `json:"plz" validate:"required"`
	Ort        string `json:"ort" validate:"required"`
	Bundesland string `json:"bundesland,omitempty"`
}
