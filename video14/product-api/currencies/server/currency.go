package server

import (
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/building-microservices-youtube/product-api/currencies/protos/currency"
	"golang.org/x/net/context"
)

type Currency struct {
	log hclog.Logger
}

func NewCurrency1(l hclog.Logger) *Currency {
	return &Currency{l}
}

func (c *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handle GetRate", rr.GetBase(), rr.GetDestination())

	return &currency.RateResponse{Rate: 0.5}, nil
}
