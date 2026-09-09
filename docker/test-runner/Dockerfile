FROM golang:1.26.6-bookworm@sha256:116d58cbd88c1297624acc6e967a060012422bacf9930927e23fb719189c6f36

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

RUN go tool pact-go -l DEBUG install
