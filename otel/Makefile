DOCKER_COMPOSE_CMD ?= docker compose
DOCKER_COMPOSE_ENV=--env-file .env
SECRET_ENV=--env-file secret.env

base:
	$(DOCKER_COMPOSE_CMD) $(DOCKER_COMPOSE_ENV) $(SECRET_ENV) -f example/base/docker-compose.yaml up -d