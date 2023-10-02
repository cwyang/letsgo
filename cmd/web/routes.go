package main

import (
	"net/http"

	"github.com/bmizerany/pat"	// router
	"github.com/justinas/alice"	// for chaining middleware
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	
	mux := pat.New()

	// Pat matches pattern in the registration order
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/note/create", http.HandlerFunc(app.createNoteForm))
	mux.Post("/note/create", http.HandlerFunc(app.createNote))
	mux.Get("/note/:id", http.HandlerFunc(app.showNote))


	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
