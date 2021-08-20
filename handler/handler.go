package Handler

import (
	"EcommerceApi/products"
	"encoding/json"
	"net/http"
)

//var jwtkey = []byte("sage-jutsu")
//
//var User = map[string]string {
//	"User1" : "password1",
//	"User2" : "password2",
//}
//
//type Credentials struct {
//	Username string `json:"username"`
//	Password string `json:"password"`
//}
//
//type Claims struct {
//	Username string `json:"username"`
//	jwt.StandardClaims
//}
//
//func Login(w http.ResponseWriter , r *http.Request) {
//	var credentials Credentials
//
//	err := json.NewDecoder(r.Body).Decode(&credentials) // decode the bye request body to json and assign to credentials
//	if err!=nil {
//		w.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	expectedPassword , available  := User[credentials.Username]
//
//	if !available || expectedPassword!=credentials.Password {	// check is credentials exists and matches
//		w.WriteHeader(http.StatusUnauthorized)
//		return
//	}
//q
//	expirationTime := time.Now().Add(time.Minute*10)
//
//	claims := &Claims{
//		Username: credentials.Username,
//		StandardClaims : jwt.StandardClaims{
//			ExpiresAt: expirationTime.Unix(),
//		},
//	}
//
//	// token := jwt.NewWithClaims(jwt.)
//
//
//
//
//}

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



	//w.Header().Set("Context-Type", "application/json")
	//err := json.NewEncoder(w).Encode(prod)
	//if err != nil {
	//	return
	//}

	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(prod)
	_, err := w.Write(data)
	if err != nil {
		return
	}

}
