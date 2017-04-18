package redis

import (
	"fmt"
	"log"
	"time"

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

func Increment(key string) {
	keys <- key
}
