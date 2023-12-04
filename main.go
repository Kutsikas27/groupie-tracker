package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
)

type Body struct {
	Index []struct {
		ID      int                 `json:"id"`
		Concert map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

type Artist struct {
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	FirstAlbum   string   `json:"firstAlbum"`
	Id           int      `json:"id"`
	CreationDate int      `json:"creationDate"`
	Members      []string `json:"members"`
}

func main() {
	http.HandleFunc("/artist/", artistPageHandler)
	http.HandleFunc("/", mainPageHandler)
	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func artistPageHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path[1:]
	id, _ := strconv.Atoi(path.Base(page))
	artists, _ := fetchArtists()
	artist, _ := findArtistByID(artists, id)
	locationAndDate, _ := getConcertData()
	locationById, _ := findConcertDataByID(locationAndDate, id)
	if err := renderArtistPage(w, artist, locationById); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering HTML: %s", err), http.StatusInternalServerError)
	}
}

func findArtistByID(artists []Artist, id int) (Artist, error) {
	for _, artist := range artists {
		if artist.Id == id {
			return artist, nil
		}
	}
	return Artist{}, fmt.Errorf("Artist with Id %d not found", id)
}
func createStructForDates() {
	type ConcertByDate struct {
		Date     []string
		Location struct {
			City    string
			Country string
		}
	}

}

func findConcertDataByID(concertData Body, id int) (map[string][]string, error) {
	for _, entry := range concertData.Index {
		if entry.ID == id {
			return entry.Concert, nil
		}
	}
	return nil, fmt.Errorf("id %d not found", id)
}

//	func extractCities(str string) City {
//		cities := strings.Split(str, "-")
//		cities2 := strings.Join(cities, " ")
//		return strings.Title(cities2)
//	}
func renderArtistPage(w http.ResponseWriter, artist Artist, locationAndDate map[string][]string) error {
	data := struct {
		Artist         Artist
		DatesLocations map[string][]string
	}{
		Artist:         artist,
		DatesLocations: locationAndDate,
	}

	htmlTemplate, err := template.ParseFiles("./static/artist.html")
	if err != nil {
		return err
	}

	return htmlTemplate.Execute(w, data)
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := fetchArtists()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching artists: %s", err), http.StatusInternalServerError)
		return
	}

	if err := renderMainPage(w, artists); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering HTML: %s", err), http.StatusInternalServerError)
	}
}

func fetchArtists() ([]Artist, error) {
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

func renderMainPage(w http.ResponseWriter, artists []Artist) error {
	htmlTemplate, err := template.ParseFiles("./static/index.html")
	if err != nil {
		return err
	}

	return htmlTemplate.Execute(w, artists)
}
func getConcertData() (Body, error) {

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
