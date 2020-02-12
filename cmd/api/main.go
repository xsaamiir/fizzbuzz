package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	lbchttp "github.com/sharkyze/lbc/http"
	"github.com/sharkyze/lbc/metrics"
)

const (
	httpServerPort = ":8000"

	shutdownWait = time.Second * 15
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", err)
		os.Exit(1) // nolint: gomnd
	}
}

func run() error {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	m := metrics.NewInMemoryMetrics()

	srv := lbchttp.New(httpServerPort, logger, &m)

	// start the server in a goroutine so that we can continue
	// listening to events in the main goroutine.
	go func() {
		logger.Println("Starting server on port " + httpServerPort)

		if err := srv.ListenAndServe(); err != nil {
			logger.Fatalf("Error starting server: %s\n", err)
		}
	}()

	backgroundCtx := context.Background()

	// Check for a closing signal for graceful shutdown
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)

	// This will block until a closing signal is received to exit
	sig := <-sigquit

	logger.Println("🛑 caught sig: " + sig.String())
	logger.Println("👋 starting graceful server shutdown")

	// Create a deadline to use for server shutdown.
	srvShutdownCtx, srvShutdownCtxCancel := context.WithTimeout(backgroundCtx, shutdownWait)
	defer srvShutdownCtxCancel()

	// Doesn't block if there are no open connections to the server,
	// but will otherwise wait until the timeout deadline.
	if err := srv.Shutdown(srvShutdownCtx); err != nil {
		return fmt.Errorf("⚠️ unable to shut down server: %w", err)
	}

	logger.Println("✅ server shutdown gracefully")

	return nil
}
