package main

import (
	"fmt"
	"groupie-tracker/controller"
	"log"
	"net/http"
)

func main() {
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.HandleFunc("/artist/", controller.ArtistPageHandler)
	http.HandleFunc("/", controller.MainPageHandler)
	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
