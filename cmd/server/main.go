package main

import (
	"context"
	"log"

	"github.com/tretiakov-work/tcp-pw-assigment/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := &server.Config{
		ServerAddr: "localhost:8080",
	}
	client := server.New(config)

	err := client.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
