package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/sharkyze/lbc/fizzbuzz"
	"github.com/sharkyze/lbc/metrics"
)

type (
	FizzBuzzRequest struct {
		Int1, Int2, Limit int
		Str1, Str2        string
	}

	handlers struct {
		logger *log.Logger
		metrics.FizzBuzzMetrics
	}

	FizzBuzzMetrics struct {
		mu sync.Mutex
		v  map[FizzBuzzRequest]int
	}

	FizzBuzzMetricsResult struct {
		Request FizzBuzzRequest `json:"request"`
		Hits    int             `json:"hits"`
	}
)

func newHandlers(logger *log.Logger, metrics metrics.FizzBuzzMetrics) handlers {
	return handlers{
		logger:          logger,
		FizzBuzzMetrics: metrics,
	}
}

func (m *FizzBuzzMetrics) Count(request FizzBuzzRequest) {
	m.mu.Lock()
	m.v[request]++
	m.mu.Unlock()
}

func (m *FizzBuzzMetrics) Get() []FizzBuzzMetricsResult {
	m.mu.Lock()
	defer m.mu.Unlock()

	res := make([]FizzBuzzMetricsResult, len(m.v))

	var idx int

	for request, hits := range m.v {
		res[idx] = FizzBuzzMetricsResult{
			Request: request,
			Hits:    hits,
		}
		idx++
	}

	return res
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

	h.Count(metrics.FizzBuzzRequest{
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

	res := h.FizzBuzzMetrics.Get()

	w.Header().Add("content-type", "application/json; charset=utf-8")
	// https://stackoverflow.com/questions/33903552/what-input-will-cause-golangs-json-marshal-to-return-an-error
	_ = json.NewEncoder(w).Encode(res)
}
