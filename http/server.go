package http

import (
	"log"
	"net/http"
	"time"

	"github.com/sharkyze/lbc/metrics"
)

const (
	httpServerReadTimeout  = 5 * time.Second
	httpServerWriteTimeout = 10 * time.Second
	httpServerIdleTimeout  = 120 * time.Second
)

// New returns a new http.Server, the caller of this function is responsible for
// starting and gracefully shutting down the server.
func New(port string, logger *log.Logger, metrics metrics.Metrics) http.Server {
	hs := newHandlers(logger, metrics)

	mux := http.NewServeMux()
	mux.HandleFunc("/fizzbuzz", hs.handleFizzBuzz)
	mux.HandleFunc("/metrics", hs.handleMetrics)

	return http.Server{
		Addr:         port,                           // configure the bind address
		Handler:      loggingMiddleware(logger, mux), // set the default handler
		ReadTimeout:  httpServerReadTimeout,          // max time to read request from the client
		WriteTimeout: httpServerWriteTimeout,         // max time to write response to the client
		IdleTimeout:  httpServerIdleTimeout,          // max time for connections using TCP Keep-Alive
	}
}
