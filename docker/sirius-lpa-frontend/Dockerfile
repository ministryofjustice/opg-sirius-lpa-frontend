FROM node:22.0.0-alpine3.19 AS asset-env

WORKDIR /app

RUN apk add --no-cache git \
    && mkdir -p web/static

COPY package.json .
COPY yarn.lock .
RUN yarn --ignore-scripts

COPY web/assets web/assets
RUN yarn build

FROM golang:1.24 AS build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go main.go
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/opg-sirius-lpa-frontend

FROM alpine:3

WORKDIR /go/bin

RUN apk update \
    && apk add --no-cache ca-certificates tzdata \
    && apk upgrade --no-cache busybox libcrypto3 libssl3

COPY --from=build-env /go/bin/opg-sirius-lpa-frontend opg-sirius-lpa-frontend
COPY --from=asset-env /app/web/static web/static
COPY web/template web/template

RUN addgroup -S app && \
    adduser -S -g app app && \
    chown -R app:app opg-sirius-lpa-frontend web/template web/static
USER app
ENTRYPOINT ["./opg-sirius-lpa-frontend"]
