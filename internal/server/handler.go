package server

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
)

type MessageHeader int

const (
	Quit MessageHeader = iota
	RequestChallenge
	RequestQuote
	ResponseChallenge
	ResponseQuote
	ResponseError
)

var (
	ErrChallengeResponse = fmt.Errorf("invalid challenge response")
	ErrQoute             = fmt.Errorf("error generating quote")
)

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		req, err := reader.ReadBytes('\n')
		if err != nil {
			if err.Error() != "EOF" {
				log.Println("Error reading request:", err.Error())
			} else {
				log.Println("client closed connection")
			}
			return
		}

		log.Println("Received request:", string(req))
		messageType, message, err := s.messageProtocol.Parse(req)
		if err != nil {
			log.Println("Error parsing request:", err.Error())
			continue
		}

		switch MessageHeader(messageType) {
		case Quit:
			log.Println("Received Quit message")
		case RequestChallenge:
			s.handleRequestChallenge(conn)
		case RequestQuote:
			s.handleRequestQuote(conn, message)
		}
	}
}

func (s *Server) handleRequestChallenge(conn net.Conn) {
	log.Println("Received Request challenge message")
	id := s.generateUniqChallengeID()
	challenge, err := s.challengeService.GenerateChallenge(id)
	if err != nil {
		log.Println("Error generating challenge:", err.Error())
	}
	s.cacheService.Set(id, challenge, s.config.ChallengeTTL)
	s.writeResponse(conn, ResponseChallenge, challenge)
}

func (s *Server) handleRequestQuote(conn net.Conn, message []byte) {
	log.Println("Received Request quote message")
	id, err := s.challengeService.DeserializeChallengeID(message)
	if err != nil {
		log.Println("Error desirializing challenge id:", err.Error())
		s.writeResponse(conn, ResponseError, []byte(err.Error()))
		return
	}
	challenge, ok := s.cacheService.Get(id)
	if !ok {
		log.Println("Challenge not found in cache")
		s.writeResponse(conn, ResponseError, []byte("Challenge not found"))
		return
	}
	success, err := s.challengeService.ValidateChallengeResponse(challenge, message)
	if err != nil || !success {
		if err != nil {
			log.Println("Error validating challenge:", err.Error())
		}
		s.writeResponse(conn, ResponseError, []byte(ErrChallengeResponse.Error()))
		return
	}
	s.cacheService.Delete(id)

	quote, err := s.quoteService.GetQuote()
	if err != nil {
		log.Println("Error generating quote:", err.Error())
		s.writeResponse(conn, ResponseError, []byte(ErrQoute.Error()))
		return
	}
	s.writeResponse(conn, ResponseQuote, quote)
}

func (s *Server) writeResponse(conn net.Conn, header MessageHeader, mess []byte) {
	response := s.messageProtocol.Encode(int(header), mess)
	_, err := conn.Write(response)
	if err != nil {
		log.Println("Error sending response:", err.Error())
	}
}

func (s *Server) generateUniqChallengeID() string {
	id := uuid.New().String()
	for _, ok := s.cacheService.Get(id); ok; {
		id = uuid.New().String()
	}
	return id
}
