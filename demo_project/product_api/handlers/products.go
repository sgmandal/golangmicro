package handlers

import (
	"demo/product_api/data"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

// KeyProduct is a key used for the product object in the context - no idea what this means
type KeyProduct struct{}

type Product1 struct {
	l         hclog.Logger
	v         *data.Validation // first letter capital is required to import it like this
	productdb *data.ProductsDB
}

// connetor function
func NewProducts(l hclog.Logger, v *data.Validation, cc *data.ProductsDB) *Product1 {
	return &Product1{l, v, cc}
}

// defining an ErrInvalidProductPath error message when path is not valid
var ErrInvalidProductPath = fmt.Errorf("invalid path, example: /products/[id]")

type GenericError struct {
	Message string `json:"message"`
}

// this is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

func getProductID(r *http.Request) int {
	// pasre the product id from the url
	vars := mux.Vars(r) // vars variable now hold the content of url in map fashion

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"]) // string to integer conversion of id from map
	if err != nil {
		panic(err)
	}

	return id
}
