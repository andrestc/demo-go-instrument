package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/andrestc/go-prom-talk/redis"
	"github.com/andrestc/go-prom-talk/weather"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	handlerDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "myapp_handlers_duration_seconds",
		Help: "Handlers request duration in seconds",
	}, []string{"path"})
)

func init() {
	prometheus.MustRegister(handlerDuration)
}

func instrumentHandler(pattern string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next(w, r)
		handlerDuration.WithLabelValues(pattern).Observe(time.Since(now).Seconds())
	})
}

func cityTemp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		http.Error(w, "Must provide a city name.", http.StatusBadRequest)
		return
	}
	redis.Increment(name)
	temp, err := weather.GetCityTemp(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(temp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
