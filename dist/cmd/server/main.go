package main

import (
	"log"

	"github.com/cwyang/letsgo/dist/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":4000")
	log.Fatal(srv.ListenAndServe())
}
