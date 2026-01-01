package main

import (
	"log"

	"github.com/bibektamang7/redis-clone/internals"
)

func main() {
	s := internals.NewServer(":6379")
	err := s.ListenAndServer()
	if err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}

	select {}
}
