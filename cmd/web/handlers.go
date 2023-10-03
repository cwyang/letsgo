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
	// redirect to the relevant page
	http.Redirect(w, r, fmt.Sprintf("/note/%d", id), http.StatusSeeOther)
}
