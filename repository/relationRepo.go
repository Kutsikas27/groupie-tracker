package repository

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/model"
	"net/http"
)

func GetConcertData() (model.Relation, error) {

	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return model.Relation{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return model.Relation{}, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var concertData model.Relation
	err = json.NewDecoder(resp.Body).Decode(&concertData)
	if err != nil {
		return model.Relation{}, err
	}

	return concertData, nil
}
