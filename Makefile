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
	oapi-codegen --config udl/codegen.yaml docs/openapi.json

mockgen:
	mockgen -source udl/udl.gen.go -imports meroxa=github.com/meroxa/conduit-connector-udl -package mock > udl/mock/mock_client.go