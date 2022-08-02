.PHONY: build test integration lint

build:
	go build -o conduit-connector-udl cmd/udl/main.go

test:
	go test $(GOTEST_FLAGS) -race ./...

integration:
	INTEGRATION_TEST=true go test $(GOTEST_FLAGS) -race ./...

generate:
	go generate ./...

lint:
	golangci-lint run
