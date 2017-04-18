package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
)

var (
	openConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "myapp_open_connections",
		Help: "The current number of open connections.",
	})
)

func init() {
	prometheus.MustRegister(openConnections)
}

func main() {
	s := &http.Server{
		Addr:         ":8080",
		Handler:      router(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		ConnState: func(c net.Conn, s http.ConnState) {
			switch s {
			case http.StateNew:
				openConnections.Inc()
			case http.StateHijacked | http.StateClosed:
				openConnections.Dec()
			}
		},
	}
	log.Println("HTTP server listening at :8080...")
	log.Fatal(s.ListenAndServe())
}

func router() http.Handler {
	r := mux.NewRouter().StrictSlash(true)
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/city/{name}/temp", instrumentHandler("/city/{name}/temp", cityTemp))
	n := negroni.Classic()
	n.UseHandler(r)
	return n
}
