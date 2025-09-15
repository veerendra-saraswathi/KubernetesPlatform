package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// Placeholder Auth middleware
func AuthMiddleware(c *gin.Context) {
    // In real implementation, validate token / auth here
    c.Next()
}

// Placeholder handlers
func GetClustersHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "list of clusters"})
}

func CreateClusterHandler(c *gin.Context) {
    c.JSON(http.StatusCreated, gin.H{"message": "cluster created"})
}

func main() {
    r := gin.Default()

    api := r.Group("/api", AuthMiddleware)
    {
        api.GET("/clusters", GetClustersHandler)
        api.POST("/clusters", CreateClusterHandler)
    }

    // Start server on port 8080
    r.Run(":8080")
}

