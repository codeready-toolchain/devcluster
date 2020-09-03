package cluster

import uuid "github.com/satori/go.uuid"

type Request struct {
	ID string
	// Number of clusters requested
	Requested int
}

func NewRequest(n int) *Request {
	return &Request{
		ID:        uuid.NewV4().String(),
		Requested: n,
	}
}
