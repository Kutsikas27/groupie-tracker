package model

type Relation struct {
	Index []struct {
		ID      int                 `json:"id"`
		Concert map[string][]string `json:"datesLocations"`
	} `json:"index"`
}
