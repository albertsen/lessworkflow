package grpcserver

import (
	"flag"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port       = flag.Int("p", 50051, "Port to listen on.")
	help       = flag.Bool("h", false, "This message.")
	listenAddr string
)

func init() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(1)
	}
	listenAddr = ":" + strconv.Itoa(*port)
}

func StartServer(regFunc func(server *grpc.Server)) {
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()
	regFunc(server)
	reflection.Register(server)
	log.Printf("Listening on address: %s", listenAddr)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %s", err)
	}
}
