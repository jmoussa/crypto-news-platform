package server

import (
	"fmt"
	"log"

	"github.com/jmoussa/crypto-dashboard/twittermicro/scraper"
	pb "github.com/jmoussa/crypto-dashboard/twittermicro/twitter_pb"
	"golang.org/x/net/context"
)

type Server struct {
	pb.UnimplementedTwitterScraperServer
}

func (s *Server) GetTwitterData(ctx context.Context, message *pb.GetTwitterDataRequest) (*pb.GetTwitterDataResponse, error) {
	log.Printf("Received request: %v", message)
	// Scrape data from Twitter API here
	data, err := scraper.ScrapeTwitterData(message)
	if err != nil {
		return nil, fmt.Errorf("Error when scraping data from Twitter API: %s", err)
	}
	// Return Scraped Data to Client
	return &pb.GetTwitterDataResponse{Tweets: data}, nil
}
