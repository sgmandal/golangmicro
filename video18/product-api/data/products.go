package data

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/building-microservices-youtube/product-api/currencies/protos/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"required,gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"sku"`
}

// Products defines a slice of Product
type Products []*Product

// GetProducts returns all products from the database
func (p *ProductsDB) GetProducts(currency1 string) (Products, error) {
	if currency1 == "" {
		return productList, nil
	}

	resp, err := p.getRate(currency1)
	if err != nil {
		p.log.Error("unable to get rate", currency1, err)
		return nil, err
	}
	pr := Products{}
	for _, p := range productList {
		np := *p
		np.Price = np.Price * resp
		pr = append(pr, &np)
	}

	return pr, nil
}

type ProductsDB struct {
	currency currency.CurrencyClient
	log      hclog.Logger
	rates    map[string]float64
	client1  currency.Currency_SubscribeRatesClient
}

func NewProductDB(c currency.CurrencyClient, l hclog.Logger) *ProductsDB {
	pb := &ProductsDB{c, l, make(map[string]float64), nil}
	go pb.updateHandler()
	return pb
}

func (p *ProductsDB) updateHandler() {
	sb, err := p.currency.SubscribeRates(context.Background())
	if err != nil {
		p.log.Error("error receiving message", "error", err)
	}

	p.client1 = sb
	for {
		rr, err := sb.Recv()
		if grpcError := rr.GetError(); rr.GetError() != nil {
			p.log.Info("error subscribing for rates", grpcError)
			continue
		}

		if resp := rr.GetRateResponse(); resp != nil {
			p.log.Info("received updates from the server", resp.GetDestination().String())
			if err != nil {
				p.log.Error("error receiving message", "error", err)
				return
			}
			p.rates[resp.GetDestination().String()] = resp.Rate
		}

	}
}

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func (p *ProductsDB) GetProductByID(id int, currency1 string) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	if currency1 == "" {
		return productList[i], nil
	}

	resp, err := p.getRate(currency1)
	if err != nil {
		p.log.Error("unable to get rate", currency1, err)
		return nil, err
	}

	np := *productList[i]
	np.Price = resp
	return &np, nil
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func (x *ProductsDB) UpdateProduct(p Product) error {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	// update the product in the DB
	productList[i] = &p

	return nil
}

// AddProduct adds a new product to the database
func (x *ProductsDB) AddProduct(p Product) {
	// get the next id in sequence
	maxID := productList[len(productList)-1].ID
	p.ID = maxID + 1
	productList = append(productList, &p)
}

// DeleteProduct deletes a product from the database
func (x *ProductsDB) DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

func (p *ProductsDB) getRate(destination string) (float64, error) {
	if r, ok := p.rates[destination]; ok {
		return r, nil
	}
	rr := currency.RateRequest{
		Base:        currency.Currencies(currency.Currencies_value["EUR"]),
		Destination: currency.Currencies(currency.Currencies_value[destination]),
	} // type conversion int32

	// get initial rate
	resp, err := p.currency.GetRate(context.Background(), &rr) // what does context do, still doubt
	if err != nil {
		if s, ok := status.FromError(err); ok {
			md := s.Details()[0].(*currency.RateRequest) // asserting the status code error in RateRequest
			if s.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf("unable to get rate, dest and base currencies can not be the same %s %s", md.Base.String(), md.Destination.String())
			}
			return -1, fmt.Errorf("unable to get rate from currency server %s %s", md.Base.String(), md.Destination.String())
		}
		return -1, err
	}

	p.rates[destination] = resp.Rate

	// subscribe for updates
	// err = p.currency.SubscribeRates(context.Background(), rr)

	p.client1.Send(&rr)

	return resp.Rate, err
}

var productList = []*Product{
	{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
	},
	{
		ID:          2,
		Name:        "Esspresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
	},
}
