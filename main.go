package main

import (
	"knowledgeplus/go-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter()
	_ = r.Run(":8081")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	defaultPath := r.Group("/api")

	// Call the SetupRoutes function from the routes package
	routes.SetupRoutes(defaultPath)

	return r
}
