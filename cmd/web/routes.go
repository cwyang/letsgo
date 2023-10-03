package main

import (
	"net/http"

	"github.com/bmizerany/pat"	// router
	"github.com/justinas/alice"	// for chaining middleware
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable)
	
	mux := pat.New()

	// Pat matches pattern in the registration order
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/note/create", dynamicMiddleware.ThenFunc(app.createNoteForm))
	mux.Post("/note/create", dynamicMiddleware.ThenFunc(app.createNote))
	mux.Get("/note/:id", dynamicMiddleware.ThenFunc(app.showNote))


	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
