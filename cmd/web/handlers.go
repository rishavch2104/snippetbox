package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/rishavch2104/snippetbox/internal/models"
)

func (app application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serveError(w, r, err)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serveError(w, r, err)
	}
}

func (app application) snippetView(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprintf(w, "Title is %s, Content is %s, Created on %s, Expires on %s", s.Title, s.Content, s.Created, s.Expires)
}

func (app application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating snippet"))
}
func (app application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
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
