FROM golang:1.25.5-bookworm@sha256:93760d0a33d57d7a9e635803f6f5642f3da57296fc777030434cd0d3cb04a956

RUN go install gotest.tools/gotestsum@latest

RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
RUN pact-go -l DEBUG install
