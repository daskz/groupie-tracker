package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		writeTemplate(w, "NotFound.html", "404 Not Found")
	}
}

func logRequest(r *http.Request) {
	log.Printf("%v %v requested", r.Method, r.URL.Path)
}

func writeTemplate(w http.ResponseWriter, templateName string, data interface{}) {
	t, err := template.ParseFiles(templateName)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func extractQueryID(w http.ResponseWriter, r *http.Request) (int, error) {
	keys, ok := r.URL.Query()["ID"]
	if !ok || len(keys) != 1 {
		return 0, errors.New("Url Param 'ID' is missing")
	}
	key := keys[0]
	id, err := strconv.Atoi(key)
	if err != nil {
		return 0, err
	}

	return id, nil
}
