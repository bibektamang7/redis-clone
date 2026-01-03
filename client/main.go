package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
	})
	ctx := context.Background()

	err := client.Set(ctx, "what", "world", 0).Err()
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
}
