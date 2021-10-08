package data

import (
	"encoding/json"
	"io"
	"time"
)

//product defines the structure for an api product
type Product struct {
	ID          int     `json:"id"` //remember no space between colon(:) and inverted comma allowed
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"` //no data displayed in this case
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`

	// ID          int
	// Name        string
	// Description string
	// Price       float64
	// SKU         string
	// CreatedOn   string
	// UpdatedOn   string
	// DeletedOn   string
}

type Products []*Product //creating own type

//we use encode function to dont use memory to store and display data
func (x *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w) //creating an encoder object passed with the io writer, in this case might be rw
	return e.Encode(x)
}

// func GetProducts() []*Product {
// 	return productList
// }

//using our type for linking
func GetProducts() Products {
	return productList
}

//we made a few changes
//code looks different from instructor's code
var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "frothy milk shake",
		Price:       2.45,
		SKU:         "abc234",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fgh234",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
