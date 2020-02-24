package main

import (
	"encoding/json"
	"fmt"
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

type ArtistData struct {
	Artist         Artist
	DatesLocations map[string]interface{}
}

var _artists = []Artist{}
var _relation = &Relation{}

func main() {
	syncData("https://groupietrackers.herokuapp.com/api/artists", &_artists)
	syncData("https://groupietrackers.herokuapp.com/api/relation", _relation)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/artists", artistIndexHandler)

	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	t.Execute(w, _artists)
}

func artistIndexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Artist Index Requested")
	keys, ok := r.URL.Query()["ID"]
	log.Println(keys)
	if !ok || len(keys) != 1 {
		log.Println("Url Param 'key' is missing ")
		return
	}
	key := keys[0]
	id, err := strconv.Atoi(key)
	if err != nil {
		log.Println(err.Error())
	}

	var artist = filterArtistByID(_artists, id)

	var data = ArtistData{Artist: *artist}
	var dates = filterRelationByID(_relation.Index, id)
	if dates != nil {
		data.DatesLocations = make(map[string]interface{})
		for key, value := range dates.DatesLocations {
			var locationName = strings.ReplaceAll(key, "_", " ")
			locationName = strings.ReplaceAll(locationName, "-", " - ")
			locationName = strings.ToUpper(locationName)
			data.DatesLocations[locationName] = value
		}
	}

	if artist != nil {
		t, err := template.ParseFiles("artist.html")
		if err != nil {
			log.Println(err.Error())
			return
		}

		err = t.Execute(w, data)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func syncData(api string, data interface{}) {
	log.Println("Started synchronization api " + api)
	res, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(bodyBytes, &data)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	log.Println("Completed synchronization api " + api)
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
