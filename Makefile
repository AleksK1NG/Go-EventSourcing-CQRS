.PHONY:

run_es:
	go run cmd/main.go -config=./config/config.yaml


# ==============================================================================
# Docker

dev:
	@echo Starting dev docker compose
	docker-compose -f docker-compose.yaml up -d --build

local:
	@echo Starting local docker compose
	docker-compose -f docker-compose.local.yaml up -d --build


# ==============================================================================
# Docker support

FILES := $(shell docker ps -aq)

down-local:
	docker stop $(FILES)
	docker rm $(FILES)

clean:
	docker system prune -f

logs-local:
	docker logs -f $(FILES)


# ==============================================================================
# Modules support

tidy:
	go mod tidy

deps-reset:
	git checkout -- go.mod
	go mod tidy

deps-upgrade:
	go get -u -t -d -v ./...
	go mod tidy

deps-cleancache:
	go clean -modcache


# ==============================================================================
# Linters https://golangci-lint.run/usage/install/

run-linter:
	@echo Starting linters
	golangci-lint run ./...

# ==============================================================================
# PPROF

pprof_heap:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/heap?seconds=10

pprof_cpu:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/profile?seconds=10

pprof_allocs:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/allocs?seconds=10


# ==============================================================================
# Usage:
# install local https://github.com/protocolbuffers/protobuf
# go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
# go get -u google.golang.org/grpc
# PROTO use your_file_name.proto

proto_order:
	@echo Generating es microservice order gRPC proto
	cd proto/order && protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. order.proto

# ==============================================================================
# Swagger

swagger:
	@echo Starting swagger generating
	swag init --parseDependency --parseInternal -g **/**/*.go


#docker run --name esdb-node -it -p 2113:2113 -p 1113:1113 \
#    ghcr.io/eventstore/eventstore:20.6.1-alpha.0.69-arm64v8 --insecure --run-projections=All \
#    --enable-external-tcp --enable-atom-pub-over-http