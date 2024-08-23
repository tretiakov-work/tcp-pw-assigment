package main

import (
	"context"
	"log"
	"os"

	"github.com/tretiakov-work/tcp-pw-assigment/internal/client"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := &client.Config{
		ServerAddr: "localhost:8080",
	}
	client := client.New(config, os.Stdout)

	err := client.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
