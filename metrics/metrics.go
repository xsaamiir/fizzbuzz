package metrics

import (
	"errors"
	"sync"
)

type (
	// Result represents the metrics result of FizzBuzz requests.
	Result struct {
		Request Request `json:"request"`
		Hits    int     `json:"hits"`
	}

	// Result is the input needed to do a "FizzBuzz"
	Request struct {
		Int1  int    `json:"int1"`
		Int2  int    `json:"int2"`
		Limit int    `json:"limit"`
		Str1  string `json:"str1"`
		Str2  string `json:"str2"`
	}

	// Metrics is the interface to record and retrieve service usage
	// metrics.
	// It wraps a:
	// 	- Record method that takes a request and adds it to existing
	// 	recorded metrics.
	// 	- Get method to retrieve the service's metrics.
	Metrics interface {
		Record(Request)
		Get() []Result
	}

	// InMemoryMetrics is an implementation of Metrics in memory.
	InMemoryMetrics struct {
		mu sync.Mutex
		v  map[Request]int
	}
)

func NewInMemoryMetrics() InMemoryMetrics {
	return InMemoryMetrics{v: make(map[Request]int)}
}

func (m *InMemoryMetrics) Record(request Request) {
	m.mu.Lock()
	m.v[request]++
	m.mu.Unlock()
}

func (m *InMemoryMetrics) Get() []Result {
	m.mu.Lock()
	// https://rakyll.org/inlined-defers/
	defer m.mu.Unlock()

	res := make([]Result, len(m.v))

	var idx int

	for request, hits := range m.v {
		res[idx] = Result{
			Request: request,
			Hits:    hits,
		}
		idx++
	}

	return res
}

func TopHit(m Metrics) (Result, error) {
	hits := m.Get()
	if len(hits) == 0 {
		return Result{}, errors.New("no recorded metrics")
	}

	var top Result
	for _, res := range m.Get() {
		if res.Hits > top.Hits {
			top = res
		}
	}

	return top, nil
}
