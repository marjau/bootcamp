package service

import (
	"errors"
	"testing"

	"github.com/marcos-wz/capstone-go-bootcamp/internal/entity"

	"github.com/stretchr/testify/assert"
)

var (
	testRepoErr = errors.New("test repository error")

	testCocktailsAll = []entity.Cocktail{
		{ID: 1, Name: "Foo", Alcoholic: "Alcoholic", Category: "Foo Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "soda", Measure: "80ml"}}},
		{ID: 2, Name: "Bar", Alcoholic: "Non alcoholic", Category: "Some Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "water", Measure: "50ml"}}},
		{ID: 3, Name: "Baz", Alcoholic: "Alcoholic", Category: "Some Category", Glass: "Cocktail glass", Ingredients: []entity.Ingredient{{Name: "soda", Measure: "100ml"}}},
	}
)

func TestNewCocktailFilter(t *testing.T) {
	tests := []struct {
		name   string
		filter string
		exp    cocktailFilter
	}{
		{
			name:   "Empty",
			filter: "",
			exp:    invalidFltr,
		},
		{
			name:   "Arbitrary",
			filter: "fooFltr",
			exp:    invalidFltr,
		},
		{
			name:   "ID",
			filter: "id",
			exp:    idFltr,
		},
		{
			name:   "Name",
			filter: "name",
			exp:    nameFltr,
		},
		{
			name:   "Alcoholic",
			filter: "alcOholIc",
			exp:    alcoholicFltr,
		},
		{
			name:   "Category",
			filter: "category",
			exp:    categoryFltr,
		},
		{
			name:   "Ingredient",
			filter: "ingredient",
			exp:    ingredientFltr,
		},
		{
			name:   "Glass",
			filter: "Glass",
			exp:    glassFltr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := newCocktailFltr(tt.filter)
			assert.IsType(t, cocktailFilter(""), out)
			assert.Equal(t, tt.exp, out)
		})
	}
}

func TestFindCocktail(t *testing.T) {
	type args struct {
		id   int
		recs []entity.Cocktail
	}
	tests := []struct {
		name  string
		args  args
		index int
		found bool
	}{
		{
			name: "Nil Records",
			args: args{
				id:   1,
				recs: nil,
			},
			index: 0,
			found: false,
		},
		{
			name: "Records empty",
			args: args{
				id:   1,
				recs: []entity.Cocktail{},
			},
			index: 0,
			found: false,
		},
		{
			name: "Not found",
			args: args{
				id: 123,
				recs: []entity.Cocktail{
					{ID: 1, Name: "foo"},
					{ID: 50, Name: "bar"},
					{ID: 100, Name: "baz"},
				},
			},
			index: 0,
			found: false,
		},
		{
			name: "Exists",
			args: args{
				id: 50,
				recs: []entity.Cocktail{
					{ID: 1, Name: "foo"},
					{ID: 50, Name: "bar"},
					{ID: 100, Name: "baz"},
				},
			},
			index: 1,
			found: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, found := findCocktail(tt.args.id, tt.args.recs)
			assert.Equal(t, tt.index, index)
			assert.Equal(t, tt.found, found)
		})
	}
}
