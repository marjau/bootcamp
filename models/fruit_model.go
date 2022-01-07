package models

import "time"

type Fruit struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Color       string  `json:"color"`
	Unit        string  `json:"unit"`     // kgs, lbs
	Price       float64 `json:"price"`    // price per unit
	Stock       int     `json:"stock"`    // stock
	Caducate    int     `json:"caducate"` // number days to the fruit be caducated
	Country     string  `json:"country"`  // country imported from
	CreateAt    time.Time
}

type Fruits []Fruit
