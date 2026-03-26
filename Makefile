PROJECT_NAME := embedded-lab-project
API_IMAGE := embedded-lab-api:latest
DEV_COMPOSE := docker compose -f docker-compose.dev.yml
PROD_COMPOSE := docker compose -f docker-compose.prod.yml

.PHONY: help dev-up dev-down dev-logs prod-up prod-down prod-logs build-api rebuild-api ps mongo-logs express-logs api-logs mongo-shell

help:
	@echo "Available targets:"
	@echo "  make dev-up       Start MongoDB and mongo-express for development"
	@echo "  make dev-down     Stop development services"
	@echo "  make dev-logs     Show development service logs"
	@echo "  make build-api    Build the API Docker image"
	@echo "  make rebuild-api  Rebuild the API Docker image without cache"
	@echo "  make prod-up      Start the production stack"
	@echo "  make prod-down    Stop the production stack"
	@echo "  make prod-logs    Show production stack logs"
	@echo "  make ps           List running containers for this project"
	@echo "  make mongo-logs   Show MongoDB logs"
	@echo "  make express-logs Show mongo-express logs"
	@echo "  make api-logs     Show API logs"
	@echo "  make mongo-shell  Open mongosh inside the MongoDB container"

dev-up:
	$(DEV_COMPOSE) up -d

dev-down:
	$(DEV_COMPOSE) down

dev-logs:
	$(DEV_COMPOSE) logs -f

build-api:
	docker build -t $(API_IMAGE) .

rebuild-api:
	docker build --no-cache -t $(API_IMAGE) .

prod-up:
	$(PROD_COMPOSE) up --build -d

prod-down:
	$(PROD_COMPOSE) down

prod-logs:
	$(PROD_COMPOSE) logs -f

ps:
	docker ps --filter "name=embedded-lab"

mongo-logs:
	$(DEV_COMPOSE) logs -f mongodb

express-logs:
	$(DEV_COMPOSE) logs -f mongo-express

api-logs:
	$(PROD_COMPOSE) logs -f api

mongo-shell:
	$(DEV_COMPOSE) exec mongodb mongosh -u admin -p admin --authenticationDatabase admin
