run:
	air
	
start-postgres:
	docker run -d --rm --name postgres -e POSTGRES_USER=${DATABASE_USERNAME} -p 5432:5432 -e POSTGRES_PASSWORD=${DATABASE_PASSWORD}  -e POSTGRES_DB=${DATABASE_NAME} -e PGDATA=/var/lib/postgresql/data/pgdata -v ${DATABASE_NAME}:/var/lib/postgresql/data postgres:17.11

stop-postgres:
	docker stop postgres

migrate-up:
	migrate -path ./internal/db/migrations -database "postgresql://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@localhost:5432/${DATABASE_NAME}?sslmode=disable" up

migrate-down:
	migrate -path ./internal/db/migrations -database "postgresql://${DATABASE_USERNAME}:${DATABASE_PASSWORD}@localhost:5432/${DATABASE_NAME}?sslmode=disable" down

migrate-force:
	migrate -path ./internal/db/migrations -database "postgresql://root:secret@localhost:5432/simple_store_go?sslmode=disable" force ${version}