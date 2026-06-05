FROM golang:1.26.3-bookworm@sha256:386d475a660466863d9f8c766fec64d7fdad3edac2c6a05020c09534d71edb4b

RUN go install gotest.tools/gotestsum@latest

RUN go install github.com/pact-foundation/pact-go/v2@v2.2.0
RUN pact-go -l DEBUG install
