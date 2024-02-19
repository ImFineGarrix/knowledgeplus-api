package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = os.Getenv("SECRET_KEY")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		// Remove "Bearer " prefix from tokenString
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Check the user's role
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Role not specified in token"})
			c.Abort()
			return
		}

		// Store the user's role in the context for further use if needed
		c.Set("userRole", role)

		// Check for specific routes that require "owner" role
		if isOwnerRoute(c) && role != "owner" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func isOwnerRoute(c *gin.Context) bool {
	// Define the routes that require "owner" role here
	ownerRoutes := []string{
		"/api/backoffice/admins",
		"/api/backoffice/admins/:id",
	}

	// Check if the current route is in the list
	currentRoute := c.FullPath()
	fmt.Println(currentRoute)
	for _, route := range ownerRoutes {
		if route == currentRoute {
			return true
		}
	}

	return false
}
