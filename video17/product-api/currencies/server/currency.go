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
	rates         *data.ExchangeRates
	log           hclog.Logger
	subscriptions map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest
}

func NewCurrency1(r *data.ExchangeRates, l hclog.Logger) *Currency {
	c := &Currency{r, l, make(map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest)}

	go c.handleUpdates()
	return c
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("Got Updated Rates")

		// loop over subscribed clients
		for k, v := range c.subscriptions {

			// loop over subscribed rates
			for _, rr := range v {
				r, err := c.rates.GetRae(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("unable to get updated rate", rr.GetBase().String(), rr.GetBase().String())
				}

				err = k.Send(&currency.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r})
				if err != nil {
					c.log.Error("unable to get updated rate", rr.GetBase().String(), rr.GetBase().String())
				}
			}
		}
	}
}

func (c *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handle GetRate", rr.GetBase(), rr.GetDestination())

	rate1, err := c.rates.GetRae(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	return &currency.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate1}, nil
}

func (c *Currency) SubscribeRates(xd currency.Currency_SubscribeRatesServer) error {
	for {
		rr, err := xd.Recv()
		if err == io.EOF {
			c.log.Error("Client has closed connection")
			break
		}

		if err != nil {
			c.log.Error("Unable to read from client", err)
		}

		c.log.Info("Handle client request", "request_base", rr.GetBase(), "request_dest", rr.GetDestination())

		rrs, ok := c.subscriptions[xd]
		if !ok {
			rrs = []*currency.RateRequest{}
		}

		rrs = append(rrs, rr)
		c.subscriptions[xd] = rrs
	} // blocking method, waits for user's requests

	return nil
}
