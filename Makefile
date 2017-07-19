prom:
	prometheus -config.file=prometheus.yml

run:
	go run main.go handlers.go

stress: 
	wrk -c 100 -t 20 http://localhost:8080/city/london/temp