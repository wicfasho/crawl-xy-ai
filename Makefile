# Makefile

# Load environment variables from .env file
include ./envs/global.env
export

# Default title if not provided
title ?= default_migration_title
migration_folder_path ?= database/migrations
mn ?= #migration number

# ANSI escape code for red text
RED := \033[0;31m
# ANSI escape code to reset text color
RESET := \033[0m

.PHONY: create-migration
create-migration:
	@if [ "$(title)" = "default_migration_title" ]; then \
		printf "$(RED)No title was given. Default title will be used$(RESET)\n"; \
	fi
	migrate create -ext sql -dir $(migration_folder_path) -seq $(title)

.PHONY: migrate-up
migrate-up:
	@if [ "$(mn)" = "" ]; then \
		printf "$(RED)No migration number was given. Last migration will be used$(RESET)\n"; \
	fi
	migrate -database "$(DATABASE_URL)" -path $(migration_folder_path) -verbose up ${mn}

.PHONY: migrate-down
migrate-down:
	@if [ "$(mn)" = "" ]; then \
		printf "$(RED)No migration number was given. Last migration will be used$(RESET)\n"; \
	fi
	migrate -database "$(DATABASE_URL)" -path $(migration_folder_path) -verbose down ${mn}

.PHONY: print-vars
print-vars:
	@if [ "$(title)" = "default_migration_title" ]; then \
		printf "$(RED)No title was given. Default title will be used$(RESET)\n"; \
	fi
	@echo "DATABASE_URL: $(DATABASE_URL)"
	@echo "TITLE: $(title)"
	@echo "MIGRATION_FOLDER_PATH: $(migration_folder_path)"

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: ollama-create
ollama-create:
	ollama create crawlxy.ai -f ./llms/ollama/Modelfile

