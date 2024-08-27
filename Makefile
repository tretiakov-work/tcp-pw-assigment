# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOCLEAN = $(GOCMD) clean

build: clean
	$(GOBUILD) -o ./dist/$(target) ./cmd/$(target)

build_docker:
	docker build -f Dockerfile.$(target) .

run:
	godotenv -f .env go run cmd/$(target)/main.go

# Test target
test:
	$(GOTEST) -v ./...

# Clean target
clean:
	$(GOCLEAN)
	rm -f ./dist/$(target)