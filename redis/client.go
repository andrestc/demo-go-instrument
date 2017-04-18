package redis

import (
	"gopkg.in/redis.v3"
)

func NewClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	status := client.Ping()
	return client, status.Err()
}
