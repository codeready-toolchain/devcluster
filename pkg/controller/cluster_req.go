package controller

import (
	"net/http"
	"strconv"

	"github.com/alexeykazakov/devcluster/pkg/context"

	"github.com/alexeykazakov/devcluster/pkg/cluster"
	"github.com/alexeykazakov/devcluster/pkg/configuration"
	"github.com/alexeykazakov/devcluster/pkg/errors"
	"github.com/alexeykazakov/devcluster/pkg/log"

	"github.com/gin-gonic/gin"
)

// ClusterRequest implements the cluster request endpoint
type ClusterRequest struct {
	config *configuration.Config
}

// NewClusterRequest returns a new ClusterRequest instance.
func NewClusterRequest(config *configuration.Config) *ClusterRequest {
	return &ClusterRequest{
		config: config,
	}
}

// PostHandler creates a ClusterRequest resource
func (r *ClusterRequest) PostHandler(ctx *gin.Context) {
	ns := ctx.PostForm("number-of-clusters")
	n, err := strconv.Atoi(ns)
	if err != nil {
		log.Error(ctx, err, "error requesting clusters; number of clusters param is missing or invalid")
		errors.AbortWithError(ctx, http.StatusBadRequest, err, "error requesting clusters; number of clusters param is missing or invalid")
		return
	}

	log.Infof(ctx, "Requested provisioning %s clusters", ns)
	req, err := cluster.DefaultClusterService.StartNewRequest(ctx.GetString(context.UsernameKey), n)
	if err != nil {
		log.Error(ctx, err, "error requesting clusters")
		errors.AbortWithError(ctx, http.StatusInternalServerError, err, "error requesting clusters")
	}
	ctx.JSON(http.StatusAccepted, req)
}

// GetHandler returns ClusterRequest resources
func (r *ClusterRequest) GetHandler(ctx *gin.Context) {
	reqs, err := cluster.DefaultClusterService.Requests()
	if err != nil {
		log.Error(ctx, err, "error fetching cluster requests")
		errors.AbortWithError(ctx, http.StatusInternalServerError, err, "error fetching cluster requests")
	}
	ctx.JSON(http.StatusOK, reqs)
}

// GetHandlerClusterReq returns ClusterRequest resource
func (r *ClusterRequest) GetHandlerClusterReq(ctx *gin.Context) {
	req, err := cluster.DefaultClusterService.RequestWithClusters(ctx.Param("id"))
	if err != nil {
		log.Error(ctx, err, "error fetching cluster request")
		errors.AbortWithError(ctx, http.StatusInternalServerError, err, "error fetching cluster request")
	}
	ctx.JSON(http.StatusOK, req)
}
