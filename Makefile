# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOCLEAN = $(GOCMD) clean

build: clean
	$(GOBUILD) -o $(target) ./cmd/$(target)

run: build
	./$(target)

# Test target
test:
	$(GOTEST) -v ./...

# Clean target
clean:
	$(GOCLEAN)
	rm -f $(target)