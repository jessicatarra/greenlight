include .envrc

.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN} -cors-trusted-origins=${CORS_TRUSTED_ORIGINS}