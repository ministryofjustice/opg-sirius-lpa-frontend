export DOCKER_BUILDKIT=1

help:
	@grep --no-filename -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: lint gosec unit-test build-all scan pa11y lighthouse cypress down

lint: ## Lint source code
lint: go-lint yarn-lint

go-lint:
	docker compose run --rm go-lint

yarn-lint:
	docker compose run --rm yarn
	docker compose run --rm yarn lint

gosec: ## Scan Go code for security flaws
	docker compose run --rm gosec

test-results:
	mkdir -p -m 0777 test-results .gocache pacts logs cypress/screenshots .trivy-cache

setup-directories: test-results

unit-test: ## Run Go unit tests
unit-test: setup-directories
	docker compose run --rm test-runner

build:
	docker compose build lpa-frontend

build-all: ## Build containers
	docker compose build --parallel lpa-frontend puppeteer cypress test-runner

dev: ## Start application and watch JS and SASS files for changes
	docker compose run --rm yarn
	docker compose run --rm yarn build
	docker compose -f docker-compose.yml -f docker/docker-compose.dev.yml up -d lpa-frontend
	docker compose run --rm yarn watch

up: ## Start application with mock Sirius API; mostly for use with Cypress tests
	docker compose up -d lpa-frontend

scan: setup-directories
	docker compose run --rm trivy image --format table --exit-code 0 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-lpa-frontend:latest
	docker compose run --rm trivy image --format sarif --output /test-results/trivy.sarif --exit-code 1 311462405659.dkr.ecr.eu-west-1.amazonaws.com/sirius/sirius-lpa-frontend:latest

pa11y: setup-directories
	docker compose run --entrypoint="pa11y-ci" puppeteer

lighthouse: setup-directories
	docker compose run --entrypoint="lhci autorun" puppeteer

cypress: setup-directories
	docker compose run --rm cypress

down: ## Stop everything
	docker compose down

run-structurizr:
	docker pull structurizr/lite
	docker run -it --rm -p 8020:8080 -v $(PWD)/docs/architecture/dsl/local:/usr/local/structurizr structurizr/lite
