package service

import (
	"fmt"
	"groupie-tracker/model"
)

func FindArtistByID(artists []model.Artist, id int) (model.Artist, error) {
	for _, artist := range artists {
		if artist.Id == id {
			return artist, nil
		}
	}
	return model.Artist{}, fmt.Errorf("artist with id %d not found", id)
}
