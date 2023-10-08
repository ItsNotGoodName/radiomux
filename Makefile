# export DB_DIR=./smtpbridge_data
# export DB_FILE=smtpbridge.db
# export DB_PATH="$(DB_DIR)/$(DB_FILE)"

-include .env

# snapshot:
# 	goreleaser release --snapshot --clean

_: init run

init:
	mkdir web/dist -p && touch web/dist/index.html

run:
	go run ./cmd/radiomux

build: build-web
	CGO_ENABLED=0 go build ./cmd/radiomux

build-web:
	cd web && pnpm run build

preview: build-web run

# clean:
# 	rm -rf "$(DB_DIR)" && mkdir "$(DB_DIR)"

gen: gen-proto gen-openapi gen-webrpc # db-migrate gen-jet gen-templ

tooling: tooling-air tooling-goreleaser tooling-protoc-gen-go tooling-oapi-codegen tooling-webrpc # tooling-jet tooling-goose tooling-atlas

# Development

dev:
	air

dev-web:
	cd web && pnpm install && pnpm run dev

# Database

# db-inspect:
# 	atlas schema inspect --env local

# db-migration:
# 	atlas migrate diff $(name) --env local

# db-migrate:
# 	goose -dir migrations/sql sqlite3 "$(DB_PATH)" up

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

# gen-jet:
# 	jet -source=sqlite -dsn="$(DB_PATH)" -path=./internal/jet -ignore-tables goose_db_version,_dummy
# 	rm -rf ./internal/jet/model

# Tooling

tooling-air:
	go install github.com/cosmtrek/air@latest

tooling-goreleaser:
	go install github.com/goreleaser/goreleaser@latest

tooling-protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

tooling-oapi-codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

tooling-webrpc:
	go install -ldflags="-s -w -X github.com/webrpc/webrpc.VERSION=v0.13.1" github.com/webrpc/webrpc/cmd/webrpc-gen@v0.13.1

# tooling-jet:
# 	go install github.com/go-jet/jet/v2/cmd/jet@latest
 
# tooling-goose:
# 	go install github.com/pressly/goose/v3/cmd/goose@latest
 
# tooling-atlas:
# 	go install ariga.io/atlas/cmd/atlas@latest
 
# tooling-templ:
# 	go install github.com/a-h/templ/cmd/templ@latest
