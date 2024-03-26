package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/marcos-wz/capstone-go-bootcamp/internal/config"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/controller"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/logger"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/repository"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/service"
	"github.com/marcos-wz/capstone-go-bootcamp/internal/sharedhttp"
)

// ApiHTTP represents an HTTP API with http.Server and chi.Mux support.
type ApiHTTP struct {
	cfg    config.Config
	server *http.Server
}

// NewApiHTTP returns a new ApiHTTP implementation.
// It uses the configuration instance to set the components.
// It provides the controllers, repository and service dependencies.
// Moreover, implements a http.Server instance with chi.Mux router support.
// Returns an error if any implementation of the dependencies fails.
func NewApiHTTP(cfg config.Config) (ApiHTTP, func(), error) {
	// Cocktail dependencies
	cRepo, err := repository.NewCocktail(cfg)
	if err != nil {
		return ApiHTTP{}, nil, err
	}
	cSvc := service.NewCocktail(cRepo)

	// Router
	router := sharedhttp.NewChi(cfg.Application)
	router.Add("HealthCheck", controller.NewHealthCheck())
	router.Add("Home", controller.NewHome())
	router.Add("Cocktail", controller.NewCocktail(cSvc))
	router.RegisterRoutes()

	return ApiHTTP{
		cfg:    cfg,
		server: sharedhttp.NewHTTPServer(cfg.HTTP.Server, router.Router()),
	}, nil, nil
}

// Start runs the http API and quits doing a grateful shutdown.
// To stop the server you must send a syscall.SIGINT signal usually through `CTRL+C`.
func (h ApiHTTP) Start() {
	go func() {
		logger.Log().Info().Msgf("running http server on %v", h.cfg.HTTP.Server.Address())
		err := h.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log().Fatal().Err(err).Msg("http server startup failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	h.shutdownApi()
}

// shutdownApi performs tasks of safely shutting down processes and closing connections.
func (h ApiHTTP) shutdownApi() {
	ctx, cancel := context.WithTimeout(context.Background(), h.cfg.HTTP.Server.ShutdownTimeout())
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		logger.Log().Error().Err(err).Msg("http server graceful shutdown failed")
	}

	logger.Log().Info().Msg("http api shutdown gracefully")
	os.Exit(0)
}
