package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/codeready-toolchain/devcluster/pkg/context"

	"github.com/codeready-toolchain/devcluster/pkg/cluster"
	"github.com/codeready-toolchain/devcluster/pkg/configuration"
	devclustererrors "github.com/codeready-toolchain/devcluster/pkg/errors"
	"github.com/codeready-toolchain/devcluster/pkg/log"

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
		log.Error(ctx, err, "error requesting clusters; number-of-clusters param is missing or invalid")
		devclustererrors.AbortWithError(ctx, http.StatusBadRequest, err, "error requesting clusters; number of clusters param is missing or invalid")
		return
	}

	zone := ctx.PostForm("zone")
	if zone == "" {
		log.Info(ctx, "WARNING: no zone parameter specified. \"wdc04\" will be used by default to create a new request")
		zone = "wdc04"
	}

	deleteIns := ctx.PostForm("delete-in-hours")
	deleteInHours, err := strconv.Atoi(deleteIns)
	if err != nil {
		log.Error(ctx, err, "error requesting clusters; delete-in-hours param is missing or invalid")
		devclustererrors.AbortWithError(ctx, http.StatusBadRequest, err, "error requesting clusters; delete-in-hours param is missing or invalid")
		return
	}

	log.Infof(ctx, "Requested provisioning %s clusters", ns)
	requestedBy := ctx.GetString(context.UsernameKey)

	noSubnet := ctx.PostForm("no-subnet") != ""

	req, err := cluster.DefaultClusterService.CreateNewRequest(requestedBy, n, zone, deleteInHours, noSubnet)
	if err != nil {
		log.Error(ctx, err, "error requesting clusters")
		devclustererrors.AbortWithError(ctx, http.StatusInternalServerError, err, "error requesting clusters")
		return
	}
	ctx.JSON(http.StatusAccepted, req)
}

// GetHandler returns ClusterRequest resources
func (r *ClusterRequest) GetHandler(ctx *gin.Context) {
	reqs, err := cluster.DefaultClusterService.Requests()
	if err != nil {
		log.Error(ctx, err, "error fetching cluster requests")
		devclustererrors.AbortWithError(ctx, http.StatusInternalServerError, err, "error fetching cluster requests")
		return
	}
	ctx.JSON(http.StatusOK, reqs)
}

// GetHandlerClusterReq returns ClusterRequest resource
func (r *ClusterRequest) GetHandlerClusterReq(ctx *gin.Context) {
	reqID := ctx.Param("id")
	req, err := cluster.DefaultClusterService.GetRequestWithClusters(reqID)
	if err != nil {
		log.Error(ctx, err, "error fetching cluster request")
		devclustererrors.AbortWithError(ctx, http.StatusInternalServerError, err, "error fetching cluster request")
		return
	}
	if req == nil { // Not Found
		err = errors.New(fmt.Sprintf("request with id=%s not found", reqID))
		log.Error(ctx, err, "request not found")
		devclustererrors.AbortWithError(ctx, http.StatusNotFound, err, "request not found")
		return
	}
	ctx.JSON(http.StatusOK, req)
}

// GetHandlerZones returns Zones resource
func (r *ClusterRequest) GetHandlerZones(ctx *gin.Context) {
	zones, err := cluster.DefaultClusterService.GetZones()
	if err != nil {
		log.Error(ctx, err, "error fetching zones")
		devclustererrors.AbortWithError(ctx, http.StatusInternalServerError, err, "error fetching zones")
		return
	}
	ctx.JSON(http.StatusOK, zones)
}

// DeleteHandlerCluster deletes Cluster resource
func (r *ClusterRequest) DeleteHandlerCluster(ctx *gin.Context) {
	id := ctx.Param("id")
	err := cluster.DefaultClusterService.DeleteCluster(id)
	if err != nil {
		log.Error(ctx, err, "error deleting cluster")
		devclustererrors.AbortWithError(ctx, http.StatusInternalServerError, err, "error deleting cluster")
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

// PostUsersHandler creates N number of users and returns them as an array (JSON)
func (r *ClusterRequest) PostUsersHandler(ctx *gin.Context) {
	ns := ctx.PostForm("number-of-users")
	n, err := strconv.Atoi(ns)
	if err != nil {
		log.Error(ctx, err, "error requesting users; number-of-users param is missing or invalid")
		devclustererrors.AbortWithError(ctx, http.StatusBadRequest, err, "error requesting users; number-of-users param is missing or invalid")
		return
	}

	sIndex := ctx.PostForm("start-index")
	var startIndex int
	if sIndex != "" {
		startIndex, err = strconv.Atoi(sIndex)
		if err != nil {
			log.Error(ctx, err, "error requesting users; start-index param is not an integer")
			devclustererrors.AbortWithError(ctx, http.StatusBadRequest, err, "error requesting users; start-index param is not an integer")
			return
		}
	}

	log.Infof(ctx, "Requested creating %s users", ns)
	users, err := cluster.DefaultClusterService.CreateUsers(n, startIndex)
	if err != nil {
		log.Error(ctx, err, "error requesting users")
		devclustererrors.AbortWithError(ctx, http.StatusInternalServerError, err, "error requesting users")
		return
	}
	ctx.JSON(http.StatusAccepted, users)
}

// GetUsersHandler returns all users as an array (JSON)
func (r *ClusterRequest) GetUsersHandler(ctx *gin.Context) {
	log.Infof(ctx, "Obtaining users")
	users, err := cluster.DefaultClusterService.Users()
	if err != nil {
		log.Error(ctx, err, "error obtaining users")
		devclustererrors.AbortWithError(ctx, http.StatusInternalServerError, err, "error obtaining users")
		return
	}
	ctx.JSON(http.StatusAccepted, users)
}
