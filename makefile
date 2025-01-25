
run_local:
	docker compose -f infra-docker/docker-compose.local.yml up -d 
	cd frontend && bun run dev --open

rebuild_local_backend:
	docker compose -f infra-docker/docker-compose.local.yml down
	docker compose -f infra-docker/docker-compose.local.yml up -d --build
	