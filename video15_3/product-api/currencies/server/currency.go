package server

import (
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
