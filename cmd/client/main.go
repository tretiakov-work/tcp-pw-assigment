package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/tretiakov-work/tcp-pw-assigment/internal/client"
	"github.com/tretiakov-work/tcp-pw-assigment/pkg/challenge_generator"
	"github.com/tretiakov-work/tcp-pw-assigment/pkg/message_protocol"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)

	config, input := parseCliFlagsToConfigAndInput()

	messageProtocol := message_protocol.NewZeroByteHeaderProtocol()
	challengeService := challenge_generator.NewHashcashChallengeGenerator(0, 0)

	client := client.New(&config, input, os.Stdout, messageProtocol, challengeService)

	go func() {
		if err := client.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-client.CloseCh:
		fmt.Println("Client exited")
	case sig := <-signalCh:
		fmt.Printf("Received signal: %s\n", sig)
	}

	cancel()
}

func parseCliFlagsToConfigAndInput() (client.Config, io.Reader) {
	var interactiveMode bool
	var executeCmds string
	var port int
	var host string

	flag.BoolVar(&interactiveMode, "i", true, "Enable interactive mode")
	flag.StringVar(&executeCmds, "e", "", "Comma-separated list of commands to execute")
	flag.IntVar(&port, "p", 0, "Port to connect to")
	flag.StringVar(&host, "h", "localhost", "Host to connect to")

	flag.Parse()

	commandsToExecute := strings.Split(executeCmds, ",")
	config := client.Config{
		ServerAddr:  fmt.Sprintf("%s:%d", host, port),
		Interactive: interactiveMode,
	}

	var input io.Reader
	if len(commandsToExecute) == 1 && len(commandsToExecute[0]) == 0 {
		input = os.Stdin
	} else {
		commandsToExecute = append(commandsToExecute, strconv.Itoa(int(client.UserQuit)))
		input = strings.NewReader(strings.Join(commandsToExecute, "\n") + "\n")
	}

	return config, input
}
