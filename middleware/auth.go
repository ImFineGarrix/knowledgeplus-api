package middleware

import (
	"fmt"
	"knowledgeplus/go-api/response"
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
			out := response.ErrorMsg{
				Code:    http.StatusUnauthorized,
				Field:   "Header",
				Message: "Missing Authorization header",
			}
			c.JSON(http.StatusUnauthorized, out)
			c.Abort()
			return
		}

		// Remove "Bearer " prefix from tokenString
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			out := response.ErrorMsg{
				Code:    http.StatusUnauthorized,
				Field:   "Token",
				Message: "Invalid token",
			}
			c.JSON(http.StatusUnauthorized, out)
			c.Abort()
			return
		}

		// Check the user's role
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			out := response.ErrorMsg{
				Code:    http.StatusUnauthorized,
				Field:   "Token",
				Message: "Invalid token claims",
			}
			c.JSON(http.StatusUnauthorized, out)
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			out := response.ErrorMsg{
				Code:    http.StatusUnauthorized,
				Field:   "Role",
				Message: "Role not specified in token",
			}
			c.JSON(http.StatusForbidden, out)
			c.Abort()
			return
		}

		// Store the user's role in the context for further use if needed
		c.Set("userRole", role)

		// Check for specific routes that require "owner" role
		if isOwnerRoute(c) && role != "owner" {
			out := response.ErrorMsg{
				Code:    http.StatusForbidden,
				Field:   "Role",
				Message: "Permission denied",
			}
			c.JSON(http.StatusForbidden, out)
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

// func errAuthResponse(c *gin.Context, code int, message string) {
// 	err := errors.New(message)
// 	var ve validator.ValidationErrors
// 			out = response.ErrorMsg{
// 				Code:    http.StatusBadRequest,
// 				Field:   fe.Field(),
// 				Message: response.GetErrorMsg(fe),
// 			}
// 		c.JSON(http.StatusCreated, out)
// 	} else {
// 		c.JSON(code, gin.H{"error": message})
// 	}
// }
