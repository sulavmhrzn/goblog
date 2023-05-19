include .envrc
.PHONY:api/run
api/run:
	go run ./cmd/api -dsn=${DB_DSN}

.PHONY:db/migrate/create
db/migrate/create:
	migrate create -seq -ext .sql -dir ./migrations ${name} 

.PHONY:db/migrate/up
db/migrate/up:
	@echo "migrating database..."
	@migrate -path ./migrations -database=${DB_DSN} up