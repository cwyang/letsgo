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
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{
				Form: form,
			})
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.session.Put(r, "authenticatedUserID", id)

	path := app.session.PopString(r, "orgURL")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
	}
	// redirect to the relevant page
	http.Redirect(w, r, "/note/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Pong"))
}

func (app *application) showAbout(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "about.page.tmpl", &templateData{})
}
func (app *application) userProfile(w http.ResponseWriter, r *http.Request) {
	exists := app.session.Exists(r, "authenticatedUserID")
	if !exists {
		app.serverError(w, models.ErrInvalidCredentials)
		return
	}
	u, err := app.users.Get(app.session.GetInt(r, "authenticatedUserID"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "profile.page.tmpl", &templateData{
		User: u,
	})
}
func (app *application) changePassForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "password.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
func (app *application) changePass(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("old", "new1", "new2")
	form.MinLength("new1", 10)
	if form.Get("new1") != form.Get("new2") {
		form.Errors.Add("new2", "Passwords don't match")
	}
	
	if !form.Valid() {
		app.render(w, r, "password.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	err = app.users.ChangePassword(app.session.GetInt(r, "authenticatedUserID"),
		form.Get("old"),
		form.Get("new"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("old", "current password is incorrect")
			app.render(w, r, "password.page.tmpl", &templateData{
				Form: form,
			})
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.session.Put(r, "flash", "Your password has been updated!")

	// redirect to the relevant page
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}
