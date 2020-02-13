package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/sharkyze/lbc/fizzbuzz"
	"github.com/sharkyze/lbc/metrics"
)

type (
	handlers struct {
		logger *log.Logger
		metrics.FizzBuzzMetrics
	}
)

func newHandlers(logger *log.Logger, metrics metrics.FizzBuzzMetrics) handlers {
	return handlers{
		logger:          logger,
		FizzBuzzMetrics: metrics,
	}
}

func (h *handlers) handleFizzBuzz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only the http method GET is accepted", http.StatusMethodNotAllowed)
		return
	}

	int1 := r.URL.Query().Get("int1")

	n1, err := strconv.Atoi(int1)
	if err != nil {
		h.respond(w, r, "error parsing int1", nil, http.StatusBadRequest)
		return
	}

	int2 := r.URL.Query().Get("int2")

	n2, err := strconv.Atoi(int2)
	if err != nil {
		h.respond(w, r, "error parsing int2", nil, http.StatusBadRequest)
		return
	}

	limit := r.URL.Query().Get("limit")

	l, err := strconv.Atoi(limit)
	if err != nil {
		h.respond(w, r, "error parsing limit", nil, http.StatusBadRequest)
		return
	}

	str1 := r.URL.Query().Get("str1")
	if str1 == "" {
		h.respond(w, r, "missing required parameter str1", nil, http.StatusBadRequest)
		return
	}

	str2 := r.URL.Query().Get("str2")
	if str2 == "" {
		h.respond(w, r, "missing required parameter str2", nil, http.StatusBadRequest)
		return
	}

	h.Count(metrics.FizzBuzzRequest{
		Int1:  n1,
		Int2:  n2,
		Limit: l,
		Str1:  str1,
		Str2:  str2,
	})

	res := fizzbuzz.FizzBuzz(n1, n2, l, str1, str2)
	h.respond(w, r, "", res, http.StatusOK)
}

func (h *handlers) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "only the http method GET is accepted", http.StatusMethodNotAllowed)
		return
	}

	h.respond(w, r, "", h.FizzBuzzMetrics.Get(), http.StatusOK)
}
