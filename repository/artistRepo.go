package repository

import (
	"encoding/json"
	"groupie-tracker/model"
	"io/ioutil"
	"net/http"
)

func FetchArtists() ([]model.Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var artists []model.Artist
	err = json.Unmarshal(body, &artists)
	return artists, err
}
