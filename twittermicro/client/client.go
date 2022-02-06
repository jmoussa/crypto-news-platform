package client

import (
	"log"

	pb "github.com/jmoussa/crypto-dashboard/twittermicro/twitter_pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

// Function that is called by the backend API to get the data from the CoinDesk Microservice API via gRPC.
func FetchTwitterData(total int64) ([]*pb.Tweet, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("dns:///be-srv-lb.default.svc.cluster.local", grpc.WithInsecure(), grpc.WithBalancerName(roundrobin.Name))
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
		return nil, err
	}
	defer conn.Close()

	c := pb.NewTwitterScraperClient(conn)
	res, err := c.GetTwitterData(context.Background(), &pb.GetTwitterDataRequest{MaxEntries: total})
	if err != nil {
		log.Fatalf("Error when calling GetCoinDeskData: %s", err)
		return nil, err
	}
	log.Printf("Response from server: %v", res.Tweets)
	return res.Tweets, nil
}
