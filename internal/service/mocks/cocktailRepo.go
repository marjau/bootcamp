package mocks

import (
	ct "github.com/marcos-wz/capstone-go-bootcamp/internal/customtype"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/entity"

	"github.com/stretchr/testify/mock"
)

// CocktailRepo is a mock type for the CocktailRepo dependency
type CocktailRepo struct {
	mock.Mock
}

// ReadAll provides a mock function with given fields:
func (o *CocktailRepo) ReadAll() ([]entity.Cocktail, error) {
	args := o.Called()
	return args.Get(0).([]entity.Cocktail), args.Error(1)
}

// ReadCC provides a mock function with given fields:
func (o *CocktailRepo) ReadCC(nType ct.NumberType, maxJobs, jWorker int) ([]entity.Cocktail, error) {
	args := o.Called(nType, maxJobs, jWorker)
	return args.Get(0).([]entity.Cocktail), args.Error(1)
}

// ReplaceDB provides a mock function with given fields:
func (o *CocktailRepo) ReplaceDB(recs []entity.Cocktail) error {
	args := o.Called(recs)
	return args.Error(0)
}

// Fetch provides a mock function with given fields:
func (o *CocktailRepo) Fetch() ([]entity.Cocktail, error) {
	args := o.Called()
	return args.Get(0).([]entity.Cocktail), args.Error(1)
}

// NewCocktailRepo creates a new instance of the CocktailRepo of type Mock.
func NewCocktailRepo() *CocktailRepo {
	return &CocktailRepo{}
}
