FROM golang:1.27.1-bookworm@sha256:8d48e12ec56735e9358640898b9d9b9fcca110612ed8a5567438c0a1baa24e66

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

RUN go tool pact-go -l DEBUG install
