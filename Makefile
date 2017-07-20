prom:
	prometheus -config.file=prometheus.yml

run:
	go run main.go handlers.go

stress: 
	wrk -c 100 -t 20 http://localhost:8080/city/london/temp

grafana:
	grafana-server --config=/usr/local/etc/grafana/grafana.ini --homepath /usr/local/share/grafana cfg:default.paths.logs=/usr/local/var/log/grafana cfg:default.paths.data=/usr/local/var/lib/grafana cfg:default.paths.plugins=/usr/local/var/lib/grafana/plugins