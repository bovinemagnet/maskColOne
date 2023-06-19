# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOBIN=./bin
GORUN=$(GOCMD) run

# Name of your binary executable
BINARY_NAME=maskColOne

all: test build

build:
	$(GOBUILD) -o $(GOBIN)/$(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	rm -rf $(GOBIN)

# Run the code to build to test doc
test_enc:
	$(GORUN) main.go --mode=e --in=test/input.tsv --out=test/output.tsv

.PHONY: all build test clean
