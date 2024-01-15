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

#### 1. Developer mode with "latest" image

To **run the application in developer mode** locally, first start Sirius with the
usual `make dev-up` (in the opg-sirius project root directory).

Run the server from the opg-sirius-lpa-frontend project root with:

```
make dev
```

If the image does not exist, it will be built on demand. (Although the image has a name
that looks like it belongs in ECR, the existence of the `build` property on the lpa-frontend
service in the docker-compose.yml file will mean it's actually built locally. Consider this
name a historical remnant.)

The application should be running at http://localhost:8888/ and using the running Sirius
at http://localhost:8080 for authentication and API requests. It can be stopped with Ctrl-C.

This command continually watches for changes inside the directories under `web/assets`
(see below).

Any changes you make to JavaScript or SASS files will be reflected in the
running application (after a page refresh). For this reason, this mode is most useful when
working on UI elements, especially related to JavaScript and CSS.

Note that changes to Go code (including gohtml templates) will not be reflected in the
running application, unless you stop and rebuild the service with `make build`.

#### 2. Developer mode with locally-compiled binary

The application can also be **run without docker**. Start Sirius with `make dev-up` (as above).

Then do:

```
yarn && yarn build
yarn watch & SIRIUS_PUBLIC_URL=http://localhost:8080 SIRIUS_URL=http://localhost:8080 PORT=8888 go run main.go
```

Again, Ctrl-C stops the application.

Note that this runs the application using a binary compiled by your local Go installation.
Like the previous mode, any changes to JS or SASS files are reflected in the running application.

The plus side of this mode is that you can very quickly recompile and restart the server; this
makes it most useful when working on server side code. The only (very small) downside is that
you have to have a working Go environment locally.

#### 3. Test mode with mock Sirius and no hot reload

To run the frontend server with a mock Sirius API:

```
make up
```

This mode is used for Cypress tests. It is the least useful for a developer, as it is
slow and points the frontend at a mock Sirius API containing canned data.
It cannot be used to modify data through the UI, and can only show how canned responses
would look in the UI.

It can occasionally be necessary to start the application in this mode when debugging the Cypress
tests. Otherwise, use one of the other modes above.

#### Note on the Sirius mock

The Sirius mock server is started for all three modes, but ignored by all except the third mode,
just in case you were wondering.

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
