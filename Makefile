.PHONY: test
test:
	@echo "running tests"
	@go test ./... -coverprofile=coverage.out --race
	@go tool cover -func=coverage.out | grep total
	@rm -f coverage.out

.PHONY: migrate-new
migrate-new:
	@echo "new migration"
	@cd migrations && goose create $(name) sql

.PHONY: migrate-up
include .env
migrate-up:
	@echo "migrating up"
	@cd migrations && goose postgres "user=${POSTGRES_USER} \
          password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable \
          host=localhost port=${POSTGRES_PORT}" up

.PHONY: migrate-down
include .env
migrate-down:
	@echo "migrating down"
	@cd migrations && goose postgres "user=${POSTGRES_USER} \
		  password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable \
		  host=localhost port=${POSTGRES_PORT}" down

.PHONY: run-accrual
run-accrual:
	@echo "running accrual service for mac arm64"
	cmd/accrual/accrual_darwin_arm64