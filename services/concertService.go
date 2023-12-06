package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Body struct {
	Index []struct {
		ID      int                 `json:"id"`
		Concert map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

func GetConcertData() (Body, error) {

	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return Body{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Body{}, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var concertData Body
	err = json.NewDecoder(resp.Body).Decode(&concertData)
	if err != nil {
		return Body{}, err
	}

	return concertData, nil
}
