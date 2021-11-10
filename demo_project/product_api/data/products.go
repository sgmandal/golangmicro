package data

import (
	"context"
	"demo/product_api/currencies/protos/currency"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorProductNotFound = fmt.Errorf("product not found")

type Product struct {

	// product information with their json attributes or key names
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"retuired,gt=0"`
	SKU         string  `json:"sku" validate:"sku"` // custom validation
}

type multiproducts []*Product // multiproducts is of type a slice of product struct, used for subscriberatesserver

func (pb *ProductsDB) GetProduct(currency1 string) (multiproducts, error) {
	if currency1 == "" {
		return productList, nil
	}

	r, err := pb.getRate(currency1)
	if err != nil {
		pb.log.Error("unable to get rate", currency1, err)
		return nil, err
	}

	pr := multiproducts{}
	for _, p := range productList {
		np := *p                // create a copy
		np.Price = np.Price * r // make pricing change to that copy
		pr = append(pr, &np)
	}
	return pr, nil
}

// what does this function do?
// - fetches the rate from the grpc server
func (pb *ProductsDB) getRate(dest string) (float64, error) {
	if rs, ok := pb.rt[dest]; ok {
		return rs, nil
	}

	// connection to the currencies stub for rpc
	// doubt: yes its done this way, stub is required, might get clear as we proceed
	rr := currency.RateRequest{
		Base:        currency.Currency1(currency.Currency1_value["EUR"]),
		Destination: currency.Currency1(currency.Currency1_value[dest]),
	}

	res, err := pb.current.GetRate(context.Background(), &rr) // calling the stub interface, which inturn calling the getrate method we've written in the grpc server side in its server package, hence we get a rate in grpc call way
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

	pb.rt[dest] = res.Rate // as the value is fetched in struct fashion or rateresponse fashion, we then store this rate

	pb.client1.Send(&rr) // subscribing for the same rates

	return res.Rate, err
}

// idiomatic way to provide connections
type ProductsDB struct {
	current currency.CurrencyClient
	log     hclog.Logger
	rt      map[string]float64 // rates
	client1 currency.Currency_SubscribeRatesClient
}

// a setter function which provides access to this program
// or the struct and its methods present in this file or package
func NewProductDB(c currency.CurrencyClient, l hclog.Logger) *ProductsDB {
	pb := &ProductsDB{current: c, log: l, rt: make(map[string]float64), client1: nil}
	go pb.apdateHandler()
	return pb
}

// getting subscribed rates and adding it to map of rates
// I think the caching is done here
func (p *ProductsDB) apdateHandler() {
	xs, err := p.current.SubscribeRates(context.Background()) // what does this return? an interface, Currency_SubscribeRatesClient
	if err != nil {
		p.log.Error("error receiving message", "error", err)
	}

	p.client1 = xs // same type
	for {
		rr, err := xs.Recv()
		if grperr := rr.GetError(); grperr != nil {
			p.log.Info("error subscribing for rates", grperr)
			continue
		}

		if resp := rr.GetRateResponse(); resp != nil {
			p.log.Info("received updates from the server", resp.GetDestination().String())
			if err != nil {
				p.log.Error("error receiving message", "error", err)
				return
			}
			p.rt[resp.GetDestination().String()] = resp.Rate // cacheing
		}
	}
}

func (pb *ProductsDB) GetProductById(id int, currency1 string) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrorProductNotFound
	}

	if currency1 == "" {
		return productList[i], nil
	}

	rex, err := pb.getRate(currency1)
	if err != nil {
		pb.log.Error("unable to get rate", currency1, err)
		return nil, err
	}

	// returning product by id with updated price calculation
	np := *productList[i]
	np.Price = np.Price * rex // multiplying with rate
	return &np, err
}

func (pb *ProductsDB) UpdateProduct(x Product) error {
	i := findIndexByProductID(x.ID)
	if i == -1 {
		return ErrorProductNotFound
	}

	productList[i] = &x

	return nil
}

func (pb *ProductsDB) AddProduct(x Product) {
	// getting maxid
	maxID := productList[len(productList)-1].ID
	x.ID = maxID + 1
	productList = append(productList, &x)
}

// delete product takes id as in input and deletes the id's entry
func (pb *ProductsDB) DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrorProductNotFound
	}

	productList = append(productList[:i], productList[i+1]) // deleting the product using slicing, assuming id is set in ascending order

	return nil
}

// helper function,
// finds the index of a product  in the database
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
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
	},
}
