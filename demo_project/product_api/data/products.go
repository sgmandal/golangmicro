package data

import (
	"context"
	"demo/product_api/currencies/protos/currency"

	"github.com/hashicorp/go-hclog"
)

type product struct {

	// product information with their json attributes or key names
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
}

type ProductsDB struct {
	current currency.CurrencyClient
	log     hclog.Logger
	rt      map[string]float64
	client1 currency.Currency_SubscribeRatesClient
}

// a setter function which provides access to this program
// or the struct and its methods present in this file or package
func NewProductDB(c currency.CurrencyClient, l hclog.Logger) *ProductsDB {
	pb := &ProductsDB{current: c, log: l, rt: make(map[string]float64), client1: nil}
	// updateHandler
	return
}

func (p *ProductsDB) apdateHandler() {
	xs, err := p.current.SubscribeRates(context.Background()) // what does this return? an interface
	if err != nil {
		p.log.Error("error receiving message", "error", err)
	}

	p.client1 = xs
	for {
		rr, err := xs.Recv()
		if grperr := rr.GetError(); grperr != nil {
			p.log.Info("error subscribing for rates", grperr)
			continue
		}
	}
}

// helper function,
// fines the index of a product  in the database
// helper function always should start with small alphabet letter
func findIndexByProductID(x int) int {

	// i gets the index, p gets the instance, whole instance
	for i, p := range productList {
		if p.ID == x {
			return i
		}
	}

	return -1
}

// created a product list which is a slice of products
// remember, for struct data, always use pointers
// which is go idiomatic
var productList = []*product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
	},
	{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
	},
}
