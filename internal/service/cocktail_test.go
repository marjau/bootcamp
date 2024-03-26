package service

import (
	"strconv"
	"testing"
	"time"

	ct "github.com/marcos-wz/capstone-go-bootcamp/internal/customtype"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/entity"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ CocktailRepo = &mocks.CocktailRepo{}

func TestNewCocktail(t *testing.T) {
	repo := mocks.NewCocktailRepo()
	require.NotNil(t, repo)
	out := NewCocktail(repo)
	assert.IsType(t, Cocktail{}, out)
}

func TestCocktail_GetFiltered(t *testing.T) {
	type args struct {
		filter string
		value  string
	}
	type repo struct {
		resp []entity.Cocktail
		err  error
	}
	tests := []struct {
		name string
		args args
		exp  []entity.Cocktail
		err  error
		repo repo
	}{
		{
			name: "Repository error",
			args: args{
				filter: idFltr.String(),
				value:  "4",
			},
			exp: nil,
			err: testRepoErr,
			repo: repo{
				resp: nil,
				err:  testRepoErr,
			},
		},
		{
			name: "Value empty",
			args: args{
				filter: idFltr.String(),
				value:  "",
			},
			exp:  nil,
			err:  ErrFltrValueEmpty,
			repo: repo{},
		},
		{
			name: "Arbitrary",
			args: args{filter: "foo", value: "foo-value"},
			exp:  nil,
			err:  ErrFltrInvalid,
			repo: repo{},
		},
		{
			name: "Value not found",
			args: args{filter: idFltr.String(), value: "123456"},
			exp:  []entity.Cocktail{},
			err:  nil,
			repo: repo{
				resp: []entity.Cocktail{},
				err:  nil,
			},
		},
		{
			name: "Bad ID",
			args: args{filter: idFltr.String(), value: "foo-id"},
			exp:  []entity.Cocktail{},
			err:  strconv.ErrSyntax,
			repo: repo{
				resp: testCocktailsAll,
				err:  nil,
			},
		},
		{
			name: "ID",
			args: args{filter: idFltr.String(), value: "2"},
			exp: []entity.Cocktail{
				{ID: 2, Name: "Bar", Alcoholic: "Non alcoholic", Category: "Some Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "water", Measure: "50ml"}}},
			},
			err: nil,
			repo: repo{
				resp: testCocktailsAll,
				err:  nil,
			},
		},
		{
			name: "Name",
			args: args{filter: nameFltr.String(), value: "bAz"},
			exp: []entity.Cocktail{
				{ID: 3, Name: "Baz", Alcoholic: "Alcoholic", Category: "Some Category", Glass: "Cocktail glass", Ingredients: []entity.Ingredient{{Name: "soda", Measure: "100ml"}}},
			},
			err: nil,
			repo: repo{
				resp: testCocktailsAll,
				err:  nil,
			},
		},
		{
			name: "Alcoholic",
			args: args{filter: alcoholicFltr.String(), value: "non alcoholic"},
			exp: []entity.Cocktail{
				{ID: 2, Name: "Bar", Alcoholic: "Non alcoholic", Category: "Some Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "water", Measure: "50ml"}}},
			},
			err: nil,
			repo: repo{
				resp: testCocktailsAll,
				err:  nil,
			},
		},
		{
			name: "Category",
			args: args{filter: categoryFltr.String(), value: "some"},
			exp: []entity.Cocktail{
				{ID: 2, Name: "Bar", Alcoholic: "Non alcoholic", Category: "Some Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "water", Measure: "50ml"}}},
				{ID: 3, Name: "Baz", Alcoholic: "Alcoholic", Category: "Some Category", Glass: "Cocktail glass", Ingredients: []entity.Ingredient{{Name: "soda", Measure: "100ml"}}},
			},
			err: nil,
			repo: repo{
				resp: testCocktailsAll,
				err:  nil,
			},
		},
		{
			name: "Ingredients",
			args: args{filter: ingredientFltr.String(), value: "soda"},
			exp: []entity.Cocktail{
				{ID: 1, Name: "Foo", Alcoholic: "Alcoholic", Category: "Foo Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "soda", Measure: "80ml"}}},
				{ID: 3, Name: "Baz", Alcoholic: "Alcoholic", Category: "Some Category", Glass: "Cocktail glass", Ingredients: []entity.Ingredient{{Name: "soda", Measure: "100ml"}}},
			},
			err: nil,
			repo: repo{
				resp: testCocktailsAll,
				err:  nil,
			},
		},
		{
			name: "Glass",
			args: args{filter: glassFltr.String(), value: "shot"},
			exp: []entity.Cocktail{
				{ID: 1, Name: "Foo", Alcoholic: "Alcoholic", Category: "Foo Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "soda", Measure: "80ml"}}},
				{ID: 2, Name: "Bar", Alcoholic: "Non alcoholic", Category: "Some Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "water", Measure: "50ml"}}},
			},
			err: nil,
			repo: repo{
				resp: testCocktailsAll,
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mRepo := mocks.NewCocktailRepo()
			mRepo.On("ReadAll").Return(tt.repo.resp, tt.repo.err)
			svc := NewCocktail(mRepo)
			require.NotEqual(t, Cocktail{}, svc)

			out, err := svc.GetFiltered(tt.args.filter, tt.args.value)
			if tt.err != nil {
				require.NotNil(t, err)
				require.Nil(t, out)
				assert.ErrorIs(t, err, tt.err)
				return
			}
			require.Nil(t, err)
			require.NotNil(t, out)
			assert.Len(t, tt.exp, len(out))
			assert.Equal(t, tt.exp, out)
		})
	}
}

func TestCocktail_GetAll(t *testing.T) {
	type repo struct {
		resp []entity.Cocktail
		err  error
	}
	tests := []struct {
		name string
		exp  []entity.Cocktail
		err  error
		repo repo
	}{
		{
			name: "Repository Error",
			exp:  nil,
			err:  testRepoErr,
			repo: repo{
				resp: nil,
				err:  testRepoErr,
			},
		},
		{
			name: "Not Found",
			exp:  []entity.Cocktail{},
			err:  nil,
			repo: repo{
				resp: []entity.Cocktail{},
				err:  nil,
			},
		},
		{
			name: "All records",
			exp:  testCocktailsAll,
			err:  nil,
			repo: repo{
				resp: testCocktailsAll,
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mRepo := mocks.NewCocktailRepo()
			mRepo.On("ReadAll").Return(tt.repo.resp, tt.repo.err)
			svc := NewCocktail(mRepo)
			require.NotNil(t, svc)

			out, err := svc.GetAll()
			if tt.err != nil {
				require.NotNil(t, err)
				assert.Nil(t, out)
				assert.ErrorIs(t, err, tt.err)
				return
			}
			require.Nil(t, err)
			require.NotNil(t, out)
			assert.Len(t, tt.exp, len(out))
			assert.Equal(t, tt.exp, out)
		})
	}
}

func TestCocktail_GetCC(t *testing.T) {
	type repoArgs struct {
		nType   ct.NumberType
		jobs    int
		jWorker int
	}
	type repo struct {
		args repoArgs
		resp []entity.Cocktail
		err  error
	}
	type args struct {
		nType   string
		jobs    string
		jWorker string
	}
	tests := []struct {
		name string
		args args
		exp  []entity.Cocktail
		err  error
		repo repo
	}{
		{
			name: "Repository Error",
			args: args{nType: "even", jobs: "10", jWorker: "2"},
			exp:  nil,
			err:  testRepoErr,
			repo: repo{
				args: repoArgs{nType: ct.EvenNum, jobs: 10, jWorker: 2},
				resp: nil,
				err:  testRepoErr,
			},
		},
		{
			name: "Invalid number type",
			args: args{nType: "foo", jobs: "10", jWorker: "2"},
			exp:  nil,
			err:  ErrInvalidNumType,
			repo: repo{},
		},
		{
			name: "Invalid jobs",
			args: args{nType: "even", jobs: "bad-num", jWorker: "2"},
			exp:  nil,
			err:  strconv.ErrSyntax,
			repo: repo{},
		},
		{
			name: "Invalid jobs per worker",
			args: args{nType: "even", jobs: "10", jWorker: "bad-num"},
			exp:  nil,
			err:  strconv.ErrSyntax,
			repo: repo{},
		},
		{
			name: "Zero jobs",
			args: args{nType: "odd", jobs: "0", jWorker: "2"},
			exp:  nil,
			err:  ErrZeroValue,
			repo: repo{},
		},
		{
			name: "Zero Jobs per worker",
			args: args{nType: "odd", jobs: "2", jWorker: "0"},
			exp:  nil,
			err:  ErrZeroValue,
			repo: repo{},
		},
		{
			name: "Jobs per worker higher",
			args: args{nType: "odd", jobs: "2", jWorker: "10"},
			exp:  nil,
			err:  ErrJobsWorkerHigher,
			repo: repo{},
		},
		{
			name: "Valid",
			args: args{nType: "even", jobs: "8", jWorker: "2"},
			exp: []entity.Cocktail{
				{ID: 17222, Name: "A1", Alcoholic: "Alcoholic", Category: "Cocktail", Ingredients: []entity.Ingredient{{Name: "Gin", Measure: "1 3/4 shot "}, {Name: "Grand Marnier", Measure: "1 Shot "}, {Name: "Lemon Juice", Measure: "1/4 Shot"}, {Name: "Grenadine", Measure: "1/8 Shot"}}, Instructions: "Pour all ingredients into a cocktail shaker, mix and serve over ice into a chilled glass.", Glass: "Cocktail glass", IBA: "", ImgAttribution: "", ImgSrc: "", Tags: "", Thumb: "https://www.thecocktaildb.com/images/media/drink/2x8thr1504816928.jpg", Video: "", SrcDate: time.Date(2017, time.September, 7, 21, 42, 9, 0, time.UTC), CreatedAt: time.Date(2023, time.October, 1, 0, 33, 47, 0, time.UTC), UpdatedAt: time.Date(2023, time.October, 1, 0, 33, 47, 0, time.UTC)},
				{ID: 14610, Name: "ACID", Alcoholic: "Alcoholic", Category: "Shot", Ingredients: []entity.Ingredient{{Name: "151 proof rum", Measure: "1 oz Bacardi "}, {Name: "Wild Turkey", Measure: "1 oz "}}, Instructions: "Poor in the 151 first followed by the 101 served with a Coke or Dr Pepper chaser.", Glass: "Shot glass", IBA: "", ImgAttribution: "", ImgSrc: "", Tags: "", Thumb: "https://www.thecocktaildb.com/images/media/drink/xuxpxt1479209317.jpg", Video: "", SrcDate: time.Date(2016, time.November, 15, 11, 28, 37, 0, time.UTC), CreatedAt: time.Date(2023, time.October, 1, 0, 33, 47, 0, time.UTC), UpdatedAt: time.Date(2023, time.October, 1, 0, 33, 47, 0, time.UTC)},
				{ID: 13938, Name: "AT&T", Alcoholic: "Alcoholic", Category: "Ordinary Drink", Ingredients: []entity.Ingredient{{Name: "Absolut Vodka", Measure: "1 oz "}, {Name: "Gin", Measure: "1 oz "}, {Name: "Tonic water", Measure: "4 oz "}}, Instructions: "Pour Vodka and Gin over ice, add Tonic and Stir", Glass: "Highball Glass", IBA: "", ImgAttribution: "", ImgSrc: "", Tags: "", Thumb: "https://www.thecocktaildb.com/images/media/drink/rhhwmp1493067619.jpg", Video: "", SrcDate: time.Date(2017, time.April, 24, 22, 0, 19, 0, time.UTC), CreatedAt: time.Date(2023, time.October, 1, 0, 33, 47, 0, time.UTC), UpdatedAt: time.Date(2023, time.October, 1, 0, 33, 47, 0, time.UTC)},
			},
			err: nil,
			repo: repo{
				args: repoArgs{nType: ct.EvenNum, jobs: 8, jWorker: 2},
				resp: []entity.Cocktail{
					{ID: 17222, Name: "A1", Alcoholic: "Alcoholic", Category: "Cocktail", Ingredients: []entity.Ingredient{{Name: "Gin", Measure: "1 3/4 shot "}, {Name: "Grand Marnier", Measure: "1 Shot "}, {Name: "Lemon Juice", Measure: "1/4 Shot"}, {Name: "Grenadine", Measure: "1/8 Shot"}}, Instructions: "Pour all ingredients into a cocktail shaker, mix and serve over ice into a chilled glass.", Glass: "Cocktail glass", IBA: "", ImgAttribution: "", ImgSrc: "", Tags: "", Thumb: "https://www.thecocktaildb.com/images/media/drink/2x8thr1504816928.jpg", Video: "", SrcDate: time.Date(2017, time.September, 7, 21, 42, 9, 0, time.UTC), CreatedAt: time.Date(2023, time.October, 1, 0, 33, 47, 0, time.UTC), UpdatedAt: time.Date(2023, time.October, 1, 0, 33, 47, 0, time.UTC)},
					{ID: 14610, Name: "ACID", Alcoholic: "Alcoholic", Category: "Shot", Ingredients: []entity.Ingredient{{Name: "151 proof rum", Measure: "1 oz Bacardi "}, {Name: "Wild Turkey", Measure: "1 oz "}}, Instructions: "Poor in the 151 first followed by the 101 served with a Coke or Dr Pepper chaser.", Glass: "Shot glass", IBA: "", ImgAttribution: "", ImgSrc: "", Tags: "", Thumb: "https://www.thecocktaildb.com/images/media/drink/xuxpxt1479209317.jpg", Video: "", SrcDate: time.Date(2016, time.November, 15, 11, 28, 37, 0, time.UTC), CreatedAt: time.Date(2023, time.October, 1, 0, 33, 47, 0, time.UTC), UpdatedAt: time.Date(2023, time.October, 1, 0, 33, 47, 0, time.UTC)},
					{ID: 13938, Name: "AT&T", Alcoholic: "Alcoholic", Category: "Ordinary Drink", Ingredients: []entity.Ingredient{{Name: "Absolut Vodka", Measure: "1 oz "}, {Name: "Gin", Measure: "1 oz "}, {Name: "Tonic water", Measure: "4 oz "}}, Instructions: "Pour Vodka and Gin over ice, add Tonic and Stir", Glass: "Highball Glass", IBA: "", ImgAttribution: "", ImgSrc: "", Tags: "", Thumb: "https://www.thecocktaildb.com/images/media/drink/rhhwmp1493067619.jpg", Video: "", SrcDate: time.Date(2017, time.April, 24, 22, 0, 19, 0, time.UTC), CreatedAt: time.Date(2023, time.October, 1, 0, 33, 47, 0, time.UTC), UpdatedAt: time.Date(2023, time.October, 1, 0, 33, 47, 0, time.UTC)},
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mRepo := mocks.NewCocktailRepo()
			mRepo.On("ReadCC", tt.repo.args.nType, tt.repo.args.jobs, tt.repo.args.jWorker).
				Return(tt.repo.resp, tt.repo.err)
			svc := NewCocktail(mRepo)
			require.NotNil(t, svc)

			out, err := svc.GetCC(tt.args.nType, tt.args.jobs, tt.args.jWorker)
			if tt.err != nil {
				require.NotNil(t, err)
				assert.Nil(t, out)
				assert.ErrorIs(t, err, tt.err)
				return
			}
			require.Nil(t, err)
			require.NotNil(t, out)
			assert.Len(t, tt.exp, len(out))
			assert.Equal(t, tt.exp, out)
		})
	}
}

func TestCocktail_UpdateDB(t *testing.T) {
	type repo struct {
		createArg []entity.Cocktail
		createErr error
		fetchResp []entity.Cocktail
		fetchErr  error
		readResp  []entity.Cocktail
		readErr   error
	}
	tests := []struct {
		name string
		repo repo
		exp  ct.DBOpsSummary
		err  error
	}{
		{
			name: "Read All error",
			repo: repo{
				readResp: nil,
				readErr:  testRepoErr,
			},
			exp: ct.DBOpsSummary{},
			err: testRepoErr,
		},

		{
			name: "Fetch Data error",
			repo: repo{
				readResp:  testCocktailsAll,
				readErr:   nil,
				fetchResp: nil,
				fetchErr:  testRepoErr,
			},
			exp: ct.DBOpsSummary{},
			err: testRepoErr,
		},
		{
			name: "Replace Data error",
			repo: repo{
				readResp:  []entity.Cocktail{},
				readErr:   nil,
				fetchResp: testCocktailsAll,
				fetchErr:  nil,
				createArg: []entity.Cocktail{
					{ID: 1, Name: "Foo", Alcoholic: "Alcoholic", Category: "Foo Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "soda", Measure: "80ml"}}, CreatedAt: dateTimeNow(), UpdatedAt: dateTimeNow()},
					{ID: 2, Name: "Bar", Alcoholic: "Non alcoholic", Category: "Some Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "water", Measure: "50ml"}}, CreatedAt: dateTimeNow(), UpdatedAt: dateTimeNow()},
					{ID: 3, Name: "Baz", Alcoholic: "Alcoholic", Category: "Some Category", Glass: "Cocktail glass", Ingredients: []entity.Ingredient{{Name: "soda", Measure: "100ml"}}, CreatedAt: dateTimeNow(), UpdatedAt: dateTimeNow()},
				},
				createErr: testRepoErr,
			},
			exp: ct.DBOpsSummary{},
			err: testRepoErr,
		},
		{
			name: "No changes",
			repo: repo{
				readResp:  testCocktailsAll,
				readErr:   nil,
				fetchResp: testCocktailsAll,
				fetchErr:  nil,
				createArg: []entity.Cocktail{},
				createErr: nil,
			},
			exp: ct.DBOpsSummary{
				Status:       noChangesDBStatus,
				NewRecs:      0,
				ModifiedRecs: 0,
				TotalOps:     0,
				TotalRecs:    3,
			},
			err: nil,
		},
		{
			name: "Same source date, but different record",
			repo: repo{
				readResp: []entity.Cocktail{
					{ID: 1, Name: "foo"},
					{ID: 2, Name: "bar", Category: "some-category"},
					{ID: 3, Name: "baz"},
				},
				readErr: nil,
				fetchResp: []entity.Cocktail{
					{ID: 1, Name: "foo"},
					{ID: 2, Name: "bar", Category: "other-category"},
					{ID: 3, Name: "baz"},
				},
				fetchErr: nil,
				createArg: []entity.Cocktail{
					{ID: 1, Name: "foo"},
					{ID: 2, Name: "bar", Category: "other-category", UpdatedAt: dateTimeNow()},
					{ID: 3, Name: "baz"},
				},
				createErr: nil,
			},
			exp: ct.DBOpsSummary{
				Status:       successfulUpdateDBStatus,
				NewRecs:      0,
				ModifiedRecs: 1,
				TotalOps:     1,
				TotalRecs:    3,
			},
			err: nil,
		},
		{
			name: "One updated",
			repo: repo{
				readResp: []entity.Cocktail{
					{ID: 1, Name: "foo", SrcDate: dateTimeNow()},
					{ID: 2, Name: "bar"},
					{ID: 3, Name: "baz"},
				},
				readErr: nil,
				fetchResp: []entity.Cocktail{
					{ID: 1, Name: "foo", Category: "fooCategory", SrcDate: dateTimeNow().Add(1 * time.Hour)},
					{ID: 2, Name: "bar"},
					{ID: 3, Name: "baz"},
				},
				fetchErr: nil,
				createArg: []entity.Cocktail{
					{ID: 1, Name: "foo", Category: "fooCategory", SrcDate: dateTimeNow().Add(1 * time.Hour), UpdatedAt: dateTimeNow()},
					{ID: 2, Name: "bar"},
					{ID: 3, Name: "baz"},
				},
				createErr: nil,
			},
			exp: ct.DBOpsSummary{
				Status:       successfulUpdateDBStatus,
				NewRecs:      0,
				ModifiedRecs: 1,
				TotalOps:     1,
				TotalRecs:    3,
			},
			err: nil,
		},
		{
			name: "Two new",
			repo: repo{
				readResp: []entity.Cocktail{
					{ID: 1, Name: "foo"},
				},
				readErr: nil,
				fetchResp: []entity.Cocktail{
					{ID: 1, Name: "foo"},
					{ID: 2, Name: "bar"},
					{ID: 3, Name: "baz"},
				},
				fetchErr: nil,
				createArg: []entity.Cocktail{
					{ID: 1, Name: "foo"},
					{ID: 2, Name: "bar", CreatedAt: dateTimeNow(), UpdatedAt: dateTimeNow()},
					{ID: 3, Name: "baz", CreatedAt: dateTimeNow(), UpdatedAt: dateTimeNow()},
				},
				createErr: nil,
			},
			exp: ct.DBOpsSummary{
				Status:       successfulUpdateDBStatus,
				NewRecs:      2,
				ModifiedRecs: 0,
				TotalOps:     2,
				TotalRecs:    3,
			},
			err: nil,
		},
		{
			name: "All new",
			repo: repo{
				readResp: []entity.Cocktail{},
				readErr:  nil,
				fetchResp: []entity.Cocktail{
					{ID: 1, Name: "foo"},
					{ID: 2, Name: "bar"},
					{ID: 3, Name: "baz"},
				},
				fetchErr: nil,
				createArg: []entity.Cocktail{
					{ID: 1, Name: "foo", CreatedAt: dateTimeNow(), UpdatedAt: dateTimeNow()},
					{ID: 2, Name: "bar", CreatedAt: dateTimeNow(), UpdatedAt: dateTimeNow()},
					{ID: 3, Name: "baz", CreatedAt: dateTimeNow(), UpdatedAt: dateTimeNow()},
				},
				createErr: nil,
			},
			exp: ct.DBOpsSummary{
				Status:       successfulUpdateDBStatus,
				NewRecs:      3,
				ModifiedRecs: 0,
				TotalOps:     3,
				TotalRecs:    3,
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mRepo := mocks.NewCocktailRepo()
			mRepo.On("ReadAll").Return(tt.repo.readResp, tt.repo.readErr)
			mRepo.On("ReplaceDB", tt.repo.createArg).Return(tt.repo.createErr)
			mRepo.On("Fetch").Return(tt.repo.fetchResp, tt.repo.fetchErr)
			svc := NewCocktail(mRepo)
			require.NotNil(t, svc)

			out, err := svc.UpdateDB()
			if tt.err != nil {
				require.NotNil(t, err)
				assert.Equal(t, ct.DBOpsSummary{}, out)
				assert.ErrorIs(t, err, tt.err)
				return
			}
			require.Nil(t, err)
			assert.Equal(t, tt.exp.Status, out.Status)
			assert.Equal(t, tt.exp.NewRecs, out.NewRecs)
			assert.Equal(t, tt.exp.ModifiedRecs, out.ModifiedRecs)
			assert.Equal(t, tt.exp.TotalOps, out.TotalOps)
			assert.Equal(t, tt.exp.TotalRecs, out.TotalRecs)
		})
	}
}
