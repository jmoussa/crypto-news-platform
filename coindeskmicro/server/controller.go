package server

import (
	"fmt"
	"log"

	pb "github.com/jmoussa/crypto-dashboard/coindeskmicro/pb"
	"github.com/jmoussa/crypto-dashboard/coindeskmicro/scraper"
	"golang.org/x/net/context"
)

type Server struct {
	pb.UnimplementedCoinDeskScraperServer
}

func (s *Server) GetCoinDeskData(ctx context.Context, message *pb.GetCoinDeskDataRequest) (*pb.GetCoinDeskDataResponse, error) {
	log.Printf("Received request: %v", message)
	// Scrape data from CoinDesk API here
	data, err := scraper.ScrapeCoinDeskData(message)
	if err != nil {
		return nil, fmt.Errorf("Error when scraping data from CoinDesk API: %s", err)
	}
	// Return Scraped Data to Client
	return &pb.GetCoinDeskDataResponse{Content: data}, nil
}
