package client

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
)

type Client struct {
	conn   net.Conn
	config *Config
	output *os.File
}

func New(config *Config, output *os.File) *Client {
	return &Client{
		config: config,
		output: output,
	}
}

func (c *Client) Start(ctx context.Context) error {
	var err error
	c.conn, err = net.Dial("tcp", c.config.ServerAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	defer c.conn.Close()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Fprint(c.output, "Enter message: ")
		message, _ := reader.ReadString('\n')

		fmt.Fprint(c.conn, message)

		response, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read response from server: %w", err)
		}

		fmt.Fprintln(c.output, "Response from server:", response)
	}
}
