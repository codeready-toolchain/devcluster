package cluster

import (
	"fmt"
	"hash/fnv"
	"time"

	devclustererr "github.com/codeready-toolchain/devcluster/pkg/errors"
	"github.com/codeready-toolchain/devcluster/pkg/ibmcloud"
	"github.com/codeready-toolchain/devcluster/pkg/log"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// Request represents a cluster request
type Request struct {
	ID          string
	Requested   int // Number of clusters requested
	Created     int64
	Status      string
	Error       string
	RequestedBy string
	Zone        string
}

// Request represents a cluster request with detailed information about all request clusters
type RequestWithClusters struct {
	Request  `json:",inline"`
	Clusters []Cluster
}

type Cluster struct {
	ID        string
	RequestID string
	Name      string
	URL       string
	Status    string
	Error     string
}

var DefaultClusterService *ClusterService

// ClusterService represents a registry of all cluster resources
type ClusterService struct {
	IbmCloudClient ibmcloud.ICClient
	Config         ibmcloud.Configuration
}

func InitDefaultClusterService(config ibmcloud.Configuration) {
	DefaultClusterService = &ClusterService{
		IbmCloudClient: ibmcloud.NewClient(config),
		Config:         config,
	}
}

func (s *ClusterService) Requests() ([]Request, error) {
	return getAllRequests()
}

func (s *ClusterService) GetZones() ([]ibmcloud.Location, error) {
	return s.IbmCloudClient.GetZones()
}

func (s *ClusterService) GetRequestWithClusters(requestID string) (*RequestWithClusters, error) {
	request, err := getRequest(requestID)
	if err != nil {
		return nil, err
	}
	if request == nil {
		// Not found
		return nil, nil
	}
	clusters, err := getClusters(requestID)
	if err != nil {
		return nil, err
	}
	return &RequestWithClusters{
		Request:  *request,
		Clusters: clusters,
	}, nil
}

// CreateNewRequest creates a new request and starts provisioning clusters
func (s *ClusterService) CreateNewRequest(requestedBy string, n int, zone string) (Request, error) {
	r := Request{
		ID:          uuid.NewV4().String(),
		Requested:   n,
		Created:     time.Now().Unix(),
		Status:      "provisioning",
		RequestedBy: requestedBy,
		Zone:        zone,
	}

	err := insertRequest(r)
	if err != nil {
		return Request{}, errors.Wrap(err, "unable to start new request")
	}
	for i := 0; i < r.Requested; i++ {
		go func() {
			err := s.provisionNewCluster(r)
			if err != nil {
				log.Error(nil, err, "unable to provision a cluster")
			}
		}()
	}

	return r, nil
}

// DeleteCluster deletes the cluster with the given ID
func (s *ClusterService) DeleteCluster(id string) error {
	err := s.IbmCloudClient.DeleteCluster(id)
	if err != nil {
		return err
	}
	c, err := getCluster(id)
	if err != nil {
		return err
	}
	c.Error = ""
	c.Status = "deleted"
	return replaceCluster(*c)
}

// ResumeProvisioningRequests load requests that are still provisioning and wait for their clusters to be ready to update the status
func (s *ClusterService) ResumeProvisioningRequests() error {
	requests, err := getRequestsWithStatus("provisioning")
	if err != nil {
		return err
	}
	for _, r := range requests {
		resumeRequest := r // need to use a copy in goroutine
		clusters, err := getClusters(resumeRequest.ID)
		if err != nil {
			return err
		}
		// Update provisioning clusters
		for _, cluster := range clusters {
			resumeCluster := cluster // need to use a copy in goroutine
			if clusterProvisioningPending(resumeCluster) {
				go func() {
					log.Infof(nil, "resuming provisioning cluster %s", resumeCluster.Name)
					err := s.waitForClusterToBeReady(resumeRequest, resumeCluster.ID, resumeCluster.Name)
					if err != nil {
						log.Error(nil, err, "unable to provision a cluster")
					}
				}()
			}
		}
		// Start provisioning missing clusters
		count := len(clusters)
		if count < resumeRequest.Requested {
			for i := count; i < resumeRequest.Requested; i++ {
				go func() {
					err := s.provisionNewCluster(resumeRequest)
					if err != nil {
						log.Error(nil, err, "unable to provision a cluster")
					}
				}()
			}
		}
	}

	return nil
}

func hash(s string) uint32 {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		log.Error(nil, err, "unable to generate cluster name")
	}
	return h.Sum32()
}

// provisionNewCluster creates one new cluster
func (s *ClusterService) provisionNewCluster(r Request) error {
	name := fmt.Sprintf("redhat-%d", hash(uuid.NewV4().String()))
	var id string
	var err error
	// Try to create a cluster. If failing then we will make six attempts for one minute before giving up.
	for i := 0; i < 6; i++ {
		id, err = s.IbmCloudClient.CreateCluster(name, r.Zone)
		if err == nil {
			break
		}
		log.Error(nil, err, "unable to create cluster")
		time.Sleep(10 * time.Second)
	}
	if err != nil {
		// Set request status to failed and break
		r.Status = "failed"
		r.Error = err.Error()
		err := updateRequestStatus(r.ID, "failed", err.Error())
		return err
	}
	log.Infof(nil, "starting provisioning cluster %s", name)
	return s.waitForClusterToBeReady(r, id, name)
}

// waitForClusterToBeReady for the cluster to be ready
func (s *ClusterService) waitForClusterToBeReady(r Request, clusterID, clusterName string) error {
	in3Hours := time.Now().Add(3 * time.Hour)
	for time.Now().Before(in3Hours) { // timeout in three hours
		c, err := s.IbmCloudClient.GetCluster(clusterID)
		if err != nil {
			log.Error(nil, err, "unable to get cluster")
			if devclustererr.IsNotFound(err) {
				return clusterFailed(err, "deleted", clusterID, clusterName, r.ID)
			}
			err := clusterFailed(err, "failed", clusterID, clusterName, r.ID)
			if err != nil {
				return err
			}
			// Do not exist. Try again in s.config.GetIBMCloudApiCallRetrySec() seconds.
		} else {
			clusterToAdd := convertCluster(*c, r.ID)
			err := replaceCluster(clusterToAdd)
			if err != nil {
				return err
			}
			if clusterReady(clusterToAdd) { // Ready
				err := setRequestStatusToSuccessIfDone(r)
				if err != nil {
					return err
				}
				// TODO add user
				break
			}
		}
		time.Sleep(time.Duration(s.Config.GetIBMCloudApiCallRetrySec()) * time.Second)
	}
	return nil
}

func clusterReady(c Cluster) bool {
	return c.Status == "normal" && c.URL != ""
}

// clusterProvisioningPending returns true if cluster is still provisioning
func clusterProvisioningPending(c Cluster) bool {
	return !clusterReady(c) && c.Status != "failed" && c.Status != "deleted"
}

func clusterFailed(clErr error, status, id, name, reqID string) error {
	clToUpdate, err := getCluster(id)
	if err != nil {
		return err
	}
	if clToUpdate == nil {
		clToUpdate = &Cluster{
			ID:        id,
			Name:      name,
			RequestID: reqID,
		}
	}
	clToUpdate.Error = clErr.Error()
	clToUpdate.Status = status
	return replaceCluster(*clToUpdate)
}

func convertCluster(from ibmcloud.Cluster, requestID string) Cluster {
	console := from.Ingress.Hostname
	if console != "" {
		console = "https://console-openshift-console." + console
	}
	return Cluster{
		ID:        from.ID,
		URL:       console,
		Status:    from.State,
		Name:      from.Name,
		RequestID: requestID,
	}
}
