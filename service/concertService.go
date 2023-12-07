package service

import (
	"fmt"
	"groupie-tracker/model"
	"sort"
	"strings"
)

func getIsoDate(str string) string {
	dateArr := strings.Split(str, "-")
	isoDate := dateArr[2] + "." + dateArr[1] + "." + dateArr[0]
	return isoDate
}

func normlizeDate(date string) string {
	splitDate := strings.Split(date, "-")
	joinDate := strings.Join(splitDate, ".")

	return joinDate
}
func findDatesLocationsById(concertData model.Relation, id int) (model.DatesLocations, error) {
	for _, entry := range concertData.Index {
		if entry.ID == id {
			return entry.DatesLocations, nil
		}
	}
	return model.DatesLocations{}, fmt.Errorf("id %d not found", id)
}

func FindConcertsByID(concertData model.Relation, id int) ([]model.Concert, error) {
	datesLocations, err := findDatesLocationsById(concertData, id)
	if err != nil {
		return nil, fmt.Errorf("id %d not found", id)
	}
	var concerts []model.Concert
	for location, dates := range datesLocations {
		for _, date := range dates {
			city, country := splitCountryAndCity(location)
			concert := model.Concert{
				Date:    normlizeDate(date),
				IsoDate: getIsoDate(date),
				Location: model.Location{
					City:    city,
					Country: country,
				},
			}
			concerts = append(concerts, concert)
		}
	}
	sort.Slice(concerts, func(i, j int) bool {
		return concerts[i].IsoDate > concerts[j].IsoDate
	})
	return concerts, nil
}

func splitCountryAndCity(str string) (string, string) {
	locationArr := strings.Split(str, "-")
	locationStr := strings.Join(locationArr, " ")
	cityAndCountry := strings.Split(locationStr, " ")

	city := strings.ToTitle(cityAndCountry[0])
	country := strings.ToTitle(cityAndCountry[1])
	return normlizeString(city), normlizeString(country)
}
func normlizeString(str string) string {
	splitStr := strings.Split(str, "_")
	return strings.Join(splitStr, " ")
}
