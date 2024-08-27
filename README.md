# Tcp server "Random word of wisdom" with Proof of Work protection

A simple tcp server with custom protocol protected by Proof of work algorithm

## Proof of work DDOS protection
In order to connect to the server, the client must solve a proof of work challenge. The client should requerst a challenge from the server, example response:
```
{"id":"dfd2047f-044e-4bd4-a5ab-b92cfb8bed02","prefix":10,"random":"OTVkZTQ1N2ItZTJhYi00ODdkLTg5ZWUtYjNmZDE0YTZkZjAy","version":1,"time":1724736371,"proof":0}
```
And then need to solve it, example solution:
```
{"id":"dfd2047f-044e-4bd4-a5ab-b92cfb8bed02","prefix":10,"random":"OTVkZTQ1N2ItZTJhYi00ODdkLTg5ZWUtYjNmZDE0YTZkZjAy","version":1,"time":1724736371,"proof":397}
```
Then the client should send the solved proof to the server and the server will return a random word of wisdom from the collection. The server uses [proof of work](https://en.wikipedia.org/wiki/Proof_of_work) challenge-response protocol. There are different algorithms to solve the proof of work:
Hashcash
- [Hashcash](https://en.wikipedia.org/wiki/Hashcash) - one of the earliest and most well-known PoW algorithms. It involves generating a hash of the challenge that meets certain criteria, such as having a specific number of leading zeroes.
- [Bitcoin Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work) - The PoW used in Bitcoin involves finding a nonce that, when hashed with the block data, produces a hash below a certain target. This process is computationally expensive and time-consuming.
- [ECASH](https://en.wikipedia.org/wiki/Ecash) - ECASH is a protocol that uses elliptic-curve cryptography to implement PoW. It requires clients to solve a cryptographic puzzle and prove their work to the server.
- [Scrypt](https://en.wikipedia.org/wiki/Scrypt) - Scrypt is a PoW algorithm that is designed to be memory-hard, meaning that it requires a large amount of memory to compute. This makes it more resistant to ASIC mining and other specialized hardware.
In this server it was decided to use a simple hashcash algorithm to solve the proof of work challenge, it is simple, well documented and has a lot of use cases in th ereal world.

## How to run
There are 2 applications:
- cmd/client/main.go - client cli application, can be used in interactive mode or with flags. Supported flags: -i - use interactive mode, -e - execute predefined commands, coma separated, -h host to connect to, -p port to connect to. Example of usage: `go run cmd/client/main.go -e "1,2,2,2" -h localhost -p "1080"` - execute predefined commands on the server(1 - challenge, 2 - get random quote)
- cmd/server/main.go - server application, uses env variables to run, supported env variables: SERVER_ADDR - host:port uri to listen to, CHALLENGE_TTL - time to live for the challenge in seconds, CHALLENGE_VERSION - version of the challenge, CHALLENGE_DIFFICULITY - prefix of the challenge, can be run with `make run target=server`

## How to build
To build the server and client applications you can use the makefile, there are 2 targets:
- make build target=server - build the server application, binary will be in /dist fodler
- make build target=client - build the client application, binary will be in /dist fodler

## Run with docker
docker-compose.yml file is provided to run the server and client applications in docker containers, use `docker-compose up`

## Protocol
The server uses a custom protocol using a zero byte and a new line as message start and end.