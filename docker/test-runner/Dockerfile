FROM golang:1.23-bookworm

RUN go install gotest.tools/gotestsum@latest

WORKDIR /opt
ENV PACT_CLI_VERSION=v2.4.2
RUN curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh > ./install.sh && bash ./install.sh 
ENV PATH="$PATH:/opt/pact/bin"
