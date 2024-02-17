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
	config.AllowOrigins = []string{"http://cp23sj2.sit.kmutt.ac.th/", "https://capstone23.sit.kmutt.ac.th/sj2/"} // Add your frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	r.Use(cors.New(config))
	defaultPath := r.Group(os.Getenv("APIPATH"))

	// Call the SetupRoutes function from the routes package
	routes.SetupRoutes(defaultPath)

	return r
}
