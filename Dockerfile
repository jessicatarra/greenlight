# Build stage
ARG GO_VERSION=1.21
FROM golang:${GO_VERSION}-alpine3.17 AS builder
WORKDIR /bin
COPY . .
RUN go mod download && go mod verify && go mod vendor
ARG API_VERSION="v0.0.0+unknown"
ARG API_PORT="unknown"
ARG API_ENV="unknown"
RUN go build -ldflags "-X 'main.version=$API_VERSION' -X 'main.port=$API_PORT' -X 'main.env=$API_ENV'" -o api ./cmd/api

# Run stage
FROM golang:${GO_VERSION}-alpine3.17 AS build-release-stage
WORKDIR /bin
COPY --from=builder /bin/api .

EXPOSE 8080
ENTRYPOINT [ "/bin/api" ]