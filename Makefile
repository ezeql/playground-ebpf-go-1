run:
	docker compose up --build -d
	@echo "Grafana: http://localhost:3300"
	docker compose logs -f obi

down:
	docker compose down

bpftrace:
	docker compose run --rm bpftrace

uprobe:
	docker compose run --rm --build uprobe
