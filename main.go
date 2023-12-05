package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"
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

type ConcertByDate struct {
	Date     string
	Location struct {
		City    string
		Country string
	}
}

func (c ConcertByDate) Less(other ConcertByDate) bool {
	date1Parts := strings.Split(c.Date, ".")
	date2Parts := strings.Split(other.Date, ".")

	year1, _ := strconv.Atoi(date1Parts[2])
	month1, _ := strconv.Atoi(date1Parts[1])
	day1, _ := strconv.Atoi(date1Parts[0])

	year2, _ := strconv.Atoi(date2Parts[2])
	month2, _ := strconv.Atoi(date2Parts[1])
	day2, _ := strconv.Atoi(date2Parts[0])

	if year1 != year2 {
		return year1 < year2
	}
	if month1 != month2 {
		return month1 < month2
	}
	return day1 < day2
}

func main() {
	fileServer := http.FileServer(http.Dir("./images"))
	http.Handle("/images/", http.StripPrefix("/images/", fileServer))

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

func findConcertDataByID(concertData Body, id int) ([]ConcertByDate, error) {
	for _, entry := range concertData.Index {
		if entry.ID == id {
			var concerts []ConcertByDate
			for city, dates := range entry.Concert {
				for _, date := range dates {
					city2, country := extractCities(city)
					concert := ConcertByDate{
						Date: normlizeDate(date),
						Location: struct {
							City    string
							Country string
						}{
							City:    city2,
							Country: country,
						},
					}
					concerts = append(concerts, concert)
				}
			}

			sort.Slice(concerts, func(i, j int) bool {
				return concerts[i].Less(concerts[j])
			})

			return concerts, nil
		}
	}
	return nil, fmt.Errorf("id %d not found", id)
}
func normlizeDate(date string) string {
	splitDate := strings.Split(date, "-")
	joinDate := strings.Join(splitDate, ".")

	return joinDate
}

func extractCities(str string) (string, string) {
	cities := strings.Split(str, "-")
	cities2 := strings.Join(cities, " ")
	city := strings.Split(cities2, " ")

	city2 := strings.ToTitle(city[0])
	country := strings.ToTitle(city[1])
	return normlizeString(city2), normlizeString(country)
}
func normlizeString(str string) string {
	splitStr := strings.Split(str, "_")
	return strings.Join(splitStr, " ")
}
func renderArtistPage(w http.ResponseWriter, artist Artist, concertData []ConcertByDate) error {
	data := struct {
		Artist   Artist
		Concerts []ConcertByDate
	}{
		Artist:   artist,
		Concerts: concertData,
	}

	htmlTemplate, err := template.ParseFiles("./templates/artist.html")
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
	htmlTemplate, err := template.ParseFiles("./templates/index.html")
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
