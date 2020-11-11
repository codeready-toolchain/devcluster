package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/codeready-toolchain/devcluster/pkg/cluster"
	"github.com/codeready-toolchain/devcluster/pkg/configuration"
	"github.com/codeready-toolchain/devcluster/pkg/context"
	devclustererrors "github.com/codeready-toolchain/devcluster/pkg/errors"
	"github.com/codeready-toolchain/devcluster/pkg/ibmcloud"
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

// GetHandlerClusters returns not deleted Cluster resources for the given zone
func (r *ClusterRequest) GetHandlerClusters(ctx *gin.Context) {
	zone := ctx.Query("zone")
	clusters, err := cluster.DefaultClusterService.GetClusters(zone)
	if err != nil {
		log.Error(ctx, err, "error fetching clusters")
		devclustererrors.AbortWithError(ctx, http.StatusInternalServerError, err, "error fetching clusters")
		return
	}
	ctx.JSON(http.StatusOK, clusters)
}

// GetHandlerZones returns Zones resource
func (r *ClusterRequest) GetHandlerZones(ctx *gin.Context) {
	zones, err := cluster.DefaultClusterService.GetZones()
	if err != nil {
		log.Error(ctx, err, "error fetching zones")
		devclustererrors.AbortWithError(ctx, http.StatusInternalServerError, err, "error fetching zones")
		return
	}
	ctx.JSON(http.StatusOK, r.filterZones(ctx, zones))
}

var allowedDCs = map[string]bool{"wdc04": true, "wdc06": true, "wdc07": true, "che01": true, "fra02": true, "fra04": true, "fra05": true, "ams03": true}

// filterZones returns the filtered array of the zones/DCs allowed to be used by the client
func (r *ClusterRequest) filterZones(ctx *gin.Context, zones []ibmcloud.Location) []ibmcloud.Location {
	result := make([]ibmcloud.Location, 0, len(allowedDCs))
	for _, z := range zones {
		if allowedDCs[z.ID] {
			result = append(result, z)
		}
	}
	if len(result) < len(allowedDCs) {
		log.Error(ctx, errors.New("not all allowed zones are available"), fmt.Sprintf("allowed zones: %v; available zones: %v", allowedDCs, zones))
	}
	return result
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
