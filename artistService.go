package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var _artists = []artist{}
var _relation = &relation{}

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

func getArtistByID(id int) (*artistData, error) {
	artist := filterArtistByID(_artists, id)
	if artist == nil {
		return nil, errors.New("Artist Not Found")
	}

	var data = artistData{Artist: *artist}
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

	return &data, nil
}

func filterArtistByID(artists []artist, id int) *artist {
	for _, item := range artists {
		if item.ID == id {
			return &item
		}
	}
	return nil
}

func filterRelationByID(relations []artistRelation, id int) *artistRelation {
	for _, item := range relations {
		if item.ID == id {
			return &item
		}
	}
	return nil
}
