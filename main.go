package main

import (
	"EcommerceApi/handler"
	"EcommerceApi/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {

	port := ":8080"
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	products.AppendProducts()

	r.HandleFunc("/login", Handler.Login)
	r.Get("/viewProducts", Handler.ViewAll)


	log.Fatal(http.ListenAndServe(port, r))
}

// /home/raihan/go/src/github.com/raihankhan/EcommerceApi
// /home/raihan/go/src/github.com/masudur-rahman/EcommerceApi
// github.com/raihankhan/EcommerceApi
