![Lint & Test](https://github.com/sharkyze/lbc/workflows/Lint%20&%20Test/badge.svg)

# LBC

The original fizz-buzz consists in writing all numbers from 1 to 100, and just replacing all multiples of 3 by "fizz", all multiples of 5 by "buzz", and all multiples of 15 by "fizzbuzz". The output would look like this: "1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,...".

Your goal is to implement a web server that will expose a REST API endpoint that:

- Accepts five parameters : three integers int1, int2 and limit, and two strings str1 and str2.
- Returns a list of strings with numbers from 1 to limit, where: all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2.

The server needs to be:

- Ready for production
- Easy to maintain by other developers

Bonus question :

Add a statistics endpoint allowing users to know what the most frequent request has been. This endpoint should:

- Accept no parameter
- Return the parameters corresponding to the most used request, as well as the number of hits for this request

# Getting started

## Running the application

Make sure you have Go 1.13 or a higher version installed.

You can use the following command to start the http server, the server starts on port 8000 by default.

```bash
go get github.com/sharkyze/lbc/...
go run github.com/sharkyze/lbc/cmd/api
```

Once your server is up and running you can try the two API routes:

- /fizzbuzz
- /metrics

```bash
curl 'localhost:8000/fizzbuzz?int1=3&int2=5&limit=20&str1=Fizz&str2=Buzz'
curl 'localhost:8000/fizzbuzz?int1=3&int2=5&limit=20&str1=Fizz&str2=Buzz'
curl 'localhost:8000/metrics'
```

## Running tests and the linter

Both the tests and linter checks run on the CI.
You can also run them locally using the following commands:

```bash
go test -race -v ./...
# If you have [golangci-lint](https://github.com/golangci/golangci-lint) installed you can run:
golangci-lint run
```

## Application directory structure

```
.
├── .github                        Config to run Continuous Integration tests and checks
│  └── workflows
│     └── lint-and-test.yaml
├── cmd                            All entrypoints to the service
│  └── api                         The HTTP server entrypoint to the service
│     └── main.go
├── fizzbuzz                       Package implementing fizz-buzz
│  ├── fizzbuzz.go
│  └── fizzbuzz_test.go
├── go.mod
├── go.sum
├── http                           All code related to http is under this directory
│  ├── handlers.go
│  ├── handlers_test.go
│  ├── middleware.go
│  ├── respond.go
│  └── server.go
├── metrics                        Package responisble for collecting metrics on the service
│  ├── metrics.go
│  └── metrics_test.go
└── README.md
```

## Things that could be done differently:

## Application structure

When the service is more complex, I will usually have a `domain` and a `usercase` directory.
The domain directory will hold all domain data structures and methods that don't need to
contact any outside services.

The `domain` package will also usually contain the definition of interfaces needed to contact
outside services. Follwing this structure, I would have typically put `metrics.Metrics` in the
domain package and it's implementation in the root directory.

The `usecase` package is the glue between `domain` and the application entrypoints. Typically
the HTTP package would only ever call the `usecase` package. Using dependency injection the http
package will provide the implementations of interfaces needed by any `usecase`.

For example, I would have code similar the following code in the package `usecase`:

```go
package usecase

import (
	"context"
)

func FizzBuzz(
	fb domain.FizzBuzzer,
	metrics domain.Metrics,
) func(context.Context, domain.Request) []string {
	return func(ctx context.Context, req domain.Request) []string {
		metrics.Record(req)
		return fb.FizzBuzz(req)
	}
}

// TopHitFizzBuzz returns the most used FizzBuzzRequest
func TopHitFizzBuzz(
	metrics domain.Metrics,
) func(context.Context) (domain.MetricsResult, error) {
	return func(ctx context.Context) (domain.MetricsResult, error) {
		hits := metrics.Get()
		....
	}
}
```

### Metrics

I created an interface for metrics with an in memory implementation,
it might be an overkill for this particular example but it makes it easy to swap implementations without changing the rest of the code.

In a production application I would have certainly used [opencensus](https://opencensus.io/quickstart/go/metrics/)
for metrics instead of this implemention.

### Logging

I don't use the logger in the standard library quite often, what I usually do is use a logger that can
format logs in json and has different log levels like [uber/zap](https://github.com/uber-go/zap).
I always log to stdout and let the infrastructure deliver logs to any log management service.
