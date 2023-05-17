.PHONY = build dev

PORT ?= 3000

build:
	@go build -o ./tmp/http cmd/http/main.go

dev:
	@air -c ./.air.toml
