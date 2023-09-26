SOURCES = cmd/web/main.go \
	cmd/web/handlers.go

all:	$(SOURCES)
	go run ./cmd/web
