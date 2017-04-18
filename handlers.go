package main

import (
	"encoding/json"
	"net/http"

	"github.com/andrestc/go-prom-talk/redis"
	"github.com/andrestc/go-prom-talk/weather"
	"github.com/gorilla/mux"
)

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
