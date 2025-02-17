FROM golang:1.24-bookworm

RUN go install gotest.tools/gotestsum@latest

WORKDIR /opt
ENV PACT_CLI_VERSION=v2.4.2
ADD https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh ./install.sh
RUN bash ./install.sh
ENV PATH="$PATH:/opt/pact/bin"
