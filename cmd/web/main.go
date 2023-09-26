package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP listen port")
	flag.Parse()
	
	mux := http.NewServeMux()
	// longer URL patterns are matched
	mux.HandleFunc("/", home)         // catch all
	mux.HandleFunc("/note", showNote) // note: `/note/` means `/foo/*`
	mux.HandleFunc("/note/create", createNote)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("server started on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
