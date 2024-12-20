package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rishavch2104/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serveError(w, r, err)
		return
	}
	data := app.newTemplateData()
	data.Snippets = snippets
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serveError(w, r, err)
		}
		return
	}
	data := app.newTemplateData()
	data.Snippet = s
	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating snippet"))
}
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	title := "Test Title 1"
	content := "Content of test title 1"
	expires := 4
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serveError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "created snippet with id %d", id)
}
