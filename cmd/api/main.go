package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/sharkyze/lbc/fizzbuzz"
)

const (
	shutdownWait = time.Second * 15

	httpServerPort         = ":8000"
	httpServerReadTimeout  = 5 * time.Second
	httpServerWriteTimeout = 10 * time.Second
	httpServerIdleTimeout  = 120 * time.Second
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", err)
		os.Exit(1) // nolint: gomnd
	}
}

func run() error {
	hs := newHandlers()

	mux := http.NewServeMux()
	mux.HandleFunc("/fizzbuzz", hs.handleFizzBuzz)
	mux.HandleFunc("/metrics", hs.handleMetrics)

	// create a new server
	srv := http.Server{
		Addr:         httpServerPort,         // configure the bind address
		Handler:      mux,                    // set the default handler
		ReadTimeout:  httpServerReadTimeout,  // max time to read request from the client
		WriteTimeout: httpServerWriteTimeout, // max time to write response to the client
		IdleTimeout:  httpServerIdleTimeout,  // max time for connections using TCP Keep-Alive
	}

	// start the server in a goroutine so that we can continue
	// listening to events in the main goroutine.
	go func() {
		log.Println("Starting server on port " + httpServerPort)

		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	backgroundCtx := context.Background()

	// Check for a closing signal for graceful shutdown
	sigquit := make(chan os.Signal, 1)
	signal.Notify(sigquit, os.Interrupt, syscall.SIGTERM)

	// This will block until a closing signal is received to exit
	sig := <-sigquit

	log.Println("🛑 caught sig: " + sig.String())
	log.Println("👋 starting graceful server shutdown")

	// Create a deadline to use for server shutdown.
	srvShutdownCtx, srvShutdownCtxCancel := context.WithTimeout(backgroundCtx, shutdownWait)
	defer srvShutdownCtxCancel()

	// Doesn't block if there are no open connections to the server,
	// but will otherwise wait until the timeout deadline.
	if err := srv.Shutdown(srvShutdownCtx); err != nil {
		return fmt.Errorf("⚠️ unable to shut down server: %w", err)
	}

	log.Println("✅ server shutdown gracefully")

	return nil
}

type fizzBuzzRequest struct {
	Int1, Int2, Limit int
	Str1, Str2        string
}

type handlers struct {
	mu sync.Mutex
	v  map[fizzBuzzRequest]int
}

func newHandlers() handlers {
	return handlers{v: make(map[fizzBuzzRequest]int)}
}

func (h *handlers) count(request fizzBuzzRequest) {
	h.mu.Lock()
	h.v[request]++
	h.mu.Unlock()
}

func (h *handlers) handleFizzBuzz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only the http method GET is accepted", http.StatusMethodNotAllowed)
		return
	}

	int1 := r.URL.Query().Get("int1")

	n1, err := strconv.Atoi(int1)
	if err != nil {
		http.Error(w, "error parsing int1", http.StatusBadRequest)
		return
	}

	int2 := r.URL.Query().Get("int2")

	n2, err := strconv.Atoi(int2)
	if err != nil {
		http.Error(w, "error parsing int2", http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")

	l, err := strconv.Atoi(limit)
	if err != nil {
		http.Error(w, "error parsing limit", http.StatusBadRequest)
		return
	}

	str1 := r.URL.Query().Get("str1")
	if str1 == "" {
		http.Error(w, "missing required parameter str1", http.StatusBadRequest)
		return
	}

	str2 := r.URL.Query().Get("str2")
	if str2 == "" {
		http.Error(w, "missing required parameter str2", http.StatusBadRequest)
		return
	}

	h.count(fizzBuzzRequest{
		Int1:  n1,
		Int2:  n2,
		Limit: l,
		Str1:  str1,
		Str2:  str2,
	})

	res := fizzbuzz.FizzBuzz(n1, n2, l, str1, str2)

	w.Header().Add("content-type", "application/json; charset=utf-8")
	// https://stackoverflow.com/questions/33903552/what-input-will-cause-golangs-json-marshal-to-return-an-error
	_ = json.NewEncoder(w).Encode(res)
}

func (h *handlers) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only the http method GET is accepted", http.StatusMethodNotAllowed)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	type metricsResult struct {
		Request fizzBuzzRequest `json:"request"`
		Hits    int             `json:"hits"`
	}

	res := make([]metricsResult, len(h.v))

	var idx int

	for request, hits := range h.v {
		res[idx] = metricsResult{
			Request: request,
			Hits:    hits,
		}
		idx++
	}

	w.Header().Add("content-type", "application/json; charset=utf-8")
	// https://stackoverflow.com/questions/33903552/what-input-will-cause-golangs-json-marshal-to-return-an-error
	_ = json.NewEncoder(w).Encode(res)
}
