package server

import (
	tmpl "html/template"
	"net/http"
	"path/filepath"
)

type template struct {
	*tmpl.Template
}

func newTemplate(path string) (*template, error) {
	t, err := tmpl.New(filepath.Base(path)).Delims("[[", "]]").ParseFiles(path)
	if err != nil {
		return nil, err
	}

	return &template{t}, nil
}

func (t template) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := t.Execute(w, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
