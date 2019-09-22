package main

import (
	"log"
	"lucky/initialize"
	"lucky/routes"
	"net/http"
)

func init() {
	initialize.Init()
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", routes.Router()))
}
