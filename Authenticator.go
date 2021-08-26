package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	Handler "github.com/raihankhan/EcommerceApi/handler"
	"net/http"
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
