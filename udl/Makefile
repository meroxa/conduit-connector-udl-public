.PHONY: build test integration lint generate_client

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

generate_client:
	oapi-codegen --config codegen.yaml docs/openapi.json

mockgen:
	mockgen -source udl.gen.go -imports meroxa=github.com/meroxa/udl -package mock > mock/mock_client.go