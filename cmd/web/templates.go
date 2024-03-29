package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/cwyang/letsgo/pkg/forms"
	"github.com/cwyang/letsgo/pkg/models"
)

type templateData struct {
	CurrentYear     int
	Form            *forms.Form
	Flash           string
	IsAuthenticated bool
	CSRFToken	string
	Note            *models.Note
	Notes           []*models.Note
	User		*models.User
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	// ref time: Mon Jan 2 15:04:05 -0700 MST 2006
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// init a template.FuncMap object & store it in a global var.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// ts, err := template.ParseFiles(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
