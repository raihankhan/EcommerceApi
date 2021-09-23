package Handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/raihankhan/ecommerceApi/products"
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

	// decode the bye request body to json and assign to credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check is credentials exists and matches
	expectedPassword, available := User[credentials.Username]
	if !available || expectedPassword != credentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Set a session expiration time
	expirationTime := time.Now().Add(time.Minute * 20)

	// Create a claim object
	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// generate the signed token string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Jwtkey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// set the token string and expiration time in the cookie
	http.SetCookie(w,
		&http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		},
	)
	w.WriteHeader(http.StatusOK)
}

func SetResponse(w http.ResponseWriter, prod map[string]products.Product) {
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(prod)
	_, err := w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func View(w http.ResponseWriter, r *http.Request) {

	brand := r.URL.Query().Get("brand")
	category := r.URL.Query().Get("category")

	prod := products.Products

	//filter product by brand
	if len(brand) != 0 {
		tmp := make(map[string]products.Product)
		for key, product := range products.Products {
			if product.Brand == brand {
				tmp[key] = product
			}
		}
		prod = tmp
	}

	//filter product by category
	if len(category) != 0 {
		for key, product := range products.Products {
			if product.Category != category {
				delete(prod, key)
			}
		}
	}

	SetResponse(w, prod)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	prod, ok := ctx.Value("id").(*products.Product)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(prod)
	_, err := w.Write(data)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
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
	} else if newProd.ID == "" {
		w.Write([]byte("Product ID can't be empty"))
		w.WriteHeader(http.StatusConflict)
		return
	}

	products.Products[newProd.ID] = newProd
	w.Write([]byte("Product " + newProd.ID + " Has been Added"))
	w.WriteHeader(http.StatusOK)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	prod, ok := ctx.Value("id").(*products.Product)
	if !ok {
		http.Error(w, http.StatusText(422), 422)
		return
	}

	var updatedProd products.Product
	err := json.NewDecoder(r.Body).Decode(&updatedProd)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	products.Products[prod.ID] = updatedProd
	SetResponse(w, products.Products)
}

func DelProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	prod, _ := ctx.Value("id").(*products.Product)
	delete(products.Products, prod.ID)
	SetResponse(w, products.Products)
}

//func AddFeature(w http.ResponseWriter, r *http.Request) {
//
//}
