package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type Artist struct {
	ID           int
	Name         string
	CreationDate int
	FirstAlbum   string
	Image        string
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/artists", artistsHandler)

	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")

	var fileName = "index.html"
	if id != "" {
		fileName = "artist.html"
	}
	t, _ := template.ParseFiles(fileName)

	t.Execute(w, nil)
}

func artistsHandler(w http.ResponseWriter, r *http.Request) {
	var artists []Artist

	res, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(bodyBytes, &artists)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	json.NewEncoder(w).Encode(artists)
}
