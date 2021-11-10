package handlers

import (
	"demo/product_api/data"
	"net/http"
)

func (p *Product1) Update(rw http.ResponseWriter, r *http.Request) {
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Debug("Updating Record", prod.ID) // just a stupid logger

	err := p.productdb.UpdateProduct(prod)
	if err == data.ErrorProductNotFound {
		p.l.Error("product not found", err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: "Product not found in database"}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
