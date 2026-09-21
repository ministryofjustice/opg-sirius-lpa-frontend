FROM golang:1.27.0-bookworm@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

RUN go tool pact-go -l DEBUG install
