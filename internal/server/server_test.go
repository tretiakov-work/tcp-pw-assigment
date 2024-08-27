package server

import (
	"bufio"
	"context"
	"net"
	"testing"
	"time"

	"github.com/avast/retry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tretiakov-work/tcp-pw-assigment/internal/cache"
)

const testListen = ":8000"

// Mock implementations
type MockMessageProtocol struct{}
type MockChallengeService struct{}
type MockQuoteService struct{}

func (m *MockMessageProtocol) Parse(data []byte) (int, []byte, error) {
	return int(data[0]), data[1 : len(data)-1], nil
}

func (m *MockMessageProtocol) Encode(header int, message []byte) []byte {
	message = append(message, '\n')
	return append([]byte{byte(header)}, message...)
}

func (m *MockChallengeService) GenerateChallenge(_ string) ([]byte, error) {
	return []byte("dummy-challenge"), nil
}

func (m *MockChallengeService) ValidateChallengeResponse(challenge, proof []byte) (bool, error) {
	return string(proof) == "dummy-challenge", nil
}

func (m *MockChallengeService) DeserializeChallengeID(_ []byte) (string, error) {
	return "dummy-id", nil
}

func (m *MockQuoteService) GetQuote() ([]byte, error) {
	return []byte("dummy-quote"), nil
}

func TestHandleConnection(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	mockCache := cache.New()
	server := New(&Config{ServerAddr: testListen}, &MockMessageProtocol{}, &MockChallengeService{}, &MockQuoteService{}, mockCache)
	go server.Start(ctx)

	err := retry.Do(
		func() error {
			conn, err := net.Dial("tcp", testListen)
			if err != nil {
				return err
			}
			defer conn.Close()

			return nil
		},
		retry.Attempts(10),
		retry.Delay(50*time.Millisecond),
	)
	require.NoError(t, err, "server startup")

	t.Cleanup(func() {
		cancel()
	})
	t.Parallel()

	tests := []struct {
		name        string
		message     []byte
		expectedRes string
		setup       func()
	}{
		{
			name:        "challenge request",
			message:     []byte{byte(RequestChallenge), '\n'},
			expectedRes: "dummy-challenge",
		},
		{
			name:        "quote request",
			message:     append([]byte{byte(RequestQuote)}, []byte("dummy-challenge\n")...),
			expectedRes: "dummy-quote",
			setup: func() {
				mockCache.Set("dummy-id", []byte("dummy-challenge"), 1*time.Minute)
			},
		},
		{
			name:        "invalid quote request",
			message:     append([]byte{byte(RequestQuote)}, []byte("invalid-challenge\n")...),
			expectedRes: "Challenge not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			if test.setup != nil {
				test.setup()
			}
			conn, err := net.Dial("tcp", testListen)
			assert.NoError(tt, err, "connect to server")
			defer conn.Close()

			_, err = conn.Write(test.message)
			assert.NoError(tt, err, "write message")

			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			response, err := bufio.NewReader(conn).ReadString('\n')
			assert.NoError(tt, err, "read response")

			assert.Contains(tt, response, test.expectedRes, "response")
		})
	}
}
