package ibmcloud

import (
	"errors"
	"fmt"
	"sync"
	"time"

	devclustererr "github.com/codeready-toolchain/devcluster/pkg/errors"
	"github.com/codeready-toolchain/devcluster/pkg/ibmcloud"

	uuid "github.com/satori/go.uuid"
)

type MockIBMCloudClient struct {
	clusterMux     sync.RWMutex
	cldUserMux     sync.RWMutex
	policyMux      sync.RWMutex
	clustersByID   map[string]*ibmcloud.Cluster
	clustersByName map[string]*ibmcloud.Cluster
	cldUserByID    map[string]*ibmcloud.CloudDirectoryUser
	aimUserByID    map[string]*ibmcloud.IAMUser
	policyByID     map[string]*string
}

func NewMockIBMCloudClient() *MockIBMCloudClient {
	return &MockIBMCloudClient{
		clustersByName: make(map[string]*ibmcloud.Cluster),
		clustersByID:   make(map[string]*ibmcloud.Cluster),
		cldUserByID:    make(map[string]*ibmcloud.CloudDirectoryUser),
		aimUserByID:    make(map[string]*ibmcloud.IAMUser),
		policyByID:     make(map[string]*string),
	}
}

func (c *MockIBMCloudClient) GetZones() ([]ibmcloud.Location, error) {
	return []ibmcloud.Location{{
		ID:          "lon06",
		Name:        "lon06",
		Kind:        "dc",
		DisplayName: "London 06",
	},
		{
			ID:          "sng01",
			Name:        "sng01",
			Kind:        "dc",
			DisplayName: "Singapore 01",
		},
	}, nil
}

func (c *MockIBMCloudClient) GetVlans(zone string) ([]ibmcloud.Vlan, error) {
	return []ibmcloud.Vlan{}, nil
}

func (c *MockIBMCloudClient) CreateCluster(name, zone string, _ bool) (string, error) {
	defer c.clusterMux.Unlock()
	c.clusterMux.Lock()
	if c.clustersByName[name] != nil {
		return "", errors.New("cluster already exist")
	}
	cluster := &ibmcloud.Cluster{
		ID:          uuid.NewV4().String(),
		Name:        name,
		Region:      zone,
		CreatedDate: time.Now().String(),
		State:       "deploying",
		Ingress:     ibmcloud.Ingress{},
	}
	c.clustersByName[name] = cluster
	c.clustersByID[cluster.ID] = cluster
	return cluster.ID, nil
}

func (c *MockIBMCloudClient) GetCluster(id string) (*ibmcloud.Cluster, error) {
	defer c.clusterMux.RUnlock()
	c.clusterMux.RLock()
	clst := c.clustersByID[id]
	if clst == nil {
		return nil, devclustererr.NewNotFoundError(fmt.Sprintf("cluster %s not found", id), "")
	}
	return c.clustersByID[id], nil
}

func (c *MockIBMCloudClient) DeleteCluster(id string) error {
	defer c.clusterMux.Unlock()
	c.clusterMux.Lock()
	cluster := c.clustersByID[id]
	if cluster != nil {
		c.clustersByID[id] = nil
		c.clustersByName[cluster.Name] = nil
	} else {
		return devclustererr.NewNotFoundError(fmt.Sprintf("cluster %s not found", id), "")
	}
	return nil
}

func (c *MockIBMCloudClient) UpdateCluster(cluster ibmcloud.Cluster) error {
	defer c.clusterMux.Unlock()
	c.clusterMux.Lock()
	found := c.clustersByID[cluster.ID]
	if found == nil {
		return errors.New("cluster not found")
	}
	found.State = cluster.State
	found.Ingress = cluster.Ingress
	return nil
}

func (c *MockIBMCloudClient) CreateCloudDirectoryUser(username string) (*ibmcloud.CloudDirectoryUser, error) {
	defer c.cldUserMux.Unlock()
	c.cldUserMux.Lock()
	user := &ibmcloud.CloudDirectoryUser{
		ID:        uuid.NewV4().String(),
		Username:  username,
		Emails:    []ibmcloud.Value{{uuid.NewV4().String()}},
		ProfileID: uuid.NewV4().String(),
		Password:  uuid.NewV4().String(),
	}
	c.cldUserByID[user.ID] = user
	aimUser := &ibmcloud.IAMUser{
		ID:     uuid.NewV4().String(),
		IAMID:  uuid.NewV4().String(),
		UserID: user.Username,
		Email:  user.Email(),
	}
	c.aimUserByID[aimUser.UserID] = aimUser
	return user, nil
}

func (c *MockIBMCloudClient) UpdateCloudDirectoryUserPassword(id string) (*ibmcloud.CloudDirectoryUser, error) {
	defer c.cldUserMux.Unlock()
	c.cldUserMux.Lock()
	found := c.cldUserByID[id]
	if found == nil {
		return nil, errors.New("user not found")
	}
	found.Password = uuid.NewV4().String()
	return found, nil
}

func (c *MockIBMCloudClient) GetIAMUserByUserID(userID string) (*ibmcloud.IAMUser, error) {
	defer c.cldUserMux.RUnlock()
	c.cldUserMux.RLock()
	found := c.aimUserByID[userID]
	if found == nil {
		return nil, errors.New("user not found")
	}
	return found, nil
}

func (c *MockIBMCloudClient) CreateAccessPolicy(_, _, _ string) (string, error) {
	defer c.policyMux.Unlock()
	c.policyMux.Lock()
	id := uuid.NewV4().String()
	c.policyByID[id] = &id
	return id, nil
}

func (c *MockIBMCloudClient) DeleteAccessPolicy(id string) error {
	defer c.policyMux.Unlock()
	c.policyMux.Lock()
	c.policyByID[id] = nil
	return nil
}
