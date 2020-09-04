package cluster

import (
	"sort"
	"sync"

	"github.com/alexeykazakov/devcluster/pkg/ibmcloud"
)

var DefaultRegistry *RequestRegistry

// RequestRegistry represents a registry of all cluster resources
type RequestRegistry struct {
	mux      sync.RWMutex
	Requests map[string]*Request
	client   *ibmcloud.Client
}

func InitDefaultRegistry(config ibmcloud.Configuration) {
	DefaultRegistry = &RequestRegistry{
		Requests: make(map[string]*Request),
		client:   ibmcloud.NewClient(config),
	}
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
	rs := DefaultRegistry.AllRequests()
	var tpc topics //:= make([]RequestTopic, 0, len(rs))
	for _, req := range rs {
		tpc = append(tpc, req.RequestTopic)
	}
	sort.Sort(tpc)
	return tpc
}

func RequestByID(id string) *Request {
	return DefaultRegistry.Get(id)
}
