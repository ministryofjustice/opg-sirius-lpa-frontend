export DOCKER_BUILDKIT=1

all: lint go-test build scan pa11y lighthouse cypress down

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

go-test:
	rm -f pacts/sirius-lpa-frontend-sirius.json
	rm -f pacts/ignored-ignored.json
	go test -count 1 ./...

build:
	docker-compose -f docker/docker-compose.ci.yml build --parallel app pact-stub

build-all:
	docker-compose -f docker/docker-compose.ci.yml build --parallel app pact-stub puppeteer cypress

scan:
	trivy sirius-lpa-frontend:latest

pa11y:
	docker-compose -f docker/docker-compose.ci.yml run --entrypoint="pa11y-ci" puppeteer

lighthouse:
	docker-compose -f docker/docker-compose.ci.yml run --entrypoint="lhci autorun" puppeteer

.PHONY: cypress
cypress:
	docker-compose -f docker/docker-compose.ci.yml run cypress

down:
	docker-compose -f docker/docker-compose.ci.yml down