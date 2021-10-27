package data

import "github.com/hashicorp/go-hclog"

type product struct {

	// product information with their json attributes or key names
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
}

type ProductsDB struct {
	log hclog.Logger
}

// a setter function which provides access to this program
// or the struct and its methods present in this file or package
func NewProductDB(l hclog.Logger) *ProductsDB {
	return &ProductsDB{l}
}

// helper function,
// fines the index of a product  in the database
// helper function always should start with small alphabet letter
func findIndexByProductID(x int) int {

	// i gets the index, p gets the instance, whole instance
	for i, p := range productList {
		if p.ID == x {
			return i
		}
	}

	return -1
}

// created a product list which is a slice of products
// remember, for struct data, always use pointers
// which is go idiomatic
var productList = []*product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
	},
}
