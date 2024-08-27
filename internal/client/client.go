package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Client struct {
	conn                   net.Conn
	serverConnectionReader *bufio.Reader
	config                 *Config
	output                 io.Writer
	input                  io.Reader
	messageProtocol        MessageProtocol
	CloseCh                chan struct{}
	wg                     sync.WaitGroup
	challengeService       ChallengeService
}

type MessageProtocol interface {
	Parse(data []byte) (int, []byte, error)
	Encode(messageType int, message []byte) []byte
}

type ChallengeService interface {
	SolveChallenge(challengeBytes []byte) ([]byte, error)
}

type UserCommand int
type ServerCommand int

const (
	UserQuit UserCommand = iota
	UserRequestChallenge
	UserRequestQuote
	Help
)

const (
	ServerQuit ServerCommand = iota
	ServerRequestChallenge
	ServerRequestQuote
	ServerResponseChallenge
	ServerResponseQuote
	ServerResponseError
)

func New(
	config *Config,
	input io.Reader,
	output io.Writer,
	messageProtocol MessageProtocol,
	challengeService ChallengeService) *Client {
	return &Client{
		config:           config,
		input:            input,
		output:           output,
		messageProtocol:  messageProtocol,
		CloseCh:          make(chan struct{}),
		challengeService: challengeService,
	}
}

func (c *Client) Start(ctx context.Context) error {
	var err error
	c.conn, err = net.Dial("tcp", c.config.ServerAddr)

	c.serverConnectionReader = bufio.NewReader(c.conn)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		defer c.conn.Close()

		userReader := bufio.NewReader(c.input)

		for {
			select {
			case <-ctx.Done():
				return
			case <-c.CloseCh:
				return
			default:
				if c.config.Interactive {
					fmt.Fprint(c.output, "Enter command: ")
				}

				message, _ := userReader.ReadString('\n')
				if len(message) == 0 {
					fmt.Fprintln(c.output, "Invalid command")
					continue
				}

				message = strings.TrimSpace(message)
				command, err := strconv.Atoi(message)
				if err != nil {
					fmt.Fprintln(c.output, "Invalid command")
					usage(c.output)
					continue
				}

				switch UserCommand(command) {
				case UserQuit:
					c.Stop()
					return
				case UserRequestChallenge:
					c.handleRequestChallenge()
				case UserRequestQuote:
					c.handleRequestQuote()
				case Help:
					usage(c.output)
					continue
				default:
					fmt.Fprintln(c.output, "Invalid command")
					continue
				}
			}
		}
	}()

	c.wg.Wait()

	return nil
}

func (c *Client) handleRequestChallenge() {
	fmt.Fprintln(c.output, "Performing challenge request")
	outputMessage := c.messageProtocol.Encode(int(ServerRequestChallenge), nil)
	c.conn.Write(outputMessage)
	response, err := c.serverConnectionReader.ReadString('\n')
	if err != nil && err.Error() != "EOF" {
		fmt.Fprintln(c.output, "failed to read response from server:", err)
		return
	}

	fmt.Fprintln(c.output, "Response from server:", response)
}

func (c *Client) handleRequestQuote() {
	fmt.Fprintln(c.output, "Performing quote request")
	outputMessage := c.messageProtocol.Encode(int(ServerRequestChallenge), nil)
	c.conn.Write(outputMessage)
	responseRaw, err := c.serverConnectionReader.ReadBytes('\n')
	if err != nil && err.Error() != "EOF" {
		fmt.Fprintln(c.output, "failed to read response from server:", err)
		return
	}
	header, response, err := c.messageProtocol.Parse(responseRaw)
	if err != nil {
		fmt.Fprintln(c.output, "failed to parse response from server:", err)
		return
	}
	if header != int(ServerResponseChallenge) {
		fmt.Fprintln(c.output, "unexpected header from server", header)
		return
	}
	fmt.Fprintln(c.output, "Response challenge from server:", string(response))
	solvedChallenge, err := c.challengeService.SolveChallenge(response)
	if err != nil {
		fmt.Fprintln(c.output, "failed to solve challenge:", err)
		return
	}
	outputMessage = c.messageProtocol.Encode(int(ServerRequestQuote), solvedChallenge)
	c.conn.Write(outputMessage)

	response, err = c.serverConnectionReader.ReadBytes('\n')
	if err != nil && err.Error() != "EOF" {
		fmt.Fprintln(c.output, "failed to read response from server:", err)
		return
	}

	fmt.Fprintln(c.output, "Response qoute from server:", string(response))
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  0 - Quit")
	fmt.Fprintln(w, "  1 - Request challenge")
	fmt.Fprintln(w, "  2 - Request quote")
}

func (c *Client) Stop() {
	close(c.CloseCh)
}
