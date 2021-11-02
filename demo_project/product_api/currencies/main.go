package main

import (
	"demo/product_api/currencies/data"
	"demo/product_api/currencies/protos/currency"
	"demo/product_api/currencies/server"
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default() // logging purpose

	r, err := data.NewRates(log)
	if err != nil {
		log.Error("error")
	}

	gs := grpc.NewServer() // creating a grpc server

	xs := server.NewCurrency1(r, log)

	currency.RegisterCurrencyServer(gs, xs) // connection with the stub

	reflection.Register(gs) // idiomatic grpc

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		os.Exit(1)
	}

	gs.Serve(l)

}
