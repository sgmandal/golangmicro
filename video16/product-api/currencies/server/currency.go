package server

import (
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/building-microservices-youtube/product-api/currencies/data"
	"github.com/nicholasjackson/building-microservices-youtube/product-api/currencies/protos/currency"
	"golang.org/x/net/context"
)

type Currency struct {
	rates *data.ExchangeRates
	log   hclog.Logger
}

func NewCurrency1(r *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{r, l}
}

func (c *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handle GetRate", rr.GetBase(), rr.GetDestination())

	rate1, err := c.rates.GetRae(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	return &currency.RateResponse{Rate: rate1}, nil
}

func (c *Currency) SubscribeRates(xd currency.Currency_SubscribeRatesServer) error {

	go func() {
		for {
			rr, err := xd.Recv()
			if err == io.EOF {
				c.log.Error("Client has closed connection")
			}
			if err != nil {
				c.log.Error("Unable to read from client", err)
			}

			c.log.Info("Handle client request", "request_base", rr.GetBase(), "request_dest", rr.GetDestination())
		} // blocking method, waits for user's requests
	}()

	for {
		err := xd.Send(&currency.RateResponse{Rate: 12.1}) // calling the interface methods, Send returns error
		if err != nil {
			return err
		}
		time.Sleep(10 * time.Second) // pretending to be some lines of code
	}
}
