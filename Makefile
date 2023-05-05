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

up:
	docker-compose -f docker/docker-compose.ci.yml up -d

dev:
	docker-compose -f docker/docker-compose.yml -f docker/docker-compose.dev.yml up -d
	yarn && yarn watch

scan:
	trivy image sirius-lpa-frontend:latest

pa11y:
	docker-compose -f docker/docker-compose.ci.yml run --entrypoint="pa11y-ci" puppeteer

lighthouse:
	docker-compose -f docker/docker-compose.ci.yml run --entrypoint="lhci autorun" puppeteer

down:
	docker-compose -f docker/docker-compose.ci.yml down
