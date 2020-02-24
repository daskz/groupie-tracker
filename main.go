package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	go syncData("https://groupietrackers.herokuapp.com/api/artists", &_artists)
	go syncData("https://groupietrackers.herokuapp.com/api/relation", _relation)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/artists", artistIndexHandler)

	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if r.URL.Path != "/" || r.Method != "GET" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	writeTemplate(w, "index.html", _artists)
}

func artistIndexHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if r.Method != "GET" {
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	id, err := extractQueryID(w, r)
	if err != nil {
		log.Println(err.Error())
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	data, err := getArtistByID(id)
	if err != nil {
		log.Println(err.Error())
		errorHandler(w, r, http.StatusNotFound)
		return
	}

	writeTemplate(w, "artist.html", data)
}
