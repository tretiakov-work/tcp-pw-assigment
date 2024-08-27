//go:build integration
// +build integration

package integration_tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tretiakov-work/tcp-pw-assigment/internal/cache"
	"github.com/tretiakov-work/tcp-pw-assigment/internal/client"
	"github.com/tretiakov-work/tcp-pw-assigment/internal/server"
	"github.com/tretiakov-work/tcp-pw-assigment/pkg/challenge_generator"
	"github.com/tretiakov-work/tcp-pw-assigment/pkg/message_protocol"
	"github.com/tretiakov-work/tcp-pw-assigment/pkg/quote_generator"
)

const (
	hashCashPrefixLength = 2
	hashCashVersion      = 1
	serverAddr           = "localhost:8080"
)

func TestClientServer(t *testing.T) {
	Setup(t)

	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "challenge request",
			command:  fmt.Sprintf("%d\n", client.UserRequestChallenge),
			expected: "Response from server: \x7f\x03",
		},
		{
			name:     "quote request",
			command:  fmt.Sprintf("%d\n", client.UserRequestQuote),
			expected: "quote",
		},
		{
			name:     "help",
			command:  fmt.Sprintf("%d\n", client.Help),
			expected: "Usage:\n  0 - Quit\n  1 - Request challenge\n  2 - Request quote\n",
		},
		{
			name:     "invalid command",
			command:  fmt.Sprintf("%d\n", 100),
			expected: "Invalid command\n",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			ctx := context.Background()
			inputReader := bytes.NewReader([]byte(test.command))
			output := SetupClient(ctx, tt, inputReader)

			time.Sleep(100 * time.Millisecond)
			assert.Contains(tt, output.String(), test.expected, "missing expected response")
		})
	}
}

func Setup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	setupServer(ctx, t)

	t.Cleanup(func() {
		cancel()
	})
}

func SetupClient(ctx context.Context, t *testing.T, inputReader io.Reader) *bytes.Buffer {
	ctx, cancel := context.WithCancel(ctx)
	output := &bytes.Buffer{}

	config := &client.Config{ServerAddr: serverAddr}
	messageProtocol := message_protocol.NewZeroByteHeaderProtocol()
	challengeService := challenge_generator.NewHashcashChallengeGenerator(hashCashPrefixLength, hashCashVersion)

	client := client.New(config, inputReader, output, messageProtocol, challengeService)
	go func() {
		err := client.Start(ctx)
		require.NoError(t, err, "start client")
	}()

	t.Cleanup(func() {
		cancel()
	})

	return output
}

func setupServer(ctx context.Context, t *testing.T) {
	ctx, cancel := context.WithCancel(ctx)
	config := server.Config{
		ServerAddr:   serverAddr,
		ChallengeTTL: 1 * time.Minute,
	}
	messageProtocol := message_protocol.NewZeroByteHeaderProtocol()
	challengeService := challenge_generator.NewHashcashChallengeGenerator(hashCashPrefixLength, hashCashVersion)
	quoteService := quote_generator.NewStaticMemory()
	cacheService := cache.New()

	server := server.New(&config, messageProtocol, challengeService, quoteService, cacheService)

	go func() {
		err := server.Start(ctx)
		require.NoError(t, err, "start server")
	}()

	t.Cleanup(func() {
		cancel()
	})
}
