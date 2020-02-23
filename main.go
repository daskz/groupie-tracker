package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type Artist struct {
	ID           int
	Name         string
	CreationDate int
	FirstAlbum   string
	Image        string
	Members      []string
}

type Relation struct {
	Index []ArtistRelation
}

type ArtistRelation struct {
	ID             int
	DatesLocations map[string]interface{}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/artists/", artistsHandler)
	http.HandleFunc("/api/concerts/", concertsHandler)

	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	idUrl := strings.TrimPrefix(r.URL.Path, "/")

	var fileName = "index.html"
	_, err := strconv.Atoi(idUrl)
	log.Println(idUrl)
	if err == nil {
		fileName = "artist.html"
		//json.NewEncoder(w).Encode(filterRelationByID(relation.Index, id))
	}

	t, _ := template.ParseFiles(fileName)

	t.Execute(w, nil)
}

func concertsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("/api/concerts requested")
	res, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == http.StatusOK {
		relation := Relation{}
		bodyBytes, err := ioutil.ReadAll(res.Body)
		//log.Println(string(bodyBytes))
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(bodyBytes, &relation)
		if err != nil {
			log.Fatal(err.Error())
		}
		//1og.Println(relation)

		idUrl := strings.TrimPrefix(r.URL.Path, "/api/concerts/")
		id, err := strconv.Atoi(idUrl)
		log.Println(idUrl)
		if idUrl != "" && err == nil {
			// log.Fatal(err.Error())
			json.NewEncoder(w).Encode(filterRelationByID(relation.Index, id))
		}
	}
}

func artistsHandler(w http.ResponseWriter, r *http.Request) {
	var artists []Artist
	log.Println("/api/artists requested")
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

	idUrl := strings.TrimPrefix(r.URL.Path, "/api/artists/")
	id, err := strconv.Atoi(idUrl)
	log.Println(idUrl)
	if idUrl != "" && err == nil {
		// log.Fatal(err.Error())
		json.NewEncoder(w).Encode(filterArtistByID(artists, id))
		return
	}
	json.NewEncoder(w).Encode(artists)
}

func filterArtistByID(artists []Artist, id int) *Artist {
	for _, item := range artists {
		if item.ID == id {
			return &item
		}
	}
	return nil
}

func filterRelationByID(relations []ArtistRelation, id int) *ArtistRelation {
	for _, item := range relations {
		if item.ID == id {
			return &item
		}
	}
	return nil
}
