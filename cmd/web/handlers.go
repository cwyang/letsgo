package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	w.Write([]byte("Create a new note!"))
}
func (app *application) createNote(w http.ResponseWriter, r *http.Request) {
	title := "Quick Brown Fox"
	content := "Jumps over the\nLittle Lazy Dog!"
	expires := "7"

	id, err := app.notes.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// redirect to the relevant page
	http.Redirect(w, r, fmt.Sprintf("/note?id=%d", id), http.StatusSeeOther)
}
