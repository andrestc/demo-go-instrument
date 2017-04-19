package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/redis.v3"
)

var (
	keys chan string
)

func init() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	status := client.Ping()
	if status.Err() != nil {
		log.Fatalf("failed to initialize redis client: %s", status.Err())
	}
	keys = make(chan string)
	go startWorker(client)
	prometheus.MustRegister(&redisCollector{client: client})
	prometheus.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "myapp_redis_queue_current_length",
		Help: "The current number of itens on redis queue.",
	}, func() float64 {
		return float64(len(keys))
	}))
	prometheus.MustRegister(queueWaitDuration)
}

func startWorker(client *redis.Client) {
	fmt.Printf("Starting background worker...\n")
	for {
		select {
		case key := <-keys:
			fmt.Printf("incrementing count for key: %s.\n", key)
			result := client.Incr(key)
			if result.Err() != nil {
				fmt.Printf("failed to increment key %s: %s", key, result.Err())
			}
		case <-time.After(time.Second * 10):
			fmt.Println("pinging redis...")
			client.Ping()
		}
	}
}

var (
	queueWaitDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "myapp_redis_queue_wait_duration_seconds",
		Help: "The wait duration when trying to write to the redis queue",
	})
)

func Increment(key string) {
	select {
	case keys <- key:
	default:
		now := time.Now()
		keys <- key
		queueWaitDuration.Observe(time.Since(now).Seconds())
	}
}

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
