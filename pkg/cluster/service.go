package cluster

import (
	"fmt"
	"hash/fnv"
	"time"

	"github.com/alexeykazakov/devcluster/pkg/ibmcloud"
	"github.com/alexeykazakov/devcluster/pkg/log"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// Request represents a cluster request
type Request struct {
	ID          string
	Requested   int // Number of clusters requested
	Created     string
	Status      string
	Error       string
	RequestedBy string
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
	ibmCloudClient *ibmcloud.Client
}

func InitDefaultClusterService(config ibmcloud.Configuration) {
	DefaultClusterService = &ClusterService{
		ibmCloudClient: ibmcloud.NewClient(config),
	}
}

func (s *ClusterService) Requests() ([]Request, error) {
	return getRequests()
}

func (s *ClusterService) RequestWithClusters(requestID string) (*RequestWithClusters, error) {
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

// StartNewRequest creates a new request and starts provisioning clusters
func (s *ClusterService) StartNewRequest(requestedBy string, n int) (Request, error) {
	r := Request{
		ID:          uuid.NewV4().String(),
		Requested:   n,
		Created:     time.Now().String(),
		Status:      "provisioning",
		RequestedBy: requestedBy,
	}

	err := insertRequest(r)
	if err != nil {
		return Request{}, errors.Wrap(err, "unable to start new request")
	}
	for i := 0; i < r.Requested; i++ {
		go func() {
			err := provisionRequest(r)
			if err != nil {
				log.Error(nil, err, "unable to provision a cluster")
			}
		}()
	}

	return r, nil
}

func hash(s string) uint32 {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		log.Error(nil, err, "unable to generate cluster name")
	}
	return h.Sum32()
}

// provisionRequest provisions one new cluster
func provisionRequest(r Request) error {
	name := fmt.Sprintf("redhat-%d", hash(uuid.NewV4().String()))
	var id string
	var err error
	// Try to create a cluster. If failing then we will make six attempts for one minute before giving up.
	for i := 0; i < 6; i++ {
		id, err = DefaultClusterService.ibmCloudClient.CreateCluster(name)
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
	in3Hours := time.Now().Add(3 * time.Hour)
	for time.Now().Before(in3Hours) { // timeout in three hours
		c, err := DefaultClusterService.ibmCloudClient.GetCluster(id)
		if err != nil {
			log.Error(nil, err, "unable to get cluster")
			err := clusterFailed(err, id, name, r.ID)
			if err != nil {
				return err
			}
			// Do not exist. Try again in 30 seconds.
		} else {
			clusterToAdd := convertCluster(*c, r.ID)
			err := insertCluster(clusterToAdd)
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
		time.Sleep(30 * time.Second)
	}
	return nil
}

func clusterReady(c Cluster) bool {
	return c.Status == "normal" && c.URL != ""
}

func clusterFailed(clErr error, id, name, reqID string) error {
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
	clToUpdate.Status = "failed"
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
