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
	@cd migrations && goose postgres "user=${DB_USER} \
          password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable \
          host=localhost port=${DB_PORT}" up

.PHONY: migrate-down
include .env
migrate-down:
	@echo "migrating down"
	@cd migrations && goose postgres "user=${DB_USER} \
		  password=${DB_PASSWORD} dbname=${DB_NAME} sslmode=disable \
		  host=localhost port=${DB_PORT}" down

.PHONY: run-accrual
run-accrual:
	@echo "running accrual service for mac arm64"
	cmd/accrual/accrual_darwin_arm64

.PHONY: run-docker
run-docker:
	@echo "running docker containers"
	@docker-compose -f docker/compose.yaml --env-file .env up -d

.PHONY: stop-docker
stop-docker:
	@echo "stopping docker containers"
	@docker-compose -f docker/compose.yaml down
