package metrics

import (
	"sync"
)

type (
	FizzBuzzMetricsResult struct {
		Request FizzBuzzRequest `json:"request"`
		Hits    int             `json:"hits"`
	}

	FizzBuzzRequest struct {
		Int1, Int2, Limit int
		Str1, Str2        string
	}

	FizzBuzzMetrics interface {
		Count(FizzBuzzRequest)
		Get() []FizzBuzzMetricsResult
	}

	InMemoryMetrics struct {
		mu sync.Mutex
		v  map[FizzBuzzRequest]int
	}
)

func NewInMemoryMetrics() InMemoryMetrics {
	return InMemoryMetrics{v: make(map[FizzBuzzRequest]int)}
}

func (m *InMemoryMetrics) Count(request FizzBuzzRequest) {
	m.mu.Lock()
	m.v[request]++
	m.mu.Unlock()
}

func (m *InMemoryMetrics) Get() []FizzBuzzMetricsResult {
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
