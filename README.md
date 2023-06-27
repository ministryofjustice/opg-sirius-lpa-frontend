another commit

# opg-sirius-lpa-frontend

[![codecov](https://codecov.io/gh/ministryofjustice/opg-sirius-lpa-frontend/branch/main/graph/badge.svg?token=BFGR5FBQ0T)](https://codecov.io/gh/ministryofjustice/opg-sirius-lpa-frontend)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ministryofjustice/opg-sirius-lpa-frontend)](https://pkg.go.dev/github.com/ministryofjustice/opg-sirius-lpa-frontend)

Frontend forms for Sirius.

## Quick start

### Major dependencies

- [Go](https://golang.org/) (>= 1.17)
- [Pact](https://github.com/pact-foundation/pact-ruby-standalone) (>= 1.88.82)
- [docker-compose](https://docs.docker.com/compose/install/) (>= 1.29.2)
- [Node](https://nodejs.org/en/) (>= 14.15.1)

### Running the application

```
docker-compose -f docker/docker-compose.yml up -d --build
```

This will run the application at http://localhost:8888/, and assumes that Sirius
is running at http://localhost:8080/.

Alternatively the application can be run without the use of Docker

```
yarn && yarn build
SIRIUS_PUBLIC_URL=http://localhost:8080 SIRIUS_URL=http://localhost:8080 PORT=8888 go run main.go
```

### Testing

#### Unit tests

Make sure that `pact` is available before running the tests, on a Mac with Homebrew you can do:

```
brew tap pact-foundation/pact-ruby-standalone
brew install pact-ruby-standalone
```

Then to run the tests:

```
make go-test
```

#### E2E tests (Cypress)

You can run the end-to-end tests locally with Cypress. This will start a copy
of the service with a mock backend so that you don't need to start all of the
Sirius backend and can get reliable responses.

```
make build up
yarn && yarn cypress
```

## Development

For Go we lint using [golangci-lint](https://golangci-lint.run/). It may be
useful to install locally to check changes. This will include a check on
formatting so it is recommended to setup your editor to use `go fmt`.

JavaScript files are similarly linted with [Prettier](https://prettier.io/).
You can run this locally with `yarn prettier`.
