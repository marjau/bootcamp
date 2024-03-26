package config

import (
	"fmt"
	"time"
)

// HTTP holds the configurations regarding HTTP items.
type HTTP struct {
	Server  HttpServer
	DataAPI DataAPI
}

// HttpServer is the config used to set the http server instance.
type HttpServer struct {
	host            string
	port            int
	shutdownTimeout time.Duration
}

// Address returns the TCP address for the server to listen on, in the form "host:port".
func (h HttpServer) Address() string {
	return fmt.Sprintf("%v:%d", h.host, h.port)
}

// ShutdownTimeout returns the timeout duration for shutting the server down.
func (h HttpServer) ShutdownTimeout() time.Duration {
	return h.shutdownTimeout
}

// DataAPI represents a data API configuration.
type DataAPI struct {
	url string
}

// NewDataAPI returns a new DataAPI implementation
func NewDataAPI(url string) DataAPI {
	return DataAPI{
		url: url,
	}
}

// URL returns the configured DataAPI url.
func (a DataAPI) URL() string {
	return a.url
}
