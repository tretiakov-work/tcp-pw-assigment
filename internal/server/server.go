package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	config           *Config
	messageProtocol  MessageProtocol
	challengeService ChallengeService
	quoteService     QuoteService
	cacheService     CacheService
	connections      chan net.Conn
	listener         net.Listener
	wg               sync.WaitGroup
}

type MessageProtocol interface {
	Parse(data []byte) (int, []byte, error)
	Encode(messageType int, message []byte) []byte
}

type ChallengeService interface {
	GenerateChallenge(id string) ([]byte, error)
	ValidateChallengeResponse(challenge, proof []byte) (bool, error)
	DeserializeChallengeID(challengeBytes []byte) (string, error)
}

type QuoteService interface {
	GetQuote() ([]byte, error)
}

type CacheService interface {
	Set(key string, value []byte, ttl time.Duration)
	Get(key string) ([]byte, bool)
	Delete(key string)
}

func New(
	config *Config,
	messageProtocol MessageProtocol,
	challengeService ChallengeService,
	quoteService QuoteService,
	cacheService CacheService) *Server {
	return &Server{
		config:           config,
		messageProtocol:  messageProtocol,
		challengeService: challengeService,
		quoteService:     quoteService,
		connections:      make(chan net.Conn),
		cacheService:     cacheService,
	}
}

func (s *Server) Start(ctx context.Context) error {
	var err error
	s.listener, err = net.Listen("tcp", s.config.ServerAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	log.Println("Server listening on", s.config.ServerAddr)

	s.connections = make(chan net.Conn)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Op == "accept" {
					// Handle server shutdown or network errors here
					log.Printf("accept error: %v", err)
					return
				}
				log.Printf("failed to accept connection: %v", err)
				return
			}
			select {
			case s.connections <- conn:
			case <-ctx.Done():
				conn.Close()
				return
			}
		}
	}()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-ctx.Done():
				close(s.connections)
				return
			case conn := <-s.connections:
				if conn != nil {
					go s.handleConnection(conn)
				}
			}
		}
	}()

	<-ctx.Done()
	s.listener.Close()
	s.wg.Wait()

	return nil
}
