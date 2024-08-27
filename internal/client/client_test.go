package client

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockTcpServer struct {
}

// Mock implementation of MessageProtocol
type MockMessageProtocol struct{}

func (m *MockMessageProtocol) Parse(data []byte) (int, []byte, error) {
	return int(data[0]), data[1:], nil
}

func (m *MockMessageProtocol) Encode(messageType int, message []byte) []byte {
	message = append(message, '\n')
	return append([]byte{byte(messageType)}, message...)
}

type MockChallengeService struct{}

func (m *MockChallengeService) SolveChallenge(challengeBytes []byte) ([]byte, error) {
	return challengeBytes, nil
}

// Create a mock server to simulate responses
func startMockServer(t *testing.T) net.Addr {
	ln, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to start mock server: %v", err)
	}

	go func() {
		defer ln.Close()
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(conn net.Conn) {
				defer conn.Close()

				reader := bufio.NewReader(conn)
				message, err := reader.ReadBytes('\n')
				if err != nil {
					return
				}

				fmt.Fprintln(conn, string(message))
			}(conn)
		}
	}()

	return ln.Addr()
}

func TestClient(t *testing.T) {
	serverAddr := startMockServer(t)

	// Test commands
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{
			name:     "challenge request",
			command:  fmt.Sprintf("%d\n", UserRequestChallenge),
			expected: "Response from server: \x01\n",
		},
		{
			name:     "help",
			command:  fmt.Sprintf("%d\n", Help),
			expected: "Usage:\n  0 - Quit\n  1 - Request challenge\n  2 - Request quote\n",
		},
		{
			name:     "invalid command",
			command:  fmt.Sprintf("%d\n", 100),
			expected: "Invalid command\n",
		},
	}

	t.Parallel()
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			input := &bytes.Buffer{}

			output, client := SetupClient(tt, serverAddr.String(), input)
			defer client.Stop()

			_, err := input.Write([]byte(test.command))
			assert.NoError(tt, err, "write command to input file")

			time.Sleep(100 * time.Millisecond)

			assert.Contains(tt, output.String(), test.expected)
		})
	}
}

func SetupClient(t *testing.T, serverAddr string, inputReader io.Reader) (*bytes.Buffer, *Client) {
	output := &bytes.Buffer{}

	protocol := &MockMessageProtocol{}
	challengeService := &MockChallengeService{}

	client := New(&Config{ServerAddr: serverAddr}, inputReader, output, protocol, challengeService)
	go func() {
		err := client.Start(context.Background())
		require.NoError(t, err, "start client")
	}()

	return output, client
}
