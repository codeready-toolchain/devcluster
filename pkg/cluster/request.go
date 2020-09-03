package cluster

import (
	"sort"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

var registry *RequestRegistry = &RequestRegistry{Requests: make(map[string]*Request)}

// RequestTopic represents a cluster request without details about actual clusters provisioned for this request
type RequestTopic struct {
	ID          string
	Requested   int // Number of clusters requested
	Created     time.Time
	Status      string
	RequestedBy string
}

type Request struct {
	RequestTopic `json:",inline"`
	Clusters     map[string]*Cluster // Clusters by ID
}

type Cluster struct {
	ID     string
	Name   string
	URL    string
	Status string
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

	// FIXME mock data
	id := uuid.NewV4().String()
	r.Clusters[id] = &Cluster{
		ID:     id,
		Name:   "cluster1",
		Status: "provisioning",
		URL:    "https://console.openshift-1.com",
	}
	id = uuid.NewV4().String()
	r.Clusters[id] = &Cluster{
		ID:     id,
		Name:   "cluster2",
		Status: "provisioning",
		URL:    "https://console.openshift-2.com",
	}
	registry.Add(r)
	return r
}

// RequestRegistry represents a registry of all cluster resources
type RequestRegistry struct {
	mux      sync.RWMutex
	Requests map[string]*Request
}

func (r *RequestRegistry) Add(req *Request) {
	defer r.mux.Unlock()
	r.mux.Lock()
	r.Requests[req.ID] = req
}

func (r *RequestRegistry) Remove(id string) {
	defer r.mux.Unlock()
	r.mux.Lock()
	r.Requests[id] = nil
}

func (r *RequestRegistry) Get(id string) *Request {
	defer r.mux.RUnlock()
	r.mux.RLock()
	return r.Requests[id]
}

// AllRequests returns all Requests from the registry
func (r *RequestRegistry) AllRequests() []Request {
	defer r.mux.RUnlock()
	r.mux.RLock()
	rs := make([]Request, 0, len(r.Requests))
	for _, req := range r.Requests {
		rs = append(rs, *req)
	}
	return rs
}

// RequestTopics returns all existing Cluster Request Topics from the default registry  sorted by creation time
func RequestTopics() []RequestTopic {
	rs := registry.AllRequests()
	var tpc topics //:= make([]RequestTopic, 0, len(rs))
	for _, req := range rs {
		tpc = append(tpc, req.RequestTopic)
	}
	sort.Sort(tpc)
	return tpc
}

func RequestByID(id string) *Request {
	return registry.Get(id)
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
