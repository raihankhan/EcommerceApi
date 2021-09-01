package main

import (
	"log"
	"net/http"
)

func main() {
	port := ":8080"
	r := GetRouter()
	log.Println("Starting Server")
	log.Fatal(http.ListenAndServe(port, r))
}
