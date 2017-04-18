package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	s := &http.Server{
		Addr:         ":8080",
		Handler:      router(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
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
