package handlers

import (
	"demo/product_api/data"
	"net/http"
)

func (p *Product1) Create(rw http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Debug("Inserting product")
	p.productdb.AddProduct(prod)
}
