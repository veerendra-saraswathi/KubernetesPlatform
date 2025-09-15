package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Simple middleware to check Authorization header with a dummy token.
// Later, we can replace this with Keycloak or real JWT validation.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Example: "Bearer mysecrettoken"
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != "mysecrettoken" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// If token is valid, continue
		c.Next()
	}
}

