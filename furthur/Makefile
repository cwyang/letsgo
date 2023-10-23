include .envrc

.PHONY: run psql psql_root create_migrate migrate migrate_rollback confirm help

## run: run the cmd/api application
run:
	go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN}

## build: build the cmd/api application
.PHONY: build
build:
	go build -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

## psql: connect to the database using psql
psql:
	echo psql --host=localhost --dbname=greenlight --username=greenlight
	psql ${GREENLIGHT_DB_DSN}

## psql_root: connect to the postgreSQL with admin privelige 
psql_root:
	sudo -u postgres psql --dbname=greenlight

## create_migrate: create migration files
create_migrate: confirm
	echo migrate create -seq -ext=.sql -dir=./migrations migration_name
	echo make create_migrate name=migration_name
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## migrate: migrate up
migrate:
	migrate -path=./migrations -database=${GREENLIGHT_DB_DSN} up
## migrate_rollback: migrate down
migrate_rollback:
	migrate -path=./migrations -database=${GREENLIGHT_DB_DSN} down

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit: vendor
	@echo 'Formatting code'
	go fmt ./...
	@echo 'Vetting code'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies'
	go mod vendor

confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
