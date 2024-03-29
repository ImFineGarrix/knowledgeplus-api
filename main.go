package main

import (
	"knowledgeplus/go-api/initializers"
	"knowledgeplus/go-api/routes"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter()
	_ = r.Run(":8081")
}

func setupRouter() *gin.Engine {
	initializers.LoadEnvVariables()
	r := gin.Default()
	// r.Use(cors.Default())
	// Customize CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Add your frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	r.Use(cors.New(config))
	defaultPath := r.Group(os.Getenv("APIPATH"))

	// Call the SetupRoutes function from the routes package
	routes.SetupRoutes(defaultPath)

	return r
}
