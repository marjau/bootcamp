package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/marcos-wz/capstone-go-bootcamp/internal/controller/mocks"
	ct "github.com/marcos-wz/capstone-go-bootcamp/internal/customtype"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/entity"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/repository"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ CocktailSvc = &mocks.CocktailSvc{}

func TestNewCocktail(t *testing.T) {
	mSvc := mocks.NewCocktailSvc()
	require.NotNil(t, mSvc)
	out := NewCocktail(mSvc)
	assert.IsType(t, Cocktail{}, out)
}

func TestCocktail_GetFiltered(t *testing.T) {
	type params struct {
		filter string
		value  string
	}
	type svc struct {
		resp []entity.Cocktail
		err  error
	}
	tests := []struct {
		name    string
		params  params
		code    int
		err     errHTTP
		svc     svc
		wantErr bool
	}{
		{
			name:   "Repository CSV error",
			params: params{filter: "id", value: "50"},
			code:   http.StatusInternalServerError,
			err: errHTTP{
				Code:      http.StatusInternalServerError,
				ErrorType: repoCsvErrType,
				Message:   "csv: %!s(<nil>)",
			},
			svc: svc{
				resp: nil,
				err:  &repository.CsvErr{},
			},
			wantErr: true,
		},
		{
			name:   "Service error",
			params: params{filter: "id", value: "50"},
			code:   http.StatusBadRequest,
			err: errHTTP{
				Code:      http.StatusBadRequest,
				ErrorType: errType(reflect.TypeOf(testSvcErr).String()),
				Message:   testSvcErr.Error(),
			},
			svc: svc{
				resp: nil,
				err:  testSvcErr,
			},
			wantErr: true,
		},
		{
			name:    "Empty",
			params:  params{filter: "", value: ""},
			code:    http.StatusNotFound,
			err:     errHTTP{},
			svc:     svc{},
			wantErr: true,
		},
		{
			name:   "Arbitrary",
			params: params{filter: "foo", value: "some-value"},
			code:   http.StatusUnprocessableEntity,
			err: errHTTP{
				Code:      http.StatusUnprocessableEntity,
				ErrorType: svcFilterErrType,
				Message:   "service filter: invalid filter",
			},
			svc: svc{
				resp: nil,
				err:  &service.FilterErr{Err: service.ErrFltrInvalid},
			},
			wantErr: true,
		},
		{
			name:   "Bad Value",
			params: params{filter: "id", value: "asd"},
			code:   http.StatusUnprocessableEntity,
			svc: svc{
				resp: nil,
				err:  &service.FilterErr{},
			},
			wantErr: true,
		},
		{
			name:   "Not Found",
			params: params{filter: "id", value: "123456"},
			code:   http.StatusOK,
			svc: svc{
				resp: []entity.Cocktail{},
				err:  nil,
			},
			wantErr: false,
		},
		{
			name:   "Valid",
			params: params{filter: "id", value: "2"},
			code:   http.StatusOK,
			svc: svc{
				resp: []entity.Cocktail{
					{ID: 2, Name: "Bar", Alcoholic: "Non alcoholic", Category: "Some Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "water", Measure: "50ml"}}},
				},
				err: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mSvc := mocks.NewCocktailSvc()
			mSvc.On("GetFiltered", tt.params.filter, tt.params.value).
				Return(tt.svc.resp, tt.svc.err)
			ctrl := Cocktail{svc: mSvc}

			// Request
			path := fmt.Sprintf("/cocktail/%v/%v", tt.params.filter, tt.params.value)
			req, err := http.NewRequest("GET", path, nil)
			require.Nil(t, err)

			// Server instance
			rr := httptest.NewRecorder()
			srv := newTestRouter(ctrl)
			srv.ServeHTTP(rr, req)

			// Tests
			assert.Equal(t, tt.code, rr.Code)
			if tt.wantErr {
				if tt.err != (errHTTP{}) {
					var errMsg errHTTP
					require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &errMsg))
					assert.Equal(t, tt.err, errMsg)
				}
				return
			}

			var resp []entity.Cocktail
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			assert.Len(t, tt.svc.resp, len(resp))
			assert.Equal(t, tt.svc.resp, resp)
		})
	}
}

