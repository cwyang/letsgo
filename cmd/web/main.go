package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// application-wide dependency
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP listen port")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}
	mux := http.NewServeMux()
	// longer URL patterns are matched
	mux.HandleFunc("/", app.home)         // catch all
	mux.HandleFunc("/note", app.showNote) // note: `/note/` means `/foo/*`
	mux.HandleFunc("/note/create", app.createNote)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// err := http.ListenAndServe(*addr, mux)
	// redirect http.Server log to standard error
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	infoLog.Printf("server started on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
