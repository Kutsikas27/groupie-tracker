package service

import (
	"fmt"
	"groupie-tracker/model"
	"sort"
	"strconv"
	"strings"
)

func sortConcerts(c1, c2 *model.Concert) bool {
	date1Parts := strings.Split(c1.Date, ".")
	date2Parts := strings.Split(c2.Date, ".")

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
func normlizeDate(date string) string {
	splitDate := strings.Split(date, "-")
	joinDate := strings.Join(splitDate, ".")

	return joinDate
}
func FindConcertDataByID(concertData model.Relation, id int) ([]model.Concert, error) {
	for _, entry := range concertData.Index {
		if entry.ID == id {
			var concerts []model.Concert
			for city, dates := range entry.Concert {
				for _, date := range dates {
					city2, country := splitCountryAndCity(city)
					concert := model.Concert{
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
				return sortConcerts(&concerts[i], &concerts[j])
			})
			return concerts, nil
		}
	}
	return nil, fmt.Errorf("id %d not found", id)
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
