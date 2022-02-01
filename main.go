package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	cdm_client "github.com/jmoussa/crypto-dashboard/coindeskmicro/client" // TODO: move with startAPI() to separate package
	cdm_server "github.com/jmoussa/crypto-dashboard/coindeskmicro/server"
)

// TODO: move and reference crypto-dashboard/api/api.go
func startAPI() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/coindesk", func(c *gin.Context) {
		// microservice client handles request via gRPC
		content, err := cdm_client.FetchCoinDeskData()
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, content)
	})
	router.Run(":3000") // listen and serve on 0.0.0.0:3000 (for windows "localhost:3000")
}

func main() {
	var server bool
	var api bool
	flag.BoolVar(&server, "server", false, "Run the server")
	flag.BoolVar(&api, "api", false, "Run the api")
	flag.Parse()
	log.Printf("Server flag: %v", server)
	log.Printf("api flag: %v", api)
	if server {
		log.Println("Starting Coindesk Microservice server...")
		cdm_server.StartServer()
	} else if api {
		log.Println("Starting API...")
		startAPI()
	}
}
