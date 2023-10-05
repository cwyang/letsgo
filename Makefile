SOURCES = cmd/web/main.go \
	cmd/web/handlers.go

all:	$(SOURCES)
	go run ./cmd/web

test:
	go test -v ./cmd/web

gencert:
	(cd tls; go run /usr/lib/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost)
