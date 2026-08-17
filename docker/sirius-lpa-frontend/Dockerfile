FROM node:24.19.0-alpine3.23@sha256:244cc2b53f46f9e876304391d17682b0ddae9ac33491f4857e25e35a36ba7995 AS asset-env

WORKDIR /app

RUN apk add --no-cache git \
    && mkdir -p web/static

COPY package.json .
COPY package-lock.json .
RUN npm ci --ignore-scripts

COPY web/assets web/assets
RUN npm run build

FROM golang:1.26.5-alpine@sha256:0178a641fbb4858c5f1b48e34bdaabe0350a330a1b1149aabd498d0699ff5fb2 AS dev

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

ENTRYPOINT ["go", "tool", "air"]

FROM golang:1.26.5@sha256:3aff6657219a4d9c14e27fb1d8976c49c29fddb70ba835014f477e1c70636647 AS build-env

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid 65532 \
    app

ARG TARGETARCH
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY main.go main.go
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o /go/bin/opg-sirius-lpa-frontend

FROM build-env AS healthcheck-build
WORKDIR /app

COPY healthcheck healthcheck

WORKDIR /app/healthcheck
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/healthcheck

FROM scratch

WORKDIR /go/bin

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build-env /etc/passwd /etc/passwd
COPY --from=build-env /etc/group /etc/group

COPY --from=build-env /go/bin/opg-sirius-lpa-frontend opg-sirius-lpa-frontend
COPY --from=healthcheck-build /go/bin/healthcheck healthcheck
COPY --from=asset-env /app/web/static web/static
COPY web/template web/template

USER app

HEALTHCHECK --interval=5s --timeout=5s --start-period=5s --retries=3 CMD [ "/go/bin/healthcheck" ]

ENTRYPOINT ["./opg-sirius-lpa-frontend"]
