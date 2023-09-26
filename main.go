package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello, world!"))
}

func showNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	//w.Write([]byte("Here comes the note!"))
	fmt.Fprintf(w, "note #%d", id)
}
func createNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// w.WriteHeader(405)
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	w.Write([]byte("Created a new note!"))
}

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
