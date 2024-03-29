include .envrc

.PHONY: run psql psql_root create_migrate migrate migrate_rollback confirm help

## run: run the cmd/api application
run:
	go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN}

## build: build the cmd/api application
current_time = $(shell date --iso-8601=seconds)
git_desc = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_desc}'
.PHONY: build
build:
	go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api

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

## production/connect: connect to the production server
production_host_ip = '10.0.0.1'
.PHONY: production/connect
production/connect:
	ssh greenlight@${production_host_ip}
## production/deploy: deploy to the production server
.PHONY: production/deploy
production/deploy:
	rsync -rP --delete ./bin/linux_amd64/api ./migrations greenlight@${production_host_ip}:~
	ssh -t greenlight@{production_host_ip} 'migrate -path ~/migrations -database $$GREENLIGHT_DB_DSN up'
## production/configure: configure production systemd file
production/configure:
	rsync -P ./remote/production/api.service greenlight@${production_host_ip}:~
	ssh -t greenlight@{production_host_ip} '\
	sudo mv ~/api.service /etc/systemd/sysmem \
	&& sudo systemctl enable api \
	&& sudo systemctl restart api \
	'

## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
