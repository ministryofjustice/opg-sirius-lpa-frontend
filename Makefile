export DOCKER_BUILDKIT=1

all: lint go-test build scan pa11y lighthouse cypress down

lint:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

go-test:
	rm -f pacts/sirius-lpa-frontend-sirius.json
	rm -f pacts/ignored-ignored.json
	go test -count 1 ./...

build:
	docker-compose -f docker/docker-compose.ci.yml build --parallel app

build-all:
	docker-compose -f docker/docker-compose.ci.yml build --parallel app puppeteer cypress

up:
	docker-compose -f docker/docker-compose.ci.yml up -d app

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

run-structurizr:
	docker pull structurizr/lite
	docker run -it --rm -p 8020:8080 -v $(PWD)/docs/architecture/dsl/local:/usr/local/structurizr structurizr/lite
