DOCKER_COMPOSE_RUN ?= docker-compose -f docker-compose.yml

## Default goal
.DEFAULT_GOAL := default

default: create_volumes lint test

run-server: ## Run server
	${DOCKER_COMPOSE_RUN} run --rm --service-ports server /bin/sh -c "go run cmd/tcp-work/main.go server"

test: create_volumes ## Run client PoW test
	${DOCKER_COMPOSE_RUN} run --rm --service-ports client /bin/sh -c "go run cmd/tcp-work/main.go client"
	${DOCKER_COMPOSE_RUN} down

down: ## Down infra
	${DOCKER_COMPOSE_RUN} down

.PHONY: lint
lint: ## Run linter
	${DOCKER_COMPOSE_RUN} run --rm linter /bin/sh -c "golangci-lint run ./... -c .golangci.yml -v"

.PHONY: create_volumes
create_volumes: ## Create docker cache volumes
	docker volume create go-mod-cache
	docker volume create go-build-cache
	docker volume create go-lint-cache
