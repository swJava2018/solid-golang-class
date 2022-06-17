.PHONEY: up, down

up:
	docker compose -f docker-compose.yaml -f docker-compose-dependencies.yaml up -d
down:
	docker compose -f docker-compose.yaml -f docker-compose-dependencies.yaml down