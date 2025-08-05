DOCKER_COMPOSE_CMD ?= docker compose
DOCKER_COMPOSE_ENV=--env-file .env

base:
	$(DOCKER_COMPOSE_CMD) $(DOCKER_COMPOSE_ENV) -f example/base/docker-compose.yaml up -d