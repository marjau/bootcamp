package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

var _ HTTP = &Home{}

// Home set the routes and handler functions related to the Home controller.
type Home struct{}

// NewHome returns a new Home implementation.
func NewHome() Home {
	return Home{}
}

// SetRoutes sets a fresh middleware stack for the Home handle functions and mounts them to the provided sub router.
func (h Home) SetRoutes(r chi.Router) {
	r.Get("/", h.HomePage)
}

// HomePage is a handler function that responses the home page in HTML format.
func (Home) HomePage(w http.ResponseWriter, r *http.Request) {
	render.HTML(w, r, `
		<H1>Welcome to the Capstone API</H1>
		<H2>GO Bootcamp - Academy </H2>
		<B>Gopher:</B> Marcos Jauregui <BR>
		<B>Email:</B> marcos.jauregui@wizeline.com
	`)
}
