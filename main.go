package main

import (
	"production-warehouse-api/config"
	"production-warehouse-api/routes"

	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	config.ConnectDB()

	r := gin.Default()
	routes.SetupRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
