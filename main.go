package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/jmoussa/crypto-dashboard/docs"

	// gRPC Microservices
	cdm_client "github.com/jmoussa/crypto-dashboard/coindeskmicro/client" // TODO: move with startAPI() to separate package
	cdm_server "github.com/jmoussa/crypto-dashboard/coindeskmicro/server"
	twitter_client "github.com/jmoussa/crypto-dashboard/twittermicro/client"
	twitter_server "github.com/jmoussa/crypto-dashboard/twittermicro/server"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// TODO: move and reference crypto-dashboard/api/api.go
func startAPI() {
	router := gin.Default()
	// Routes
	router.GET("/", HealthCheck)
	url := ginSwagger.URL("http://localhost:3000/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
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
	router.POST("/twitter", func(c *gin.Context) {
		// microservice client handles request via gRPC
		max_entries := c.PostForm("max_entries")
		max_entries_int, err := strconv.ParseInt(max_entries, 0, 64)
		if err != nil {
			log.Printf("Error when parsing max_entries: %s\nUsing default value of 100", err)
			max_entries_int = 100
		}
		content, err := twitter_client.FetchTwitterData(max_entries_int)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Could not fetch data from Twitter API",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(200, content)
	})
	router.Run(":3000") // listen and serve on 0.0.0.0:3000 (for windows "localhost:3000")
}

// @title Gin Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	var server bool
	var api bool
	var coindesk bool
	var twitter bool

	flag.BoolVar(&server, "server", false, "Run the server")
	flag.BoolVar(&coindesk, "coindesk", false, "Run the server")
	flag.BoolVar(&twitter, "twitterscraper", false, "Run the server")
	flag.BoolVar(&api, "api", false, "Run the api")
	flag.Parse()
	log.Printf("Server flag: %v", server)
	log.Printf("api flag: %v", api)
	if server {
		if twitter {
			log.Println("Starting Twitter Microservice server...")
			twitter_server.StartServer()
		} else if coindesk {
			log.Println("Starting Coindesk Microservice server...")
			cdm_server.StartServer()
		} else {
			log.Println("Please specify a microservice to run ('coindesk', 'twitter') as a flag")
		}
	} else if api {
		log.Println("Starting Top Level API...")
		startAPI()
	}
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c *gin.Context) {
	res := map[string]interface{}{
		"data": "Server is up and running",
	}

	c.JSON(http.StatusOK, res)
}
