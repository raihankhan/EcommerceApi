package products

import "errors"

type Category struct {
	Categories []Product `json:"Categories"`
}

type Product struct {
	ID          string `json:"id"`
	ProductName string `json:"product_name"`
	Brand       string `json:"brand"`
	Category    string `json:"category"`
	IsAvailable bool   `json:"is_available"`
	Features    []Feature
	Price
}

type Feature struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Price struct {
	BasePrice       float64 `json:"base_price"`
	Discount        float64 `json:"discount"`
	DiscountedPrice float64 `json:"discounted_price"`
}

var Products map[string]Product

func calcDiscount(p Product) Product {
	p.DiscountedPrice = p.BasePrice - (p.BasePrice*p.Discount)/100
	return p
}
func calcDiscountedPrice() {
	for id, item := range Products {
		Products[id] = calcDiscount(item)
	}
}

func AppendProducts() {
	Products = make(map[string]Product)
	Products["MP01"] = Product{
		ID:          "MP01",
		ProductName: "Samsung J7",
		Brand:       "Samsung",
		Category:    "Mobile Phone",
		IsAvailable: true,
		Features:    nil,
		Price: Price{
			BasePrice:       22000,
			Discount:        3,
			DiscountedPrice: 0,
		},
	}
	Products["MP02"] = Product{
		ID:          "MP02",
		ProductName: "Samsung J7",
		Brand:       "Samsung",
		Category:    "Mobile Phone",
		IsAvailable: true,
		Features:    nil,
		Price: Price{
			BasePrice:       22000,
			Discount:        3,
			DiscountedPrice: 0,
		},
	}
	Products["MP03"] = Product{
		ID:          "MP03",
		ProductName: "OnePlus Nord",
		Brand:       "OnePLus",
		Category:    "Mobile Phone",
		IsAvailable: true,
		Features:    nil,
		Price: Price{
			BasePrice:       36000,
			Discount:        2,
			DiscountedPrice: 0,
		},
	}
	Products["LT01"] = Product{
		ID:          "LT01",
		ProductName: "HP Pavilion X46",
		Brand:       "HP",
		Category:    "Computer",
		IsAvailable: true,
		Features:    nil,
		Price: Price{
			BasePrice:       22000,
			Discount:        3,
			DiscountedPrice: 0,
		},
	}
	calcDiscountedPrice()
}

func Find(key string) (*Product, error) {
	item, exists := Products[key]
	if !exists {
		return nil, errors.New("product not found")
	}
	return &item, nil
}
