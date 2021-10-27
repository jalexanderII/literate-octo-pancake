package main

import (
	"log"
	"net"

	"github.com/hashicorp/go-hclog"
	"github.com/jalexanderII/literate-octo-pancake/currency/protos/currency"
	"github.com/jalexanderII/literate-octo-pancake/currency/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const port = ":9092"

func main() {
	hlog := hclog.Default()

	// create a TCP socket for inbound server connections
	lstnr, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to start server:", err)
	}

	// setup and register currency service
	// create a new gRPC server, use WithInsecure to allow http connections
	grpcServer := grpc.NewServer()
	// create an instance of the Currency server
	curService := server.NewCurrency(hlog)
	currency.RegisterCurrencyServer(grpcServer, curService)

	// register the reflection service which allows clients to determine the methods
	// for this gRPC service
	reflection.Register(grpcServer)

	// start service's server
	log.Println("starting currency rpc service on", port)
	if err := grpcServer.Serve(lstnr); err != nil {
		log.Fatal(err)
	}
}
