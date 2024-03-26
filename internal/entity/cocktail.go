package entity

import "time"

// Cocktail is the representation of a Cocktail recipe used to hold the business logic.
type Cocktail struct {
	ID             int          `json:"id"`
	Name           string       `json:"name"`
	Alcoholic      string       `json:"alcoholic"`
	Category       string       `json:"category"`
	Ingredients    []Ingredient `json:"ingredients"`
	Instructions   string       `json:"instructions"`
	Glass          string       `json:"glass"`
	IBA            string       `json:"iba"`
	ImgAttribution string       `json:"image_attribution"`
	ImgSrc         string       `json:"image_source"`
	Tags           string       `json:"tags"`
	Thumb          string       `json:"thumb"`
	Video          string       `json:"video"`

	SrcDate   time.Time `json:"source_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Ingredient provides the ingredient name and its measure.
type Ingredient struct {
	Name    string `json:"name"`
	Measure string `json:"measure"`
}
