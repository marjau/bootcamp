package controller

import (
	"errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// testSvcErr represents a non-custom error of the service
var testSvcErr = errors.New("test service error")

// newTestRouter returns a new pre-allocated chi.Mux instance with the registered routes provided by the given controller.
func newTestRouter(ctrl HTTP) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	ctrl.SetRoutes(r)
	return r
}
