GOCACHE := $(CURDIR)/.cache/go-build
GOMODCACHE := $(CURDIR)/.cache/gomod

.PHONY: run swag test test-integration test-db-up test-db-down test-all seed-start seed-reset generate-entity generate-usecase docker-up docker-down docker-logs docker-logs-tail docker-ps docker-reset restart restart-build db-recreate

run:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go run ./cmd/api

swag:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/api/main.go -o docs/swagger --parseDependency --parseInternal

test:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go test ./...

test-integration:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) DATABASE_DSN_TEST="$${DATABASE_DSN_TEST}" go test ./... -p 1 -count=1 -v

test-db-up:
	docker compose -f docker-compose.test.yml up -d

test-db-down:
	docker compose -f docker-compose.test.yml down -v

test-all:
	@set -e; \
	cleanup() { \
		echo ">> Bajando DB de tests..."; \
		docker compose -f docker-compose.test.yml down -v; \
	}; \
	trap cleanup EXIT; \
	echo ">> Levantando DB de tests..."; \
	docker compose -f docker-compose.test.yml up -d; \
	echo ">> Esperando DB saludable..."; \
	for i in $$(seq 1 45); do \
		status=$$(docker inspect -f '{{.State.Health.Status}}' photogallery_db_test 2>/dev/null || true); \
		if [ "$$status" = "healthy" ]; then \
			echo ">> DB lista"; \
			break; \
		fi; \
		if [ $$i -eq 45 ]; then \
			echo "ERROR: DB de tests no alcanzó estado healthy"; \
			exit 1; \
		fi; \
		sleep 2; \
	done; \
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) DATABASE_DSN_TEST="host=localhost user=admin password=1234 dbname=photogallery_test port=5434 sslmode=disable" bash cmd/scripts/run_pretty_tests.sh .cache/test-all.log

seed-reset:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go run ./cmd/seedreset

seed-start:
	GOCACHE=$(GOCACHE) GOMODCACHE=$(GOMODCACHE) go run ./cmd/seedstart

generate-entity:
	bash cmd/scripts/generate_entity.sh $(name)

generate-usecase:
	bash cmd/scripts/generate_use_case.sh $(entity) $(name)

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f app

docker-logs-tail:
	docker compose logs --tail=120 app

docker-ps:
	docker compose ps

docker-reset:
	docker compose down -v
	docker compose up -d --build

restart:
	docker compose restart app
	docker compose logs --tail=60 app

restart-build:
	docker compose up -d --build app
	docker compose logs --tail=60 app

db-recreate:
	@set -e; \
	echo ">> Eliminando contenedores y volumen de Postgres..."; \
	docker compose down -v; \
	echo ">> Levantando DB, pgAdmin y API..."; \
	docker compose up -d --build; \
	echo ">> Esperando DB saludable..."; \
	for i in $$(seq 1 45); do \
		status=$$(docker inspect -f '{{.State.Health.Status}}' photogallery_db 2>/dev/null || true); \
		if [ "$$status" = "healthy" ]; then \
			echo ">> DB lista"; \
			break; \
		fi; \
		if [ $$i -eq 45 ]; then \
			echo "ERROR: la DB no alcanzó estado healthy"; \
			exit 1; \
		fi; \
		sleep 2; \
	done; \
	echo ">> Listo. GORM ejecuta AutoMigrate al iniciar la API."
