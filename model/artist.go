package model

type Artist struct {
	Name         string   `json:"name"`
	Image        string   `json:"image"`
	FirstAlbum   string   `json:"firstAlbum"`
	Id           int      `json:"id"`
	CreationDate int      `json:"creationDate"`
	Members      []string `json:"members"`
}
