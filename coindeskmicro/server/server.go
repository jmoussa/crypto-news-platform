package server

import (
	"log"
	"net"

	pb "github.com/jmoussa/crypto-dashboard/coindeskmicro/pb"
	"google.golang.org/grpc"
)

func StartServer() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Fail to listen on port 9000: %v", err)
	}

	s := Server{}
	grpcServer := grpc.NewServer()
	log.Println("Registering gRPC server...")
	pb.RegisterCoinDeskScraperServer(grpcServer, &s)
	log.Println("Starting server on port 9000...")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Fail to serve gRPC Server over port 9000: %v", err)
	}
	log.Println("Server started on localhost:9000")
}
