DB_URL ?= postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
SQL_DIR ?= ./backend/seed/sql

.PHONY: up
up:
	docker compose up

.PHONY: sql
sql:
	sqlc generate -f ./backend/sqlc.yaml 

.PHONY: seed_user_tasks
seed_user_tasks:
	psql $(DB_URL) -f $(SQL_DIR)/user_tasks.sql

.PHONY: seed_tasks
seed_tasks:
	psql $(DB_URL) -f $(SQL_DIR)/tasks.sql

.PHONY: seed_users
seed_users:
	psql $(DB_URL) -f $(SQL_DIR)/users.sql

.PHONY: seed_ayahs
seed_ayahs:
	psql $(DB_URL) -f $(SQL_DIR)/ayahs.sql

.PHONY: seed
seed:
	make seed_ayahs
	make seed_users
	make seed_tasks
	make seed_user_tasks

.PHONY: migrate
migrate:
	atlas schema apply --url "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" --to "file://backend/schema.sql" --dev-url "docker://postgres/17"