FROM node:21.1.0-alpine3.17 as asset-env

WORKDIR /app

RUN apk add --no-cache git

RUN mkdir -p web/static

COPY package.json .
COPY yarn.lock .
RUN yarn

COPY web/assets web/assets
RUN yarn build

FROM golang:1.22.2 as build-env

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go main.go
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/opg-sirius-lpa-frontend

FROM alpine:3

WORKDIR /go/bin

RUN apk --update --no-cache add \
    ca-certificates \
    tzdata

# Patch vulnerabilities
RUN apk update && apk upgrade --no-cache libcrypto3 libssl3

COPY --from=build-env /go/bin/opg-sirius-lpa-frontend opg-sirius-lpa-frontend
COPY --from=asset-env /app/web/static web/static
COPY web/template web/template

RUN addgroup -S app && \
    adduser -S -g app app && \
    chown -R app:app opg-sirius-lpa-frontend web/template web/static
USER app
ENTRYPOINT ["./opg-sirius-lpa-frontend"]
