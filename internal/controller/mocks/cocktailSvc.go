package mocks

import (
	ct "github.com/marcos-wz/capstone-go-bootcamp/internal/customtype"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/entity"

	"github.com/stretchr/testify/mock"
)

// CocktailSvc is a mock type for the CocktailSvc dependency
type CocktailSvc struct {
	mock.Mock
}

// GetFiltered provides a mock function with given fields:
func (o *CocktailSvc) GetFiltered(filter, value string) ([]entity.Cocktail, error) {
	args := o.Called(filter, value)
	return args.Get(0).([]entity.Cocktail), args.Error(1)
}

// GetAll provides a mock function with given fields:
func (o *CocktailSvc) GetAll() ([]entity.Cocktail, error) {
	args := o.Called()
	return args.Get(0).([]entity.Cocktail), args.Error(1)
}

// GetCC provides a mock function with given fields:
func (o *CocktailSvc) GetCC(nType, jobs, jWorker string) ([]entity.Cocktail, error) {
	args := o.Called(nType, jobs, jWorker)
	return args.Get(0).([]entity.Cocktail), args.Error(1)
}

// UpdateDB provides a mock function with given fields:
func (o *CocktailSvc) UpdateDB() (ct.DBOpsSummary, error) {
	args := o.Called()
	return args.Get(0).(ct.DBOpsSummary), args.Error(1)
}

// NewCocktailSvc creates a new instance of the CocktailSvc of type Mock.
func NewCocktailSvc() *CocktailSvc {
	return &CocktailSvc{}
}
