package main

import (
	"net"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/building-microservices-youtube/product-api/currencies/data"
	"github.com/nicholasjackson/building-microservices-youtube/product-api/currencies/protos/currency"
	"github.com/nicholasjackson/building-microservices-youtube/product-api/currencies/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()

	rates, err := data.NewRates(log)
	if err != nil {
		log.Error("error")
	}

	gs := grpc.NewServer()
	cs := server.NewCurrency1(rates, log)

	currency.RegisterCurrencyServer(gs, cs)

	reflection.Register(gs)

	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		os.Exit(1)
	}

	gs.Serve(l)
}
