package weather

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	apiURL = "http://api.openweathermap.org"
)

var (
	client *weatherClient
)

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

type weatherClient struct {
	httpClient *http.Client
	apiKey     string
}

func (w *weatherClient) do(method string, path string, params map[string]string) (resp *http.Response, err error) {
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
	url := fmt.Sprintf("%s/%s?appid=%s", apiURL, path, w.apiKey)
	for k, v := range params {
		url += fmt.Sprintf("&%s=%s", k, v)
	}
	fmt.Printf("DEBUG: %s\n", url)
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	return w.httpClient.Do(request)
}

func getClient() (*weatherClient, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, errors.New("must set WEATHER_API_KEY env")
	}
	if client == nil {
		client = &weatherClient{
			httpClient: &http.Client{Timeout: time.Second * 15},
			apiKey:     apiKey,
		}
	}
	return client, nil
}
