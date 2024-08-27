package server

import "time"

type Config struct {
	ServerAddr           string        `env:"SERVER_ADDR, default=localhost:8080"`
	ChallengeTTL         time.Duration `env:"CHALLENGE_TTL, default=2m"`
	ChallengeVersion     int           `env:"CHALLENGE_VERSION, default=1"`
	ChallengeDifficulity int           `env:"CHALLENGE_DIFFICULITY, default=10"`
}
