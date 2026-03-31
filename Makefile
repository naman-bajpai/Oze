.PHONY: build install clean test

BINARY := oze
INSTALL_DIR := /usr/local/bin

build:
	go build -o $(BINARY) .

install: build
	mv $(BINARY) $(INSTALL_DIR)/$(BINARY)
	@echo "Installed to $(INSTALL_DIR)/$(BINARY)"

clean:
	rm -f $(BINARY)

test:
	go test ./...

vet:
	go vet ./...

