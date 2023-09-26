package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	// longer URL patterns are matched
	mux.HandleFunc("/", home)         // catch all
	mux.HandleFunc("/note", showNote) // note: `/note/` means `/foo/*`
	mux.HandleFunc("/note/create", createNote)

	log.Println("server started on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
