FROM golang:1.26.5-bookworm@sha256:1ecb7edf62a0408027bd5729dfd6b1b8766e578e8df93995b225dfd0944eb651

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go get -tool github.com/pact-foundation/pact-go/v2 gotest.tools/gotestsum

RUN go tool pact-go -l DEBUG install
