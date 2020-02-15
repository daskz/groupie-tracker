package main

import (
	"log"
	"net/http"
	"text/template"
)

func main() {
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	page := struct {
		Title string
		Body  string
	}{
		Title: "this is Title",
		Body:  "this is Body",
	}

	t.Execute(w, page)
}
