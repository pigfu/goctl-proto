.PHONY: all

all: build run goctl rpc

build:
	go vet ./...
	go build -ldflags "-X main.buildTime=`date +%Y-%m-%dT%H:%M:%S`" -o ./releases/goctl-proto ./cmd/proto

run:
	@if [ -x "./releases/goctl-proto" ]; \
    then \
        ./releases/goctl-proto proto --input ./example/api/service.api --output ./example/proto; \
    else \
        go run ./cmd/proto/*.go proto --input ./example/api/service.api --output ./example/proto; \
    fi

goctl:
	@if ! [ -x "./releases/goctl-proto" ]; \
    then \
        make build; \
    fi
	goctl api plugin -plugin ./releases/goctl-proto="proto" -api ./example/api/service.api -dir ./example/proto

rpc:
	goctl rpc protoc ./example/proto/service.proto --go_out=./example/service --go-grpc_out=./example/service --zrpc_out=./example/service --client=true

clean:
	rm -rf ./releases
	rm -rf ./example/proto/service.proto
	rm -rf ./example/service
