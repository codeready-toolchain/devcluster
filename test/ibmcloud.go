package test

import (
	"errors"
	"sync"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/codeready-toolchain/devcluster/pkg/ibmcloud"
)

type MockIBMCloudClient struct {
	mux            sync.RWMutex
	clustersByID   map[string]*ibmcloud.Cluster
	clustersByName map[string]*ibmcloud.Cluster
}

func NewMockIBMCloudClient() *MockIBMCloudClient {
	return &MockIBMCloudClient{
		clustersByName: make(map[string]*ibmcloud.Cluster),
		clustersByID:   make(map[string]*ibmcloud.Cluster),
	}
}

func (c *MockIBMCloudClient) CreateCluster(name string) (string, error) {
	defer c.mux.Unlock()
	c.mux.Lock()
	if c.clustersByName[name] != nil {
		return "", errors.New("cluster already exist")
	}
	newCluster := &ibmcloud.Cluster{
		ID:          uuid.NewV4().String(),
		Name:        name,
		Region:      "dc",
		CreatedDate: time.Now().String(),
		State:       "deploying",
		Ingress:     ibmcloud.Ingress{},
	}
	c.clustersByName[name] = newCluster
	c.clustersByID[newCluster.ID] = newCluster
	return newCluster.ID, nil
}

func (c *MockIBMCloudClient) GetCluster(id string) (*ibmcloud.Cluster, error) {
	defer c.mux.RUnlock()
	c.mux.RLock()
	return c.clustersByID[id], nil
}

func (c *MockIBMCloudClient) UpdateCluster(cluster ibmcloud.Cluster) error {
	defer c.mux.Unlock()
	c.mux.Lock()
	found := c.clustersByID[cluster.ID]
	if found == nil {
		return errors.New("cluster not found")
	}
	found.State = cluster.State
	found.Ingress = cluster.Ingress
	return nil
}
