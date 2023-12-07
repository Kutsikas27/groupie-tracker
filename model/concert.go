package model

type Concert struct {
	Date     string
	Location struct {
		City    string
		Country string
	}
}
