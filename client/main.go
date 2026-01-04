package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	ctx := context.Background()

	err := client.Set(ctx, "new", "world", 0).Err()
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	str, err := client.Get(ctx, "new").Result()
	if err != nil {
		log.Fatal("ERROR: ", err)
	}

	fmt.Println("value after get", str)

}
