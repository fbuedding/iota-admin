# syntax=docker/dockerfile:1

# Build the application from source
FROM --platform=linux/amd64 golang:1.22.1-bookworm AS build-stage
ENV CGO_ENABLED=1

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      build-essential \
      libsqlite3-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" CGO_ENABLED=1 GOOS=linux go build -o /cmd/iota-admin -a -ldflags '-linkmode external -extldflags "-static"' ./cmd/iota-admin
RUN mkdir /db
RUN mkdir /db/sqlite

# Run the tests in the container
#FROM build-stage AS run-test-stage
#RUN go test -v ./...

# Deploy the application binary into a lean image
FROM --platform=linux/amd64 gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /cmd/iota-admin /iota-admin
COPY --from=build-stage --chown=nonroot:nonroot /db /db
COPY --from=build-stage --chown=nonroot:nonroot /app/assets /assets

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["./iota-admin"]

