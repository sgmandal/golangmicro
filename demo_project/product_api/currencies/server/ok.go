package server

import (
	"context"
	"demo/product_api/currencies/data"
	"demo/product_api/currencies/protos/currency"
	"io"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type currency2 struct {
	xs            *data.ExchangeRates
	log           hclog.Logger
	subscriptions map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest // currency package RateRequest struct
}

func NewCurrency1(r *data.ExchangeRates, l hclog.Logger) *currency2 {
	x := &currency2{r, l, make(map[currency.Currency_SubscribeRatesServer][]*currency.RateRequest)}

	go x.handleUpdates()
	return x
}

func (cx *currency2) handleUpdates() {
	ru := cx.xs.MonitorRates(5 * time.Second) // this monitor rate function starts with ticker as 5 seconds and we get updated rates in xs

	for range ru {
		cx.log.Info("Got Updated Rates")

		// loop over subscribed clients
		// main part of the code
		for x, v := range cx.subscriptions { // x stores interface reference

			// loop over subscribed rates
			for _, rr := range v {
				r, err := cx.xs.GetRae(rr.GetBase().String(), rr.Destination.String())
				if err != nil {
					cx.log.Error("unable to get updated rate", rr.GetBase().String(), rr.GetBase().String())
				}

				// rr is fetched from the rpc stub interfaces, and we compute rate according to that and send it via our inside written codes
				err = x.Send(&currency.StreamingRateResponse{
					Message: &currency.StreamingRateResponse_RateResponse{
						RateResponse: &currency.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r},
					},
				})
				if err != nil {
					cx.log.Error("unable to get updated rate", rr.GetBase().String(), rr.GetBase().String())
				}
			}
		}
	}
}

// following code to implement stub interfaces

func (cx *currency2) GetRate(ctx context.Context, rr *currency.RateRequest) (*currency.RateResponse, error) {
	cx.log.Info("Handle GetRate", rr.GetBase(), rr.GetDestination()) // rr.GetBase and destination is fetched from the client

	// main code starts from here for this method
	if rr.Base == rr.Destination {

		// rpc error handling
		// we use status package provided in golang
		err := status.Newf(
			codes.InvalidArgument,
			"base cannot be %s the same as destination %s",
			rr.Base.String(),
			rr.Destination.String(),
		)

		err, wde := err.WithDetails(rr) // can add metdata to error, also throws its own error if not handled properly

		if wde != nil {
			return nil, wde
		}

		return nil, err.Err()

	}

	rate1, err := cx.xs.GetRae(rr.GetBase().String(), rr.GetDestination().String()) // we call this function to get rate
	if err != nil {
		return nil, err
	}
	return &currency.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate1}, nil

}

// not understood this method
// doubt, this method is not implementing the currency.pb.go method
// how is it getting used
// assumption: maybe it doesn't depend upon what is returned or passed
func (cx *currency2) SubscribeRates(ls currency.Currency_SubscribeRatesServer) error {
	for {
		rr, err := ls.Recv()
		// EOF signifies end of file
		if err == io.EOF {
			cx.log.Error("Client has closed connection")
			break
		} else if err != nil {
			cx.log.Error("Unable to read from client", err)
		}

		cx.log.Info("Handle client request", "request_base", rr.GetBase(), "request_dest", rr.GetDestination())

		xrs, ok := cx.subscriptions[ls]
		if !ok {
			xrs = []*currency.RateRequest{}
		}

		var validerr *status.Status
		for _, v := range xrs {

			// ranging over until we get what we need
			if v.Base == rr.Base && v.Destination == rr.Destination {

				// rpc error
				validerr = status.Newf(
					codes.AlreadyExists,
					"Unable to subscribe for currency as subscription already exists",
				)
				validerr, err = validerr.WithDetails(rr)
				if err != nil {
					cx.log.Error("unable to add metadata to error", "error", err)
					break
				}

				if validerr != nil {
					ls.Send(&currency.StreamingRateResponse{Message: &currency.StreamingRateResponse_Error{
						Error: validerr.Proto(),
					}})
					continue
				}
			}
		}

		// executes when all ok
		xrs = append(xrs, rr)
		cx.subscriptions[ls] = xrs
	}
	return nil
}
