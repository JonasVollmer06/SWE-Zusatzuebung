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
	ID             int            `json:"id" gorm:"column:id;primaryKey"`
	Version        int            `json:"version" gorm:"column:version"`
	Nachname       string         `json:"nachname" gorm:"column:nachname"`
	Nationalitaet  string         `json:"nationalitaet" gorm:"column:nationalitaet"`
	Position       *Position      `json:"position,omitempty" gorm:"column:position;type:fussballer.position_enum"`
	Geburtsdatum   time.Time      `json:"geburtsdatum" gorm:"column:geburtsdatum"`
	Username       string         `json:"username" gorm:"column:username"`
	Erzeugt        time.Time      `json:"erzeugt" gorm:"column:erzeugt"`
	Aktualisiert   time.Time      `json:"aktualisiert" gorm:"column:aktualisiert"`
	Adresse        *Adresse       `json:"adresse,omitempty" gorm:"foreignKey:FussballerID"`
	Auszeichnungen []Auszeichnung `json:"auszeichnungen,omitempty" gorm:"foreignKey:FussballerID"`
}

func (Fussballer) TableName() string {
	return "fussballer.fussballer"
}

type Adresse struct {
	ID           int    `json:"id" gorm:"column:id;primaryKey"`
	PLZ          string `json:"plz" gorm:"column:plz"`
	Ort          string `json:"ort" gorm:"column:ort"`
	Bundesland   string `json:"bundesland,omitempty" gorm:"column:bundesland"`
	FussballerID int    `json:"fussballerId" gorm:"column:fussballer_id"`
}

func (Adresse) TableName() string {
	return "fussballer.adresse"
}

type Auszeichnung struct {
	ID           int    `json:"id" gorm:"column:id;primaryKey"`
	Bezeichnung  string `json:"bezeichnung" gorm:"column:bezeichnung"`
	Saison       string `json:"saison" gorm:"column:saison"`
	FussballerID int    `json:"fussballerId" gorm:"column:fussballer_id"`
}

func (Auszeichnung) TableName() string {
	return "fussballer.auszeichnung"
}

type CreateFussballerRequest struct {
	Nachname      string                `json:"nachname" validate:"required"`
	Nationalitaet string                `json:"nationalitaet" validate:"required"`
	Position      Position              `json:"position" validate:"required"`
	Geburtsdatum  time.Time             `json:"geburtsdatum" validate:"required"`
	Username      string                `json:"username" validate:"required"`
	Adresse       *CreateAdresseRequest `json:"adresse,omitempty"`
}

type UpdateFussballerRequest struct {
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
