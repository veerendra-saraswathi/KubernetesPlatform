package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler - simple health endpoint
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// GetClustersHandler - placeholder for listing clusters
func GetClustersHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"clusters": []string{"cluster1", "cluster2"},
	})
}

// CreateClusterHandler - placeholder for creating cluster
func CreateClusterHandler(c *gin.Context) {
	type ClusterRequest struct {
		Name string `json:"name"`
	}

	var req ClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cluster created",
		"name":    req.Name,
	})
}

