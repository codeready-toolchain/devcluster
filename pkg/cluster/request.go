package cluster

import (
	"sort"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"
)

var registry *RequestRegistry = &RequestRegistry{Requests: make(map[string]*Request)}

type Request struct {
	ID string
	// Number of clusters requested
	Requested   int
	Created     time.Time
	State       string
	RequestedBy string
}

func NewRequest(requestedBy string, n int) *Request {
	r := &Request{
		ID:          uuid.NewV4().String(),
		Requested:   n,
		Created:     time.Now(),
		State:       "provisioning",
		RequestedBy: requestedBy,
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
func (r *RequestRegistry) AllRequests() requests {
	defer r.mux.RUnlock()
	r.mux.RLock()
	rs := make([]Request, 0, len(r.Requests))
	for _, req := range r.Requests {
		rs = append(rs, *req)
	}
	return rs
}

// AllRequests returns all existing Cluster Requests from the default registry  sorted by creation time
func AllRequests() []Request {
	rs := registry.AllRequests()
	sort.Sort(rs)
	return rs
}

type requests []Request

func (s requests) Len() int {
	return len(s)
}

func (s requests) Less(i, j int) bool {
	return s[i].Created.Before(s[j].Created)
}

func (s requests) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