func TestCocktail_GetAll(t *testing.T) {
	type svc struct {
		resp []entity.Cocktail
		err  error
	}

	tests := []struct {
		name    string
		code    int
		err     errHTTP
		svc     svc
		wantErr bool
	}{
		{
			name: "Repository CSV error",
			code: http.StatusInternalServerError,
			err: errHTTP{
				Code:      http.StatusInternalServerError,
				ErrorType: repoCsvErrType,
				Message:   "csv: %!s(<nil>)",
			},
			svc: svc{
				resp: nil,
				err:  &repository.CsvErr{},
			},
			wantErr: true,
		},
		{
			name: "Service error",
			code: http.StatusBadRequest,
			err: errHTTP{
				Code:      http.StatusBadRequest,
				ErrorType: errType(reflect.TypeOf(testSvcErr).String()),
				Message:   testSvcErr.Error(),
			},
			svc: svc{
				resp: nil,
				err:  testSvcErr,
			},
			wantErr: true,
		},
		{
			name: "Not records",
			code: http.StatusOK,
			svc: svc{
				resp: []entity.Cocktail{},
				err:  nil,
			},
			wantErr: false,
		},
		{
			name: "All Records",
			code: http.StatusOK,
			svc: svc{
				resp: []entity.Cocktail{
					{ID: 1, Name: "Foo", Alcoholic: "Alcoholic", Category: "Foo Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "soda", Measure: "80ml"}}},
					{ID: 2, Name: "Bar", Alcoholic: "Non alcoholic", Category: "Some Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "water", Measure: "50ml"}}},
					{ID: 3, Name: "Baz", Alcoholic: "Alcoholic", Category: "Some Category", Glass: "Cocktail glass", Ingredients: []entity.Ingredient{{Name: "soda", Measure: "100ml"}}},
				},
				err: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mSvc := mocks.NewCocktailSvc()
			mSvc.On("GetAll").Return(tt.svc.resp, tt.svc.err)
			ctrl := Cocktail{svc: mSvc}

			// Request
			req, err := http.NewRequest("GET", "/cocktails", nil)
			require.Nil(t, err)

			// Server instance
			rr := httptest.NewRecorder()
			srv := newTestRouter(ctrl)
			srv.ServeHTTP(rr, req)

			// Tests
			assert.Equal(t, tt.code, rr.Code)
			if tt.wantErr {
				var errMsg errHTTP
				require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &errMsg))
				assert.Equal(t, tt.err, errMsg)
				return
			}

			var resp []entity.Cocktail
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			assert.Len(t, tt.svc.resp, len(resp))
			assert.Equal(t, tt.svc.resp, resp)
		})
	}
}

