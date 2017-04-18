package main

import (
	"log"
	"net/http"
	"time"

	"encoding/json"

	"github.com/andrestc/go-prom-talk/weather"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var keys = make(chan string)

func main() {
	s := &http.Server{
		Addr:         ":8080",
		Handler:      router(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	go Process(keys)
	log.Println("HTTP server listening at :8080...")
	log.Fatal(s.ListenAndServe())
}

func router() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/city/{name}/temp", cityTemp)
	n := negroni.Classic()
	n.UseHandler(r)
	return n
}

func cityTemp(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		http.Error(w, "Must provide a city name.", http.StatusBadRequest)
		return
	}
	keys <- name
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
