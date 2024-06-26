DOCKER_COMPOSE_RUN ?= docker-compose -f docker-compose.yml

run-server: ## Run server
	${DOCKER_COMPOSE_RUN} run --rm --service-ports server /bin/sh -c "go run cmd/tcp-work/main.go server"

test: ## Run client PoW test
	${DOCKER_COMPOSE_RUN} run --rm --service-ports client /bin/sh -c "go run cmd/tcp-work/main.go client"
	${DOCKER_COMPOSE_RUN} down

down: ## Down infra
	${DOCKER_COMPOSE_RUN} down

.PHONY: lint
lint: ## Run linter
	${DOCKER_COMPOSE_RUN} run --rm linter /bin/sh -c "golangci-lint run ./... -c .golangci.yml -v"
