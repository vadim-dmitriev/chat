package main

import (
	"io/ioutil"
	"net/http"
)

type page string

func (t page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pageContent, err := ioutil.ReadFile(string(t))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(pageContent)
}
