package main

import (
	"net/http"

	"github.com/bmizerany/pat"  // router
	"github.com/justinas/alice" // for chaining middleware
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()

	// Pat matches pattern in the registration order
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/note/create", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.createNoteForm))
	mux.Post("/note/create", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.createNote))
	mux.Get("/note/:id", dynamicMiddleware.ThenFunc(app.showNote))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.logoutUser))

	mux.Get("/ping", http.HandlerFunc(ping))
	mux.Get("/about", dynamicMiddleware.ThenFunc(app.showAbout))
	mux.Get("/user/profile", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.userProfile))
	mux.Get("/user/changepass", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.changePassForm))
	mux.Post("/user/changepass", dynamicMiddleware.Append(app.requireAuth).ThenFunc(app.changePass))
	
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
