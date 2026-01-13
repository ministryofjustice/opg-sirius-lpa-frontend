FROM node:22.0.0-alpine3.19@sha256:9459e243f620fff19380c51497493e91db1454f5a30847efe5bc5e50920748d5 AS asset-env

WORKDIR /app

RUN apk add --no-cache git \
    && mkdir -p web/static

COPY package.json .
COPY package-lock.json .
RUN npm ci --ignore-scripts

COPY web/assets web/assets
RUN npm run build

FROM golang:1.25.5-alpine@sha256:ac09a5f469f307e5da71e766b0bd59c9c49ea460a528cc3e6686513d64a6f1fb AS dev

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

# We need to copy the source files and config for Air to work when the container is run standalone for Sirius
COPY main.go main.go
COPY internal internal
COPY web/template web/template

COPY .air.toml .air.toml

RUN apk update \
    && apk add --no-cache ca-certificates tzdata \
    && apk upgrade --no-cache busybox libcrypto3 libssl3

COPY --from=asset-env /app/web/static web/static

RUN go install github.com/cosmtrek/air@v1.52.1 && go install github.com/go-delve/delve/cmd/dlv@latest

ENTRYPOINT ["air"]

FROM golang:1.25.5@sha256:0f406d34b7cb7255d0700af02ec28a2c88f1e00701055f4c282aa4c3ec0b3245 AS build-env

ARG TARGETARCH
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go main.go
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o /go/bin/opg-sirius-lpa-frontend

FROM alpine:3@sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62

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
