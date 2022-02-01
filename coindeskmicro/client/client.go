package client

import (
	"log"

	pb "github.com/jmoussa/crypto-dashboard/coindeskmicro/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

// Function that is called by the backend API to get the data from the CoinDesk Microservice API via gRPC.
func FetchCoinDeskData() ([]*pb.Content, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dns:///be-srv-lb.default.svc.cluster.local", grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
		return nil, err
	}
	defer conn.Close()

	c := pb.NewCoinDeskScraperClient(conn)
	res, err := c.GetCoinDeskData(context.Background(), &pb.GetCoinDeskDataRequest{MaxEntries: 10})
	if err != nil {
		log.Fatalf("Error when calling GetCoinDeskData: %s", err)
		return nil, err
	}
	log.Printf("Response from server: %v", res.Content)
	return res.Content, nil
}
