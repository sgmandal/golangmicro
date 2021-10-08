package handlers

import (
	"log"
	"micro/video3/data"
	"net/http"
)

type Products struct {
	l *log.Logger
}

func NewProducts(x *log.Logger) *Products {
	return &Products{x}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	//simplified way of writing the code
	if r.Method == http.MethodGet {
		p.getProducts(rw, r) //function call down below
		return
	}
	// lp := data.GetProducts()
	// err := lp.ToJSON(rw)
	// // d, err := json.Marshal(lp)
	// if err != nil {
	// 	http.Error(rw, "unable to marshall", http.StatusInternalServerError)
	// }
	// // rw.Write(d)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to marshall", http.StatusInternalServerError)
	}
}
