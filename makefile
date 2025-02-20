MIGRATIONS_PATH = ./cmd/migrate/migrations
DB_MIGRATOR_ADDR = postgres://myusername:mypassword1234@localhost/social?sslmode=disable

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migration-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_MIGRATOR_ADDR) up

.PHONY: migrate-down
migration-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_MIGRATOR_ADDR) down


.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go

start-app:
	@air


.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal,types && swag fmt