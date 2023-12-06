package services

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Artist struct {
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	FirstAlbum   string   `json:"firstAlbum"`
	Id           int      `json:"id"`
	CreationDate int      `json:"creationDate"`
	Members      []string `json:"members"`
}

func FetchArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var artists []Artist
	err = json.Unmarshal(body, &artists)
	return artists, err
}
 