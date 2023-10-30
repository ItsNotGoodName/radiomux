-include .env

_: init run

init:
	mkdir web/dist -p && touch web/dist/index.html

run:
	go run ./cmd/radiomux

gen: gen-proto gen-openapi gen-webrpc

tooling: tooling-air tooling-goreleaser tooling-protoc-gen-go tooling-oapi-codegen tooling-webrpc tooling-java tooling-taskfile

# Docker

docker-build-demo:
	docker build . -f docker/radiomux-demo.Dockerfile -t itsnotgoodname/radiomux-demo

# Development

dev:
	air

dev-demo:
	air -build.cmd="go build -o ./tmp/main ./cmd/radiomux-demo"

dev-web:
	cd web && pnpm install && pnpm run dev

# Generation

gen-proto:
	protoc --java_out=lite:android/app/src/main/java shared/message.proto
	protoc --go_out=. shared/message.proto

gen-openapi:
	oapi-codegen -config shared/server.cfg.yaml shared/openapi.yaml
	cd web && pnpm run generate-openapi

gen-webrpc:
	webrpc-gen -schema=./shared/api.ridl -target=golang -pkg=webrpc -server -out=./internal/webrpc/webrpc.gen.go
	webrpc-gen -schema=./shared/api.ridl -target=typescript -client -out=./web/src/api/client.gen.ts
	# Convert interface to type unless it is a Service
	awk '/interface/ && !/Service / { gsub(/interface/, "type"); gsub(/{/, "= {");} 1' ./web/src/api/client.gen.ts > tmp-3842984 && mv tmp-3842984 ./web/src/api/client.gen.ts

# Tooling

tooling-air:
	go install github.com/cosmtrek/air@latest

tooling-protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

tooling-oapi-codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

tooling-webrpc:
	go install -ldflags="-s -w -X github.com/webrpc/webrpc.VERSION=v0.13.1" github.com/webrpc/webrpc/cmd/webrpc-gen@v0.13.1

tooling-java:
	$(info Please install Java 17.)

tooling-taskfile:
	go install github.com/go-task/task/v3/cmd/task@latest
