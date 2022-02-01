package main

import (
	"log"

	"github.com/jmoussa/crypto-dashboard/coindeskmicro/client"
)

func test_grpc_client_request() {
	content, err := client.FetchCoinDeskData()
	if err != nil {
		log.Fatalf("Error when fetching data from CoinDesk Microservice: %s", err)
	}
	log.Printf("Content from CoinDesk Microservice: %v", content)
}

func main() {
	test_grpc_client_request()
}
