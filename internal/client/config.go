package client

type Config struct {
	ServerAddr  string `env:"CLIENT_SERVER_ADDR, default=localhost:8080"`
	Interactive bool   `env:"CLIENT_INTERACTIVE, default=true"`
}
