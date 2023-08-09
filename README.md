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
make up
```

This will run the application at http://localhost:8888/, running against the mock api server.

```
make dev
```

This will run the application at http://localhost:8888/, assumes sirius is at http://localhost:8080, and will hot reload the javascript into the app.

Alternatively the application can be run without the use of Docker

```
yarn && yarn build
SIRIUS_PUBLIC_URL=http://localhost:8080 SIRIUS_URL=http://localhost:8080 PORT=8888 go run main.go
```

### Testing

#### Unit tests

```
make unit-test
```

#### E2E tests (Cypress)

You can run the end-to-end tests locally with Cypress. This will start a copy
of the service with a mock backend so that you don't need to start all of the
Sirius backend and can get reliable responses.

```
make cypress
```

## Development

For CI Like linting locally you can run

```
make go-lint
make yarn-lint
# or to run both linters simply
make lint
```

To run the entire CI build locally just run

```
make
```
