prom:
	prometheus -config.file=prometheus.yml

run:
	go run main.go handlers.go