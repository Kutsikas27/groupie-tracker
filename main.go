package main

import (
	"fmt"
	"groupie-tracker/services"
	"html/template"
	"log"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"
)

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
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/artist/", artistPageHandler)
	http.HandleFunc("/", mainPageHandler)
	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func artistPageHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Path[1:]
	id, _ := strconv.Atoi(path.Base(page))
	artists, _ := services.FetchArtists()
	artist, _ := findArtistByID(artists, id)

	locationAndDate, _ := services.GetConcertData()
	locationById, _ := findConcertDataByID(locationAndDate, id)
	data := struct {
		Artist   services.Artist
		Concerts []ConcertByDate
	}{
		Artist:   artist,
		Concerts: locationById,
	}

	if err := renderPage(w, data, "./templates/artist.html"); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering HTML: %s", err), http.StatusInternalServerError)

	}
}

func findArtistByID(artists []services.Artist, id int) (services.Artist, error) {
	for _, artist := range artists {
		if artist.Id == id {
			return artist, nil
		}
	}
	return services.Artist{}, fmt.Errorf("artist with id %d not found", id)
}

func findConcertDataByID(concertData services.Body, id int) ([]ConcertByDate, error) {
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

func renderPage(w http.ResponseWriter, data interface{}, templatePath string) error {
	htmlTemplate, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}
	return htmlTemplate.Execute(w, data)
}

func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := services.FetchArtists()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching artists: %s", err), http.StatusInternalServerError)
		return
	}

	if err := renderPage(w, artists, "./templates/index.html"); err != nil {
		http.Error(w, fmt.Sprintf("Error rendering HTML: %s", err), http.StatusInternalServerError)
	}
}
