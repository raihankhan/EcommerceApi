package main

import (
	Handler "EcommerceApi/handler"
	"EcommerceApi/products"
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func IDCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var prod *products.Product
		var err error

		if prodID := chi.URLParam(r, "id"); prodID != "" {
			prod, err = products.Find(prodID)
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "id", prod)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	products.AppendProducts()

	r.Post("/login", Handler.Login)

	r.Group(func(r chi.Router) {

		r.Use(Authenticator)

		r.Route("/products", func(r chi.Router) {

			r.Get("/", Handler.View)
			r.Post("/", Handler.AddProduct)

			r.Route("/{id}", func(r chi.Router) {

				r.Use(IDCtx)
				r.Get("/", Handler.GetProduct)
				r.Delete("/", Handler.DelProduct)
				r.Put("/", Handler.UpdateProduct)
			})

		})
	})

	return r
}
