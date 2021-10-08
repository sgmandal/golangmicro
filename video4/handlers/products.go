package handlers

import (
	"log"
	"micro/video4/data"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(x *log.Logger) *Products {
	return &Products{x}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	//simplified way of writing the code
	//handling Get request
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

	//for reading json
	//method post
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	//demonstration of PUT method
	if r.Method == http.MethodPut {
		p.l.Println("PUT")
		//expect the id in uri
		//ID is difficult to fetch using only the standard library go provides
		//hence manual work is required for that, example, mention the client that
		//put the input in so and so format
		//we can use regex instead
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			p.l.Println("invalid url more than one id")
			http.Error(rw, "invalid url", http.StatusBadRequest)
			return

		}

		if len(g[0]) != 2 {
			p.l.Println("invalid url more than one capture group") //need to know its significance, must be something to do with regex
			http.Error(rw, "invalid url", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString) //converts string to integer
		if err != nil {
			p.l.Println("Invalid URI unable to convert to number", idString)
			http.Error(rw, "invalid url", http.StatusBadRequest)
			return
		}

		p.l.Println("got id ", id)

		p.updateProducts(id, rw, r)
		return
	}
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	//handling get products
	p.l.Println("handle get peroducts")

	//fetch products from data package for now, or a database
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to marshall", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle POST products")

	prod := &data.Product{}      //empty initialization, receiving block
	err := prod.FromJSON(r.Body) //calling from json function which decodes, and sets the result into prod variable
	if err != nil {
		http.Error(rw, "unable to unmarshall json", http.StatusBadRequest) //using the correct http status code
	}

	p.l.Printf("prod: %#v", prod)

	data.Addproduct(prod)
}

//method to update products
func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("method to handle PUT product")

	//creating empty struct
	prod := &data.Product{}

	//calling FromJson method created in data/products.go
	//reading the data from request body with r.Body
	//decoding from json
	err := prod.FromJSON(r.Body)
	//testing for errors for decoding function if any
	//now prod has data unmarshalled if no error is thrown
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound {
		http.Error(rw, "product not found", http.StatusNotFound)
		return
	}

	p.l.Println("PUT successful")
}
