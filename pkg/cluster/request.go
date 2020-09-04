package cluster

import (
	"sync"
	"time"

	"github.com/alexeykazakov/devcluster/pkg/ibmcloud"
	"github.com/alexeykazakov/devcluster/pkg/log"

	uuid "github.com/satori/go.uuid"
)

// RequestTopic represents a cluster request without details about actual clusters provisioned for this request
type RequestTopic struct {
	ID          string
	Requested   int // Number of clusters requested
	Created     time.Time
	Status      string
	Error       string
	RequestedBy string
}

type Request struct {
	RequestTopic `json:",inline"`
	Clusters     map[string]*Cluster // Clusters by ID
	clusterMux   sync.RWMutex
	statusMux    sync.RWMutex
}

type Cluster struct {
	ID     string
	Name   string
	URL    string
	Status string
	Error  string
}

func NewRequest(requestedBy string, n int) *Request {
	r := &Request{
		RequestTopic: RequestTopic{
			ID:          uuid.NewV4().String(),
			Requested:   n,
			Created:     time.Now(),
			Status:      "provisioning",
			RequestedBy: requestedBy,
		},
		Clusters: make(map[string]*Cluster),
	}

	DefaultRegistry.Add(r)
	return r
}

func (r *Request) Add(c *Cluster) {
	defer r.clusterMux.Unlock()
	r.clusterMux.Lock()
	r.Clusters[c.ID] = c
}

func (r *Request) Remove(id string) {
	defer r.clusterMux.Unlock()
	r.clusterMux.Lock()
	r.Clusters[id] = nil
}

func (r *Request) Get(id string) *Cluster {
	defer r.clusterMux.RUnlock()
	r.clusterMux.RLock()
	return r.Clusters[id]
}

func (r *Request) GetClusters() []Cluster {
	clusters := make([]Cluster, 0, r.Len())
	defer r.clusterMux.RUnlock()
	r.clusterMux.RLock()
	for _, c := range r.Clusters {
		clusters = append(clusters, *c)
	}
	return clusters
}

func (r *Request) Len() int {
	defer r.clusterMux.RUnlock()
	r.clusterMux.RLock()
	return len(r.Clusters)
}

func (r *Request) SetStatus(status, error string) {
	defer r.statusMux.Unlock()
	r.statusMux.Lock()
	r.Status = status
	r.Error = error
}

func (r *Request) GetStatus() string {
	defer r.statusMux.RUnlock()
	r.statusMux.RLock()
	return r.Status
}

func (r *Request) SetStatusToSuccessIfDone() {
	defer r.statusMux.Unlock()
	r.statusMux.Lock()
	if r.Len() < r.Requested {
		return
	}
	for _, c := range r.GetClusters() {
		if !clusterReady(c) {
			return
		}
	}
	r.Status = "ready"
	r.Error = ""
}

// Start starts provisioning clusters
func (r *Request) Start() {
	for i := 0; i < r.Requested; i++ {
		go r.provision()
	}
}

type topics []RequestTopic

func (t topics) Len() int {
	return len(t)
}

func (t topics) Less(i, j int) bool {
	return t[i].Created.Before(t[j].Created)
}

func (t topics) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// provision provisions one cluster
func (r *Request) provision() {
	name := uuid.NewV4().String() // TODO generate a better name
	var id string
	var err error
	// Try to create a cluster. If failing then we will make six attempts for one minute before giving up.
	for i := 0; i < 6; i++ {
		id, err = DefaultRegistry.client.CreateCluster(name)
		if err == nil {
			break
		}
		log.Error(nil, err, "unable to create cluster")
		time.Sleep(10 * time.Second)
	}
	if err != nil {
		// Set request status to failed and break
		r.SetStatus("failed", err.Error())
		return
	}
	in3Hours := time.Now().Add(3 * time.Hour)
	for time.Now().Before(in3Hours) { // timeout in three hours
		c, err := DefaultRegistry.client.GetCluster(id)
		if err != nil {
			log.Error(nil, err, "unable to get cluster")
			r.clusterFailed(err, id, name)
			// Do not exist. Try again in 30 seconds.
		} else {
			clusterToAdd := convertCluster(*c)
			r.Add(clusterToAdd)
			if clusterReady(*clusterToAdd) { // Ready
				r.SetStatusToSuccessIfDone()
				// TODO add user
				break
			}
		}
		time.Sleep(30 * time.Second)
	}

}

func clusterReady(c Cluster) bool {
	return c.Status == "normal" && c.URL != ""
}

func (r *Request) clusterFailed(err error, id, name string) {
	clToUpdate := r.Get(id)
	if clToUpdate == nil {
		clToUpdate = &Cluster{
			ID:   id,
			Name: name,
		}
	}
	clToUpdate.Error = err.Error()
	clToUpdate.Status = "failed"
	r.Add(clToUpdate)
}

func convertCluster(from ibmcloud.Cluster) *Cluster {
	console := from.Ingress.Hostname
	if console != "" {
		console = "https://console-openshift-console." + console
	}
	return &Cluster{
		ID:     from.ID,
		URL:    console,
		Status: from.State,
		Name:   from.Name,
	}
}
