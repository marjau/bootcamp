package mocks

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// HttpClient is a mock type for the HttpClient dependency
type HttpClient struct {
	mock.Mock
}

// Do provide a mock function with given fields:
func (o *HttpClient) Do(req *http.Request) (*http.Response, error) {
	args := o.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// NewHttpClient creates a new instance of the HttpClient of type Mock.
func NewHttpClient() *HttpClient {
	return new(HttpClient)
}
