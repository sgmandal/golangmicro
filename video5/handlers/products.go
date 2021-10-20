package handlers

import (
	"context"
	"log"
	"micro/video5/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(x *log.Logger) *Products {
	return &Products{x}
}

//while using gorillamux
//I've noticed a slight difference
//that difference is the method names, should be capital first letter

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	//handling get products
	p.l.Println("handle GET products")

	//fetch products from data package for now, or a database
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to marshall", http.StatusInternalServerError)
	}
	p.l.Println("GET successful")
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle POST products")

	prod := r.Context().Value(KeyProduct{}).(*data.Product) //casting it as data.Product

	p.l.Printf("prod: %#v", prod)

	data.Addproduct(prod)
	p.l.Println("POST successful")
}

//method to update products
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("method to handle PUT product", id)

	//creating empty struct

	//calling FromJson method created in data/products.go
	//reading the data from request body with r.Body
	//decoding from json
	//testing for errors for decoding function if any
	//now prod has data unmarshalled if no error is thrown

	prod := r.Context().Value(KeyProduct{}).(*data.Product) //casting it as data.Product
	//because it is of type interface, hence such conversion is required

	err = data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "product not found", http.StatusNotFound)
		return
	}

	p.l.Println("PUT successful")
}

type KeyProduct struct {
}

//validating the product
//deserializing

//to use context is to use key
//we use types as key
func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		//add product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
