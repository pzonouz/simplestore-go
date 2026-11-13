# Force serial execution so `make -j up` can't launch deps in parallel
.NOTPARALLEL:

# ============================================================================
# Configuration
# ============================================================================
NETWORK           = search-net

# Host-facing URL — used by migrate-* (runs on your machine)
POSTGRES_URL  = postgresql://$(DATABASE_USER):$(DATABASE_PASSWORD)@localhost:5432/$(DATABASE_NAME)?sslmode=disable

# Container-facing URL — used by PGSync (runs inside the Docker network)
PGSYNC_PG_URL = postgresql://$(DATABASE_USER):$(DATABASE_PASSWORD)@postgres:5432/$(DATABASE_NAME)?sslmode=disable

# ============================================================================
# Phony targets
# ============================================================================
.PHONY: dev build run \
        network \
        up down status \
        start-postgres stop-postgres \
        start-opensearch stop-opensearch \
        start-redis stop-redis \
        start-pgsync stop-pgsync logs-pgsync \
        migrate-up migrate-down migrate-force

# ============================================================================
# Application
# ============================================================================
dev:
	air

build:
	go build -o goserver ./cmd/api/.

run: build
	./goserver

# ============================================================================
# Network
# ============================================================================
network:
	@docker network inspect $(NETWORK) >/dev/null 2>&1 || \
		(echo "creating network $(NETWORK)" && docker network create $(NETWORK))

# ============================================================================
# PostgreSQL
# ============================================================================
start-postgres: network
	@docker container inspect -f '{{.State.Running}}' postgres 2>/dev/null | grep -q true || \
		docker run -d --rm --name postgres \
			--network $(NETWORK) \
			-p 5432:5432 \
			-e POSTGRES_USER=$(DATABASE_USER) \
			-e POSTGRES_PASSWORD=${DATABASE_PASSWORD} \
			-e POSTGRES_DB=$(DATABASE_NAME) \
			-e PGDATA=/var/lib/postgresql/data/pgdata \
			-v $(DATABASE_NAME):/var/lib/postgresql/data \
			postgres:17.11 \
			-c wal_level=logical \
			-c max_replication_slots=4 \
			-c max_wal_senders=4
	@echo "waiting for postgres to accept connections..."
	@for i in $$(seq 1 30); do \
		docker exec postgres pg_isready -U $(DATABASE_USER) -d $(DATABASE_NAME) >/dev/null 2>&1 && break; \
		sleep 1; \
	done
	@docker exec postgres pg_isready -U $(DATABASE_USER) -d $(DATABASE_NAME) >/dev/null 2>&1 || \
		(echo "postgres failed to become ready" && exit 1)

stop-postgres:
	-docker stop postgres 2>/dev/null || true

# ============================================================================
# OpenSearch
# ============================================================================
start-opensearch: network
	@docker container inspect -f '{{.State.Running}}' opensearch 2>/dev/null | grep -q true || \
		docker run -d --rm --name opensearch \
			--network $(NETWORK) \
			-p 9200:9200 \
			-p 9600:9600 \
			-e discovery.type=single-node \
			-e DISABLE_SECURITY_PLUGIN=true \
			opensearchproject/opensearch:latest
	@echo "waiting for opensearch to accept connections..."
	@for i in $$(seq 1 90); do \
		curl -sf http://localhost:9200 >/dev/null 2>&1 && break; \
		sleep 1; \
	done
	@curl -sf http://localhost:9200 >/dev/null 2>&1 || \
		(echo "opensearch failed to become ready" && exit 1)

stop-opensearch:
	-docker stop opensearch 2>/dev/null || true

# ============================================================================
# Redis (required by PGSync for replication checkpoints)
# ============================================================================
start-redis: network
	@docker container inspect -f '{{.State.Running}}' redis 2>/dev/null | grep -q true || \
		docker run -d --rm --name redis \
			--network $(NETWORK) \
			-p 6379:6379 \
			redis:7-alpine
	@for i in $$(seq 1 20); do \
		docker exec redis redis-cli ping 2>/dev/null | grep -q PONG && break; \
		sleep 0.5; \
	done
	@docker exec redis redis-cli ping 2>/dev/null | grep -q PONG || \
		(echo "redis failed to become ready" && exit 1)

stop-redis:
	-docker stop redis 2>/dev/null || true

# ============================================================================
# PGSync (CDC: PostgreSQL -> OpenSearch)
# ============================================================================
start-pgsync: start-postgres start-opensearch start-redis
	@test -f "$(CURDIR)/schema.json" || \
		(echo "schema.json not found in $(CURDIR)" && exit 1)
	@mkdir -p "$(CURDIR)/checkpoints"
	@docker container inspect -f '{{.State.Running}}' pgsync 2>/dev/null | grep -q true || \
		( \
			echo "running bootstrap (creates index + replication slot)..."; \
			docker run --rm \
				--entrypoint bootstrap \
				--network $(NETWORK) \
				-e PG_URL=$(PGSYNC_PG_URL) \
				-e ELASTICSEARCH_URL=http://opensearch:9200 \
				-e OPENSEARCH=true \
				-e REDIS_HOST=redis \
				-e CHECKPOINT_PATH=/app/checkpoints \
				-v "$(CURDIR)/schema.json:/app/schema.json" \
				-v "$(CURDIR)/checkpoints:/app/checkpoints" \
				toluaina1/pgsync:latest --config schema.json || exit 1; \
			echo "starting pgsync daemon..."; \
			docker run -d --name pgsync \
				--network $(NETWORK) \
				-e PG_URL=$(PGSYNC_PG_URL) \
				-e ELASTICSEARCH_URL=http://opensearch:9200 \
				-e OPENSEARCH=true \
				-e REDIS_HOST=redis \
				-e CHECKPOINT_PATH=/app/checkpoints \
				-v "$(CURDIR)/schema.json:/app/schema.json" \
				-v "$(CURDIR)/checkpoints:/app/checkpoints" \
				toluaina1/pgsync:latest --config schema.json --daemon \
		)

stop-pgsync:
	-docker stop pgsync 2>/dev/null || true
	-docker rm pgsync 2>/dev/null || true

logs-pgsync:
	@docker container inspect pgsync >/dev/null 2>&1 \
		&& docker logs -f pgsync \
		|| echo "pgsync container does not exist — run 'make start-pgsync' first"

# ============================================================================
# Aggregate targets
# ============================================================================
up: start-postgres start-opensearch start-redis start-pgsync
	@echo ""
	@echo "stack is up:"
	@echo "  PostgreSQL : localhost:5432"
	@echo "  OpenSearch : http://localhost:9200"
	@echo "  Redis      : localhost:6379"
	@echo "  PGSync     : streaming  (make logs-pgsync to tail)"

down: stop-pgsync stop-postgres stop-opensearch stop-redis

status:
	@echo "containers on $(NETWORK):"
	@docker ps --filter network=$(NETWORK) --format '  {{.Names}}\t{{.Status}}'

# ============================================================================
# Migrations
# ============================================================================
migrate-up: start-postgres
	migrate -path ./internal/db/migrations -database "$(POSTGRES_URL)" up

migrate-down: start-postgres
	migrate -path ./internal/db/migrations -database "$(POSTGRES_URL)" down

migrate-force: start-postgres
	@test -n "$(version)" || (echo "usage: make migrate-force version=N" && exit 1)
	migrate -path ./internal/db/migrations -database "$(POSTGRES_URL)" force $(version)
