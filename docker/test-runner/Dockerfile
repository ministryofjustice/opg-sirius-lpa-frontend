FROM golang:1.21-bookworm@sha256:640f192f24852cd875514e18a01977ae87692dca466abed7705ad3c4670b6993

RUN go install gotest.tools/gotestsum@latest

ADD https://github.com/pact-foundation/pact-ruby-standalone/releases/download/v1.88.82/pact-1.88.82-linux-x86_64.tar.gz /usr/bin
RUN cd /usr/bin; tar xzf pact-1.88.82-linux-x86_64.tar.gz
ENV PATH="$PATH:/usr/bin/pact/bin"
