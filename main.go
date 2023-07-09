package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	// MaxIdleConnections ...
	MaxIdleConnections = 10
)

func init() {
	var err error

	envConfig, err = readConfig("config.env", ".", map[string]interface{}{
		"dbengine": "memory",
		"port":     "8080",
	})

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	serverPort = envConfig.GetString("port")

	ctx = context.TODO()

	// Initialize DB:
	urlDAO = factoryURLDao(envConfig)
	userDAO = factoryUserDAO(envConfig)
	statsDAO = factoryStatsDao(envConfig)

	gob.Register(&UserInMemory{})
}

func main() {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// Set the router as the default one provided by Gin
	router = gin.Default()

	router.Static("/assets", "./assets")

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("templates/*")

	// Initialize the routes
	initializeRoutes(envConfig)

	// Start serving the applications
	if err := router.Run(net.JoinHostPort("", serverPort)); err != nil {
		log.Fatal(err)
	}
}