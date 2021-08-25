package main

import (
	"log"
	"net/http"
)

func main() {

	port := ":8080"
	r := GetRouter()
	log.Fatal(http.ListenAndServe(port, r))
}
