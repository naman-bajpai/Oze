.PHONY: build test clean install lint

BINARY := oze
MODULE := github.com/yourusername/oze

build:
	go build -o $(BINARY) .

install:
	go install .

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -f $(BINARY)

# Quick smoke tests (requires the binary to be built first)
smoke: build
	./$(BINARY) --help
	./$(BINARY) --dry-run "test feature"

