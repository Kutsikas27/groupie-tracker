package controller

import (
	"fmt"
	"groupie-tracker/model"
	"groupie-tracker/repository"
	"groupie-tracker/service"
	"net/http"
	"path"
	"strconv"
	"text/template"
)

func ArtistPageHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path[1:]
	id, _ := strconv.Atoi(path.Base(page))
	artists, _ := repository.FetchArtists()
	artist, _ := service.FindArtistByID(artists, id)

	locationAndDate, _ := repository.GetConcertData()
	locationById, _ := service.FindConcertDataByID(locationAndDate, id)
	data := struct {
		Artist   model.Artist
		Concerts []model.Concert
	}{
		Artist:   artist,
		Concerts: locationById,
	}

	if err := renderPage(w, data, "./templates/artist.html"); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering HTML: %s", err), http.StatusInternalServerError)

	}
}
func renderPage(w http.ResponseWriter, data interface{}, templatePath string) error {
	htmlTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	return htmlTemplate.Execute(w, data)
}

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := repository.FetchArtists()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching artists: %s", err), http.StatusInternalServerError)
		return
	}

	if err := renderPage(w, artists, "./templates/index.html"); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering HTML: %s", err), http.StatusInternalServerError)
	}
}
