FROM golang:1.24-bookworm

RUN go install gotest.tools/gotestsum@latest

ENV PATH="$PATH:/opt/pact/bin"

RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
RUN pact-go -l DEBUG install