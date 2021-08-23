package main

import (
	Handler "EcommerceApi/handler"
	"EcommerceApi/products"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized) // 401
				return
			}
		}

		token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(tk *jwt.Token) (interface{}, error) {
			return Handler.Jwtkey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		//Token is authenticated, pass it through
		fmt.Println("authenticator working")
		next.ServeHTTP(w, r)
	})
}

func main() {

	port := ":8080"
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	products.AppendProducts()

	r.Post("/login", Handler.Login)

	r.Group(func(r chi.Router) {

		r.Use(Authenticator)
		r.Get("/view_products", Handler.ViewAll)
		r.Post("/add_product", Handler.AddProduct)

		//r.Route("/{ID}", func(r chi.Router) {
		//	r.Use(ArticleCtx)
		//	r.Get("/", getArticle)       // GET /articles/123
		//	r.Put("/", updateArticle)    // PUT /articles/123
		//	r.Delete("/", deleteArticle) // DELETE /articles/123
		//})

		r.Put("/update_product/{id}", Handler.UpdateProduct)
		r.Delete("/delete_product/{id}", Handler.DelProduct)
	})

	log.Fatal(http.ListenAndServe(port, r))
}

// /home/raihan/go/src/github.com/raihankhan/EcommerceApi
// /home/raihan/go/src/github.com/masudur-rahman/EcommerceApi
// github.com/raihankhan/EcommerceApi
