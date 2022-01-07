package store

import "bootcamp/models"

// *******************************

type Store interface {
	CreateFruit(obj *models.Fruit) error
	ListFruit(filter FruitFilter) *models.Fruits
}

type store struct{}

func NewStore() Store {
	return &store{}
}

// *******************************

func (store) CreateFruit(obj *models.Fruit) error {
	return nil
}

// fruit filter key value unit
type FruitFilter struct {
	Key   string
	Value string
}

func (store) ListFruit(filter FruitFilter) *models.Fruits {
	return &models.Fruits{}
}
