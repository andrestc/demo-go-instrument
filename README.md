# Go code instrumentation using Prometheus

This repo contains a demo of code instrumentation using prometheus official
Go client.

# Introduction

The application is a single endpoint API that fetches the weather for a given city using the [openweathermap API](http://api.openweathermap.org). The code responsible for fetching the weather information is on the `weather` package.

The application also uses a redis server to store the number of requests done for each city, this logic is on the `redis` package and uses a channel to serialize the calls to the redis backend. This is done just to give different examples on how to instrument those.

# Timeline

Use the different branches to navigate between the code as we go from a completely non instrumented code to a code with some instrumentation coverage.

* Branch 01

This is the base code that implements the logic described in the introduction section.

* Branch 02

In this step we introduced the prometheus client to our vendored dependencies in included it's HTTP handler on our server:

```go
r.Handle("/metrics", promhttp.Handler())
```

With this single line of code we add several metrics about the Go runtime to our `/metrics` endpoint.

* Branch 03

* Branch 04

In this step we added instrumentation to our `weather` package. The instrumentation on this package follows the RED pattern:

* **R**equest Rate
* **E**rror rate
* **D**uration

This is done by adding the following metrics and initializing them:

```go
var (
	requestsDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "myapp_weather_request_duration_seconds",
		Help: "The duration of the requests to the weather service.",
	})

	requestsCurrent = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_weather_requests_current",
		Help: "The current number of requests to the weather service.",
	})

	requestsStatus = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "myapp_weather_requests_total",
		Help: "The total number of requests to the weather service by status.",
	}, []string{"status"})

	clientErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "myapp_weather_errors",
		Help: "The total number of weather client errors",
	})
)

func init() {
	prometheus.MustRegister(requestsDuration)
	prometheus.MustRegister(requestsCurrent)
	prometheus.MustRegister(requestsStatus)
	prometheus.MustRegister(clientErrors)
}
```

And updating the `(w *weatherClient) do` func to properly update those metrics, by adding the following snippet to the begin of the func:

```go
now := time.Now()
requestsCurrent.Inc()
defer func() {
    requestsDuration.Observe(time.Since(now).Seconds())
    requestsCurrent.Dec()
    if resp != nil {
        requestsStatus.WithLabelValues(strconv.Itoa(resp.StatusCode)).Inc()
    }
    if err != nil {
        clientErrors.Inc()
    }
}()
```



* Branch 05
* Branch 06
* Branch 07
* Branch 08