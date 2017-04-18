package main

import (
	"fmt"
	"time"

	"github.com/andrestc/go-prom-talk/redis"
)

func Process(ch chan string) {
	fmt.Printf("Starting background worker...\n")
	redisClient, err := redis.NewClient()
	if err != nil {
		panic(err)
	}
	for {
		select {
		case key := <-ch:
			fmt.Printf("incrementing count for key: %s.\n", key)
			result := redisClient.Incr(key)
			if result.Err() != nil {
				fmt.Printf("failed to increment key %s: %s", key, result.Err())
			}
		case <-time.After(time.Second * 10):
			fmt.Println("pinging redis...")
			redisClient.Ping()
		}
	}
}
