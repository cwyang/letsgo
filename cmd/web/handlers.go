package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cwyang/letsgo/pkg/forms"
	"github.com/cwyang/letsgo/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	n, err := app.notes.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Notes: n,
	})
}

func (app *application) showNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	n, err := app.notes.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Note: n,
	})
}
func (app *application) createNoteForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) createNote(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.notes.Insert(
		form.Get("title"),
		form.Get("content"),
		form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Note created successfully!")

	// redirect to the relevant page
	http.Redirect(w, r, fmt.Sprintf("/note/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	}
	err = app.users.Insert(form.Get("name"),
		form.Get("email"),
		form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Duplicate Email")
			app.render(w, r, "signup.page.tmpl", &templateData{
				Form: form,
			})
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.session.Put(r, "flash", "You are signed up. Please log in.")

	// redirect to the relevant page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
}
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
}
