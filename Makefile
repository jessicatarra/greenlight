include .envrc

.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN} -cors-trusted-origins=${CORS_TRUSTED_ORIGINS}

.PHONY: run/api/help
run/api/help:
	go run ./cmd/api/ -help

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}