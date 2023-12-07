package model

type Concert struct {
	IsoDate  string
	Date     string
	Location Location
}
type Location struct {
	City    string
	Country string
}
