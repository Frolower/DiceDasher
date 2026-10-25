.DEFAULT_GOAL := help

COMPOSE ?= docker compose
GO ?= go
NPM ?= npm
SERVICE ?=
GO_PACKAGES := ./pkg/... ./services/backend/... ./services/character/... ./services/resolve/...
GO_DIRS := pkg services/backend services/character services/resolve

.PHONY: help build up down logs ps restart clean rebuild nuke fresh dev \
	install frontend-dev frontend-build frontend-preview test test-go test-frontend \
	lint lint-go lint-frontend fmt fmt-go fmt-frontend fmt-check

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*## "; print "Usage: make <target> [SERVICE=backend]"; print ""} /^[a-zA-Z_-]+:.*## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build Docker images (optionally SERVICE=...)
	$(COMPOSE) build $(SERVICE)

up: ## Start containers in the background (optionally SERVICE=...)
	$(COMPOSE) up -d $(SERVICE)

down: ## Stop and remove containers, keeping database volumes
	$(COMPOSE) down

logs: ## Follow container logs (optionally SERVICE=...)
	$(COMPOSE) logs -f --tail=100 $(SERVICE)

ps: ## Show container status
	$(COMPOSE) ps

restart: ## Restart containers (optionally SERVICE=...)
	$(COMPOSE) restart $(SERVICE)

clean: ## Remove containers, local images AND database volumes (deletes data)
	$(COMPOSE) down --rmi local --volumes --remove-orphans

rebuild: ## Build Docker images without cache (optionally SERVICE=...)
	$(COMPOSE) build --no-cache $(SERVICE)

nuke: ## Remove containers AND database volumes (deletes data)
	$(COMPOSE) down --volumes --remove-orphans

fresh: ## Rebuild without cache and start containers
	$(MAKE) rebuild
	$(MAKE) up

dev: ## Build and run Docker containers with attached logs
	$(COMPOSE) up --build $(SERVICE)

install: ## Install frontend dependencies from the lockfile
	$(NPM) --prefix frontend ci

frontend-dev: ## Run Vite with live reload (stop the frontend container first)
	$(NPM) --prefix frontend run dev

frontend-build: ## Type-check and build the frontend
	$(NPM) --prefix frontend run build

frontend-preview: ## Preview the built frontend
	$(NPM) --prefix frontend run preview

test: test-go test-frontend ## Run Go and frontend tests

test-go: ## Run tests across all Go workspace modules
	$(GO) test $(GO_PACKAGES)

test-frontend: ## Run frontend tests
	$(NPM) --prefix frontend test

lint: lint-go lint-frontend ## Run Go vet and frontend lint

lint-go: ## Run Go vet across all workspace modules
	$(GO) vet $(GO_PACKAGES)

lint-frontend: ## Lint the frontend
	$(NPM) --prefix frontend run lint

fmt: fmt-go fmt-frontend ## Format Go and frontend files

fmt-go: ## Format Go files across all workspace modules
	$(GO) fmt $(GO_PACKAGES)

fmt-frontend: ## Format frontend files with Prettier
	$(NPM) --prefix frontend run format

fmt-check: ## Check Go and frontend formatting without modifying files
	@files=$$(gofmt -l $(GO_DIRS)) || exit $$?; if [ -n "$$files" ]; then printf '%s\n' "$$files"; exit 1; fi
	$(NPM) --prefix frontend run format:check
