.PHONY: build up down logs restart clean rebuild fresh dev

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

restart: down up

clean:
	docker compose down --rmi local --volumes --remove-orphans

rebuild:
	docker compose build --no-cache

fresh: rebuild up

dev:
	docker compose build --no-cache
	docker compose up