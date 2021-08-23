package Handler

import (
	"EcommerceApi/products"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var Jwtkey = []byte("sage-jutsu")

var User = map[string]string{
	"User1": "password1",
	"User2": "password2",
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(w http.ResponseWriter, r *http.Request) {

	var credentials Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials) // decode the bye request body to json and assign to credentials
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, available := User[credentials.Username]

	if !available || expectedPassword != credentials.Password { // check is credentials exists and matches
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)

	claims := &Claims{ // Create a claim object
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(Jwtkey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		},
	)
	w.WriteHeader(http.StatusOK)
}

func ViewAll(w http.ResponseWriter, r *http.Request) {

	brand := r.URL.Query().Get("brand")
	prod := products.Products
	if len(brand) != 0 {
		tmp := make(map[string]products.Product)
		for key, product := range products.Products {
			if product.Brand == brand {
				tmp[key] = product
			}
		}
		prod = tmp
	}

	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(prod)
	_, err := w.Write(data)
	if err != nil {
		return
	}

}

func AddProduct(w http.ResponseWriter, r *http.Request) {
	var newProd products.Product
	err := json.NewDecoder(r.Body).Decode(&newProd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, doesExist := products.Products[newProd.ID]
	if doesExist {
		w.Write([]byte("Product already exists"))
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	products.Products[newProd.ID] = newProd

}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var prod products.Product
	err := json.NewDecoder(r.Body).Decode(&prod)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, doesExist := products.Products[prod.ID]
	if !doesExist {
		w.Write([]byte("Product doesn't exists"))
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	products.Products[prod.ID] = prod

}

func DelProduct(w http.ResponseWriter, r *http.Request) {
	var prod products.Product
	err := json.NewDecoder(r.Body).Decode(&prod)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, doesExist := products.Products[prod.ID]
	if !doesExist {
		w.Write([]byte("Product doesn't exists"))
		w.WriteHeader(http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	delete(products.Products, prod.ID)

}
