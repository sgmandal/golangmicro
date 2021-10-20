package handlers

import (
	"context"
	"hello/presentation/data"
	"log"
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

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "unable to convert id", http.StatusBadRequest)
		return
	}
	p.l.Println("method to handle PUT product", id)
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	err = data.NiggaProduct(id, prod)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "product not found", http.StatusNotFound)
		return
	}

	p.l.Println("PUT successful")
}

func (p *Products) GetRaw(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle get peroducts")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to marshall", http.StatusInternalServerError)
	}
}

func (p *Products) AddRaw(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle POST products")
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	p.l.Printf("prod: %#v", prod)

	data.Addproduct(prod)
}

type KeyProduct struct {
}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		prod.FromJSON(r.Body)

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
