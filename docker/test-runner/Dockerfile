FROM golang:1.22-bookworm

RUN go install gotest.tools/gotestsum@latest

WORKDIR /opt
RUN curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh > ./install.sh && bash ./install.sh tag=v1.92.0
ENV PATH="$PATH:/opt/pact/bin"
