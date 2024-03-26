package sharedhttp

import (
	"github.com/marcos-wz/capstone-go-bootcamp/internal/config"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/controller"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Chi configures a chi.Mux instance.
type Chi struct {
	basePath    string
	router      *chi.Mux
	controllers map[string]controller.HTTP
}

// NewChi returns a Chi implementation.
// It allocates a pre-configured chi.Mux instance.
func NewChi(cfg config.Application) Chi {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	return Chi{
		basePath: cfg.BasePath(),
		router:   r,
	}
}

// Add appends the controller.HTTP given to the controller map list.
func (m *Chi) Add(name string, ctrl controller.HTTP) {
	if m.controllers == nil {
		m.controllers = make(map[string]controller.HTTP)
	}
	m.controllers[name] = ctrl
}

// RegisterRoutes register the routes defined on the controller map list.
// All controller's routes are prefixed with the configured base path. e.g. "/api/v0"
func (m *Chi) RegisterRoutes() {
	if len(m.controllers) == 0 {
		return
	}
	m.router.Route(m.basePath, func(r chi.Router) {
		for name, ctrl := range m.controllers {
			if ctrl != nil {
				ctrl.SetRoutes(r)
				logger.Log().Debug().
					Str("controller", name).
					Msg("registered http controller")
			}
		}
	})
}

// Router returns the configured chi.Mux instance.
func (m *Chi) Router() *chi.Mux {
	return m.router
}
