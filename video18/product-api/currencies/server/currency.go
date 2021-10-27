package server

import (
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/building-microservices-youtube/product-api/currencies/data"
	"github.com/nicholasjackson/building-microservices-youtube/product-api/currencies/protos/currency"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
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

				err = k.Send(&currency.StreamingRateResponse{
					Message: &currency.StreamingRateResponse_RateResponse{
						RateResponse: &currency.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r},
					},
				})
				if err != nil {
					c.log.Error("unable to get updated rate", rr.GetBase().String(), rr.GetBase().String())
				}
			}
		}
	}
}

func (c *Currency) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	c.log.Info("Handle GetRate", rr.GetBase(), rr.GetDestination())

	if rr.Base == rr.Destination {
		// err := status.Errorf(
		// 	codes.InvalidArgument,
		// 	"base cannot be %s the same as the destination %s",
		// 	rr.Base.String(),
		// 	rr.Destination.String(),
		// )

		// usingstatus package provided by golang
		// calling the unary operator in grpc
		err := status.Newf(
			codes.InvalidArgument,
			"base cannot be %s the same as the destination %s",
			rr.Base.String(),
			rr.Destination.String(),
		)

		// grpc error handling method
		// status.Newf has many such method which can be used
		err, wde := err.WithDetails(rr) // can add metdata to error, also throws its own error if not handled properly
		if wde != nil {
			return nil, wde
		}

		return nil, err.Err() // print error as it is
	}

	rate1, err := c.rates.GetRae(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}
	return &currency.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate1}, nil
}

func (c *Currency) SubscribeRates(xd currency.Currency_SubscribeRatesServer) error {

	// handle client messages
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

		// check that subscription does not exist
		var validationError *status.Status
		for _, v := range rrs {
			if v.Base == rr.Base && v.Destination == rr.Destination {
				// subscription exists return errors

				//convenience package Newf
				validationError = status.Newf(
					codes.AlreadyExists,
					"Unable to subscribe for currency as subscription already exists",
				)

				validationError, err = validationError.WithDetails(rr)
				if err != nil {
					c.log.Error("unable to add metadata to error", "error", err)
					break
				}

				// if a validation error return error and continue
				if validationError != nil {
					xd.Send(&currency.StreamingRateResponse{Message: &currency.StreamingRateResponse_Error{
						Error: validationError.Proto(),
					}})
					continue
				}

			}
		}

		// all ok
		rrs = append(rrs, rr)
		c.subscriptions[xd] = rrs
	} // blocking method, waits for user's requests

	return nil
}
