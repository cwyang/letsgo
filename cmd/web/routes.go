package main

import (
	"net/http"

	"github.com/justinas/alice"	// for chaining middleware
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	
	mux := http.NewServeMux()
	// longer URL patterns are matched
	mux.HandleFunc("/", app.home)         // catch all
	mux.HandleFunc("/note", app.showNote) // note: `/note/` means `/foo/*`
	mux.HandleFunc("/note/create", app.createNote)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