func TestCocktail_GetCC(t *testing.T) {
	type svc struct {
		resp []entity.Cocktail
		err  error
	}
	type args struct {
		nType       string
		items       string
		itemsWorker string
	}
	tests := []struct {
		name    string
		args    args
		code    int
		err     errHTTP
		svc     svc
		wantErr bool
	}{
		{
			name: "Repository CSV error",
			args: args{
				nType: "even", items: "10", itemsWorker: "4",
			},
			code: http.StatusInternalServerError,
			err: errHTTP{
				Code:      http.StatusInternalServerError,
				ErrorType: repoCsvErrType,
				Message:   "csv: %!s(<nil>)",
			},
			svc: svc{
				resp: nil,
				err:  &repository.CsvErr{},
			},
			wantErr: true,
		},
		{
			name: "Service Args error",
			args: args{nType: "foo", items: "10", itemsWorker: "4"},
			code: http.StatusUnprocessableEntity,
			err: errHTTP{
				Code:      http.StatusUnprocessableEntity,
				ErrorType: svcArgsErrType,
				Message:   "service arguments: invalid number type",
			},
			svc: svc{
				resp: nil,
				err:  &service.ArgsErr{Err: service.ErrInvalidNumType},
			},
			wantErr: true,
		},
		{
			name: "Not records",
			args: args{nType: "even", items: "10", itemsWorker: "4"},
			code: http.StatusOK,
			svc: svc{
				resp: []entity.Cocktail{},
				err:  nil,
			},
			wantErr: false,
		},
		{
			name: "Even",
			args: args{nType: "even", items: "10", itemsWorker: "4"},
			code: http.StatusOK,
			svc: svc{
				resp: []entity.Cocktail{
					{ID: 2, Name: "Bar", Alcoholic: "Non alcoholic", Category: "Some Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "water", Measure: "50ml"}}},
				},
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "Odd",
			args: args{nType: "even", items: "10", itemsWorker: "4"},
			code: http.StatusOK,
			svc: svc{
				resp: []entity.Cocktail{
					{ID: 1, Name: "Foo", Alcoholic: "Alcoholic", Category: "Foo Category", Glass: "Shot glass", Ingredients: []entity.Ingredient{{Name: "soda", Measure: "80ml"}}},
					{ID: 3, Name: "Baz", Alcoholic: "Alcoholic", Category: "Some Category", Glass: "Cocktail glass", Ingredients: []entity.Ingredient{{Name: "soda", Measure: "100ml"}}},
				},
				err: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mSvc := mocks.NewCocktailSvc()
			mSvc.On("GetCC", tt.args.nType, tt.args.items, tt.args.itemsWorker).
				Return(tt.svc.resp, tt.svc.err)
			ctrl := Cocktail{svc: mSvc}

			// Request
			path := fmt.Sprintf("/cocktails/%v/%v/%v", tt.args.nType, tt.args.items, tt.args.itemsWorker)
			req, err := http.NewRequest("GET", path, nil)
			require.Nil(t, err)

			// Server instance
			rr := httptest.NewRecorder()
			srv := newTestRouter(ctrl)
			srv.ServeHTTP(rr, req)

			// Tests
			assert.Equal(t, tt.code, rr.Code)
			if tt.wantErr {
				var errMsg errHTTP
				require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &errMsg))
				assert.Equal(t, tt.err, errMsg)
				return
			}

			var resp []entity.Cocktail
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			assert.Len(t, tt.svc.resp, len(resp))
			assert.Equal(t, tt.svc.resp, resp)
		})
	}
}

func TestCocktail_UpdateDB(t *testing.T) {
	type svc struct {
		summary ct.DBOpsSummary
		err     error
	}
	tests := []struct {
		name    string
		code    int
		err     errHTTP
		svc     svc
		wantErr bool
	}{
		{
			name: "CSV error",
			code: http.StatusInternalServerError,
			err: errHTTP{
				Code:      http.StatusInternalServerError,
				ErrorType: repoCsvErrType,
				Message:   "csv: %!s(<nil>)",
			},
			svc: svc{
				summary: ct.DBOpsSummary{},
				err:     &repository.CsvErr{},
			},
			wantErr: true,
		},
		{
			name: "Data API error",
			code: http.StatusBadGateway,
			err: errHTTP{
				Code:      http.StatusBadGateway,
				ErrorType: repoDataApiErrType,
				Message:   "data api: %!s(<nil>)",
			},
			svc: svc{
				summary: ct.DBOpsSummary{},
				err:     &repository.DataApiErr{},
			},
			wantErr: true,
		},
		{
			name: "Service error",
			code: http.StatusBadRequest,
			err: errHTTP{
				Code:      http.StatusBadRequest,
				ErrorType: errType(reflect.TypeOf(testSvcErr).String()),
				Message:   testSvcErr.Error(),
			},
			svc: svc{
				summary: ct.DBOpsSummary{},
				err:     testSvcErr,
			},
			wantErr: true,
		},
		{
			name: "Success",
			code: http.StatusOK,
			err:  errHTTP{},
			svc: svc{
				summary: ct.DBOpsSummary{
					Status:       "some status",
					NewRecs:      5,
					ModifiedRecs: 5,
					TotalOps:     10,
				},
				err: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mSvc := mocks.NewCocktailSvc()
			mSvc.On("UpdateDB").Return(tt.svc.summary, tt.svc.err)
			ctrl := Cocktail{mSvc}

			// Request
			req, err := http.NewRequest("GET", "/cocktail/updatedb", nil)
			require.Nil(t, err)

			// Server instance
			rr := httptest.NewRecorder()
			srv := newTestRouter(ctrl)
			srv.ServeHTTP(rr, req)

			// Tests
			assert.Equal(t, tt.code, rr.Code)
			if tt.wantErr {
				var errMsg errHTTP
				require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &errMsg))
				assert.Equal(t, tt.err, errMsg)
				return
			}

			var resp ct.DBOpsSummary
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
			assert.Equal(t, tt.svc.summary, resp)
		})
	}
}
