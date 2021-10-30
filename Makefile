#!make
include .env
export $(shell sed 's/=.*//' .env)

SHELL := /bin/bash

MODULE = $(shell go list -m)
PID_FILE := './.pid'
FSWATCH_FILE := './fswatch.cfg'
PACKAGES := $(shell go list ./... | grep -v /vendor/)

CONFIG_FILE ?= .env
APP_DSN ?= $(shell sed -n 's/^dsn:[[:space:]]*"\(.*\)"/\1/p' $(CONFIG_FILE))
MIGRATE := docker run --rm -v $(shell pwd)/migrations:/migrations --network host --user $(id -u):$(id -g) migrate/migrate -path=/migrations/ -database "$$DB_PATH"
MIGRATE_TEST_DB := docker run --rm -v $(shell pwd)/migrations:/migrations --network host --user $(id -u):$(id -g) migrate/migrate -path=/migrations/ -database "$$TEST_DB_PATH"
MIGRATE_CREATE := docker run --rm -v $(shell pwd)/migrations:/migrations --network host --user $(shell id -u):$(shell id -g) migrate/migrate create --seq -ext sql -dir /migrations/
CWD := $(shell pwd)

.PHONY: default
default: help

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build:
	CGO_ENABLED=0 go build -o weeklyprogramme

.PHONY: clean
clean:
	rm -rf weeklyprogramme

.PHONY: db-start
db-start: ## start the database
	@mkdir -p testdata/postgres
	docker run --rm --net host --name weeklyprogrammedb -d -v $(shell pwd)/testdata:/testdata \
		-v $(shell pwd)/testdata/postgres:/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD=password -e POSTGRES_DB=weeklyprogrammedb -e POSTGRES_USER=admin -d postgres:12.3

.PHONY: db-stop
db-stop: ## stop the database
	docker stop weeklyprogrammedb

.PHONY: db-login
db-login: ## login to the database
	docker exec -it weeklyprogrammedb psql -U admin -d weeklyprogrammedb

.PHONY: test-db-start
test-db-start: ## start the test database
	@mkdir -p testdata/postgres-test
	docker run --rm --name weeklyprogrammedb_test -p 5433:5432 -d -v $(shell pwd)/testdata:/testdata \
		-v $(shell pwd)/testdata/postgres-test:/var/lib/postgresql/data \
		-e POSTGRES_PASSWORD=password -e POSTGRES_DB=weeklyprogrammedb_test -e POSTGRES_USER=admin -d postgres:12.3

.PHONY: test-db-stop
test-db-stop: ## stop the database server
	docker stop weeklyprogrammedb_test

.PHONY: test-db-login
test-db-login: ## login to the database
	docker exec -it weeklyprogrammedb_test psql -U admin -d weeklyprogrammedb_test

.PHONY: migrate
migrate: ## run all new database migrations
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## revert database to the last migration step
	@echo "Reverting database to the last migration step..."
	@$(MIGRATE) down 1

.PHONY: migrate-new
migrate-new: ## create a new database migration
	@read -p "Enter the name of the new migration: " name; \
	$(MIGRATE_CREATE) $${name}

.PHONY: test-migrate
test-migrate: ## run all new database migrations
	@echo "Running all new database migrations..."
	@$(MIGRATE_TEST_DB) up