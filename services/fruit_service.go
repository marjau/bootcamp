package services

import (
	"bootcamp/models"
	"bootcamp/store"
	"log"
)

// ***************************

type FruitService interface {
	Write(fruit *models.Fruit) error
	List(filter *store.FruitFilter) (*models.Fruits, error)
}

type fruitService struct {
	store store.Store
}

func NewFruitService(s store.Store) FruitService {
	return &fruitService{s}
}

// ***************************

func (f *fruitService) Write(fruit *models.Fruit) error {
	log.Printf("SERVICE WRITE, Fruit: %+v", fruit)
	err := f.store.CreateFruit(fruit)
	if err != nil {
		return err
	}
	return nil
}

func (f *fruitService) List(filter *store.FruitFilter) (*models.Fruits, error) {
	list := f.store.ListFruit(*filter)
	log.Println("FRUIT SERVICE LIST")
	return list, nil
}
