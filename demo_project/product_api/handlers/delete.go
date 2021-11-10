package handlers

import (
	"demo/product_api/data"
	"net/http"
)

func (p *Product1) Delete(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r) // getting productid written in products.go in same package

	p.l.Debug("deleting record id", id) // nothing meaningful happening here, just logging

	err := p.productdb.DeleteProduct(id)
	//now we need to check error
	if err == data.ErrorProductNotFound {
		p.l.Error("Deleting record id does not exist") // log error

		rw.WriteHeader(http.StatusNotFound) // responsewriter error printing

		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	} else if err != nil {
		p.l.Error("deleting record", "error", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	rw.WriteHeader(http.StatusNoContent) // HTTP Status 204 (No Content) indicates that the server has successfully fulfilled the request and that there is no content to send in the response payload body.
}
