package handlers

import (
	"demo/product_api/data"
	"net/http"
)

func (p *Product1) ListAll(rw http.ResponseWriter, r *http.Request) {
	c := r.URL.Query().Get("currency") // fetching the currency from the url which was passed
	p.l.Debug("get all records")
	rw.Header().Add("Content-type", "application/json")

	prods, err := p.productdb.GetProduct(c) // as the method requires currency1 string as an input, we've passed into it
	// usual error checking
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
	}

	err = data.ToJSON(prods, rw) // feeding the productList struct we got into Json and then response writer

	// error checking, nothing fancy
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("unable to serializing product", err)
	}
}

// list single uses getProductById
func (p *Product1) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r) // a helper function to get product id
	c := r.URL.Query().Get("currency")

	p.l.Debug("get record id", id) //just little logging, indicating that we're at the right part of the code

	prod, err := p.productdb.GetProductById(id, c)

	switch err {
	case nil:

	case data.ErrorProductNotFound:
		p.l.Error("error product not found", err)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return

	default:
		p.l.Error("fetching product", err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = data.ToJSON(prod, rw) // printing the data on screen
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("serializing product", err)
	}
}
