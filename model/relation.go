package model

type DatesLocations = map[string][]string

type Relation struct {
	Index []struct {
		ID             int            `json:"id"`
		DatesLocations DatesLocations `json:"datesLocations"`
	} `json:"index"`
}
