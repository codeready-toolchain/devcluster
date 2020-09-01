package controller

import (
	"github.com/alexeykazakov/devcluster/pkg/configuration"
	"github.com/gin-gonic/gin"
)

// ClusterRequest implements the cluster request endpoint
type ClusterRequest struct {
	config        *configuration.Config
}

// NewClusterRequest returns a new ClusterRequest instance.
func NewClusterRequest(config *configuration.Config) *ClusterRequest {
	return &ClusterRequest{
		config:        config,
	}
}

// PostHandler creates a ClusterRequest resource
func (r *ClusterRequest) PostHandler(ctx *gin.Context) {
}