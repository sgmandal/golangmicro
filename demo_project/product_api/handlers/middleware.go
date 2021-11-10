/*product validation middleware, middleware as we know is a feature provided by gorillamux*/

package handlers

import (
	"context"
	"demo/product_api/data"
	"net/http"
)

func (p *Product1) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// declaring an empty container struct
		prod := &data.Product{} // capitalization of first letter was required in data/products.go to use it here

		err := data.FromJSON(prod, r.Body)
		if err != nil {
			p.l.Error("deserializing product", err)

			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		}

		errs := p.v.Validate(prod)
		if len(errs) != 0 {
			p.l.Error("validating product", errs)

			// return the validation messages as an array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod) // keyproduct is empty struct
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
