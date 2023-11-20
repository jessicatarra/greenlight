include .envrc

.PHONY: run/api
run/api:
	go run ./cmd/api -db-dsn=${DATABASE_URL} -cors-trusted-origins=${CORS_TRUSTED_ORIGINS} -jwt-secret=${JWT_SECRET}

.PHONY: run/api/help
run/api/help:
	go run ./cmd/api/ -help

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...


.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...' go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

.PHONY: generate/api/docs
generate/api/docs:
	@echo 'Remove docs...'
	rm -rf docs
	@echo 'Generate updated docs folder'
	swag init -d cmd/api --parseDependency --parseInternal --parseDepth 2

# ====================================================================================
# # BUILD
# ==================================================================================== #
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-X main.version=${git_description}'

.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
