package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sethvargo/go-envconfig"
	"github.com/tretiakov-work/tcp-pw-assigment/internal/cache"
	"github.com/tretiakov-work/tcp-pw-assigment/internal/server"
	"github.com/tretiakov-work/tcp-pw-assigment/pkg/challenge_generator"
	"github.com/tretiakov-work/tcp-pw-assigment/pkg/message_protocol"
	"github.com/tretiakov-work/tcp-pw-assigment/pkg/quote_generator"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	var config server.Config
	if err := envconfig.Process(ctx, &config); err != nil {
		log.Fatal(err)
	}
	messageProtocol := message_protocol.NewZeroByteHeaderProtocol()
	challengeService := challenge_generator.NewHashcashChallengeGenerator(config.ChallengeDifficulity, config.ChallengeVersion)
	quoteService := quote_generator.NewStaticMemory()
	cacheService := cache.New()
	defer cacheService.Stop()

	server := server.New(&config, messageProtocol, challengeService, quoteService, cacheService)

	go func() {
		if err := server.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	sig := <-signalCh
	log.Printf("Received signal: %s\n", sig)

	cancel()
}
