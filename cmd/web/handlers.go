package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/cwyang/letsgo/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	n, err := app.notes.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for _, note := range n {
		fmt.Fprintf(w, "%v\n", note)
	}
	return
	
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) showNote(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

	fmt.Fprintf(w, "%v", n);
}
func (app *application) createNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// w.WriteHeader(405)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

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
