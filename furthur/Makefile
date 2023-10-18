GREENLIGHT_DB_DSN="postgres://greenlight:password@localhost/greenlight"
all:
	go run ./cmd/api
psql:
	echo psql --host=localhost --dbname=greenlight --username=greenlight
	psql ${GREENLIGHT_DB_DSN}
psql_root:
	sudo -u postgres psql --dbname=greenlight
create_migrate:
	echo migrate create -seq -ext=.sql -dir=./migrations migration_name
migrate:
	migrate -path=./migrations -database=${GREENLIGHT_DB_DSN} up
migrate_rollback:
	migrate -path=./migrations -database=${GREENLIGHT_DB_DSN} down
