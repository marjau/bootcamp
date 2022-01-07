package services

import (
	"bootcamp/models"
	"bootcamp/store"

	"github.com/stretchr/testify/mock"
)

type mockFruitStore struct {
	mock.Mock
}

func (o mockFruitStore) CreateFruit(obj *models.Fruit) error {
	args := o.Called()
	return args.Error(0)
}

func (o mockFruitStore) ListFruit(filter store.FruitFilter) *models.Fruits {
	args := o.Called()

	return args.Get(0).(*models.Fruits)
}
