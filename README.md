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

In this step we add instrumentation to the `redis` package by implementing the `prometheus.Collector` interface to expose some of the metrics already available on the `redis` client.

The registration of the collector is done after the client is created, on the `init` function.

```go
prometheus.MustRegister(&redisCollector{client: client})
```

The collector is a simple struct that implements two methods:

```go
var (
	requestsDesc  = prometheus.NewDesc("myapp_redis_connections_requests_total", "The total number of connections requests to redis pool.", []string{}, nil)
	hitsDesc      = prometheus.NewDesc("myapp_redis_connections_hits_total", "The total number of times a free connection was found in redis pool.", []string{}, nil)
	waitsDesc     = prometheus.NewDesc("myapp_redis_connections_waits_total", "The total number of times the redis pool had to wait for a connection.", []string{}, nil)
	timeoutsDesc  = prometheus.NewDesc("myapp_redis_connections_timeouts_total", "The total number of wait timeouts in redis pool.", []string{}, nil)
	connsDesc     = prometheus.NewDesc("myapp_redis_connections_current", "The current number of connections in redis pool.", []string{}, nil)
	freeConnsDesc = prometheus.NewDesc("myapp_redis_connections_free_current", "The current number of free connections in redis pool.", []string{}, nil)
)

type redisCollector struct {
	client *redis.Client
}

func (c *redisCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- requestsDesc
	ch <- hitsDesc
	ch <- waitsDesc
	ch <- timeoutsDesc
	ch <- connsDesc
	ch <- freeConnsDesc
}

func (c *redisCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.client.PoolStats()
	ch <- prometheus.MustNewConstMetric(requestsDesc, prometheus.CounterValue, float64(stats.Requests))
	ch <- prometheus.MustNewConstMetric(hitsDesc, prometheus.CounterValue, float64(stats.Hits))
	ch <- prometheus.MustNewConstMetric(waitsDesc, prometheus.CounterValue, float64(stats.Waits))
	ch <- prometheus.MustNewConstMetric(timeoutsDesc, prometheus.CounterValue, float64(stats.Timeouts))
	ch <- prometheus.MustNewConstMetric(connsDesc, prometheus.GaugeValue, float64(stats.TotalConns))
	ch <- prometheus.MustNewConstMetric(freeConnsDesc, prometheus.GaugeValue, float64(stats.FreeConns))
}
```

Every time Prometheus fetches the `/metrics` endpoint on our API the `Collect` function is called and sends the metrics to the provided channel.

* Branch 06
* Branch 07
* Branch 08