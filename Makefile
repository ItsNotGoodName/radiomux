# export DB_DIR=./smtpbridge_data
# export DB_FILE=smtpbridge.db
# export DB_PATH="$(DB_DIR)/$(DB_FILE)"

-include .env

# snapshot:
# 	goreleaser release --snapshot --clean

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

gen: gen-proto gen-openapi # db-migrate gen-jet gen-templ

tooling: tooling-air tooling-goreleaser tooling-protoc-gen-go tooling-oapi-codegen # tooling-jet tooling-goose tooling-atlas

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

# tooling-jet:
# 	go install github.com/go-jet/jet/v2/cmd/jet@latest
 
# tooling-goose:
# 	go install github.com/pressly/goose/v3/cmd/goose@latest
 
# tooling-atlas:
# 	go install ariga.io/atlas/cmd/atlas@latest
 
# tooling-templ:
# 	go install github.com/a-h/templ/cmd/templ@latest
