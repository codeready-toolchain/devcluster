package test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/codeready-toolchain/devcluster/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/codeready-toolchain/devcluster/pkg/cluster"
	"github.com/codeready-toolchain/devcluster/pkg/configuration"
	"github.com/codeready-toolchain/devcluster/pkg/ibmcloud"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"k8s.io/apimachinery/pkg/util/wait"
)

type TestIntegrationSuite struct {
	IntegrationTestSuite
}

func TestRunDTestIntegrationSuite(t *testing.T) {
	suite.Run(t, &TestIntegrationSuite{IntegrationTestSuite{}})
}

func (s *TestIntegrationSuite) newRequest(service *cluster.ClusterService, n int) cluster.Request {
	req, err := service.CreateNewRequest("johnsmith@domain.com", n, "lon06", 100)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "johnsmith@domain.com", req.RequestedBy)
	assert.Equal(s.T(), n, req.Requested)
	assert.Equal(s.T(), "provisioning", req.Status)
	assert.Equal(s.T(), "lon06", req.Zone)
	return req
}

func (s *TestIntegrationSuite) TestRequestService() {
	mockClient := NewMockIBMCloudClient()
	service := &cluster.ClusterService{
		IbmCloudClient: mockClient,
		Config: &MockConfig{
			config: s.Config,
		},
	}

	request1 := s.newRequest(service, 10)
	request2 := s.newRequest(service, 10)

	s.Run("request is provisioning", func() {
		reqWithClusters1, err := waitForClustersToStartProvisioning(service, request1)
		require.NoError(s.T(), err)
		_, err = waitForClustersToStartProvisioning(service, request2)
		require.NoError(s.T(), err)

		s.Run("provisioned", func() {
			// Update all clusters as provisioned in the mock client
			s.markClustersAsProvisioned(service, mockClient, request1)

			// Check that the request is now also returned as provisioned
			_, err = waitForClustersToGetProvisioned(service, request1)
			require.NoError(s.T(), err)

			s.Run("resume provisioning", func() {
				// Delete some clusters from mongo to imitate the case when provisioning was interrupted (i.g. if pod was killed)
				// And set others in deploying state
				_, err := mongodb.Clusters().DeleteOne(
					context.Background(),
					bson.D{
						{"_id", reqWithClusters1.Clusters[0].ID},
					},
				)
				require.NoError(s.T(), err)
				_, err = mongodb.Clusters().UpdateOne(
					context.Background(),
					bson.D{
						{"_id", reqWithClusters1.Clusters[1].ID},
					},
					bson.D{
						{"$set", bson.D{
							{"status", "deploying"},
							{"url", ""},
						}},
					},
				)
				require.NoError(s.T(), err)

				// Now resume provisioning
				err = service.ResumeProvisioningRequests()
				require.NoError(s.T(), err)

				// Verify that all clusters are now provisioning
				_, err = waitForClustersToStartProvisioning(service, request2)
				require.NoError(s.T(), err)

				// Update all clusters as provisioned in the mock client
				s.markClustersAsProvisioned(service, mockClient, request2)

				// Verify that all clusters are now provisioned
				_, err = waitForClustersToGetProvisioned(service, request2)
				require.NoError(s.T(), err)
			})
		})
	})

	s.Run("get zones", func() {
		zones, err := service.GetZones()
		require.NoError(s.T(), err)
		expected, err := service.IbmCloudClient.GetZones()
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), zones)
		assert.Equal(s.T(), expected, zones)
	})

	s.Run("delete cluster", func() {
		// Provision some clusters
		_, serv, req, reqWithClusters := s.provisionClusters(3)

		// Now delete one
		toDelete := reqWithClusters.Clusters[1]
		err := service.DeleteCluster(toDelete.ID)
		require.NoError(s.T(), err)

		// Check the deleted cluster
		result, err := serv.GetRequestWithClusters(req.ID)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), cluster.Cluster{
			ID:        toDelete.ID,
			RequestID: req.ID,
			Name:      toDelete.Name,
			URL:       toDelete.URL,
			Status:    "deleted",
			Error:     "",
		}, result.Clusters[1])
	})
}

func (s *TestIntegrationSuite) provisionClusters(n int) (*MockIBMCloudClient, *cluster.ClusterService, cluster.Request, cluster.RequestWithClusters) {
	mockClient := NewMockIBMCloudClient()
	service := &cluster.ClusterService{
		IbmCloudClient: mockClient,
		Config: &MockConfig{
			config: s.Config,
		},
	}

	req := s.newRequest(service, n)
	_, err := waitForClustersToStartProvisioning(service, req)
	require.NoError(s.T(), err)
	s.markClustersAsProvisioned(service, mockClient, req)
	r, err := waitForClustersToGetProvisioned(service, req)
	require.NoError(s.T(), err)

	return mockClient, service, req, r
}

func (s *TestIntegrationSuite) markClustersAsProvisioned(service *cluster.ClusterService, client *MockIBMCloudClient, request cluster.Request) {
	// Update all clusters as provisioned in the mock client
	r, err := service.GetRequestWithClusters(request.ID)
	require.NoError(s.T(), err)
	for _, c := range r.Clusters {
		err := client.UpdateCluster(ibmcloud.Cluster{
			ID:      c.ID,
			State:   "normal",
			Ingress: ibmcloud.Ingress{Hostname: fmt.Sprintf("prefix-%s", c.Name)},
		})
		require.NoError(s.T(), err)
	}
}

var retryInterval = 100 * time.Millisecond
var timeout = 5 * time.Second

func waitForClustersToStartProvisioning(service *cluster.ClusterService, request cluster.Request) (cluster.RequestWithClusters, error) {
	var req cluster.RequestWithClusters
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		fmt.Println("Wait for clusters to start provisioning")
		r, err := service.GetRequestWithClusters(request.ID)
		if err != nil {
			return false, err
		}
		if r == nil {
			fmt.Println("Request not found")
			return false, nil
		}
		if len(r.Clusters) != request.Requested {
			fmt.Printf("Number of clusters in Request: %d\n", len(r.Clusters))
			return false, nil
		}
		if r.Zone != request.Zone {
			return false, errors.New("zone doesn't match")
		}
		if r.RequestedBy != request.RequestedBy {
			return false, errors.New("requestedBy doesn't match")
		}
		if r.Created < 1 && r.Created > time.Now().Unix() {
			return false, errors.New("invalid created time")
		}
		if r.DeleteInHours != request.DeleteInHours {
			return false, errors.New("deleteInHours doesn't match")
		}
		for _, c := range r.Clusters {
			ok := c.Status == "deploying" &&
				c.RequestID == request.ID &&
				c.Error == "" &&
				c.URL == "" &&
				strings.Contains(c.Name, "redhat-")
			if !ok {
				fmt.Printf("Found clusters: %v\n", r.Clusters)
				return false, nil
			}
		}
		req = *r
		return true, nil
	})
	return req, err
}

func waitForClustersToGetProvisioned(service *cluster.ClusterService, request cluster.Request) (cluster.RequestWithClusters, error) {
	var req cluster.RequestWithClusters
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		fmt.Println("Wait for clusters to get provisioned")
		r, err := service.GetRequestWithClusters(request.ID)
		if err != nil {
			return false, err
		}
		if r == nil {
			fmt.Println("Request not found")
			return false, nil
		}
		if r.Status != "ready" || len(r.Clusters) != request.Requested {
			fmt.Printf("Found request: %v\n", r)
			return false, nil
		}
		if r.Zone != request.Zone {
			return false, errors.New("zone doesn't match")
		}
		if r.RequestedBy != request.RequestedBy {
			return false, errors.New("requestedBy doesn't match")
		}
		if r.Created < 1 && r.Created > time.Now().Unix() {
			return false, errors.New("invalid created time")
		}
		if r.DeleteInHours != request.DeleteInHours {
			return false, errors.New("deleteInHours doesn't match")
		}
		for _, c := range r.Clusters {
			ok := c.Status == "normal" &&
				c.RequestID == request.ID &&
				c.Error == "" &&
				c.URL == fmt.Sprintf("https://console-openshift-console.prefix-%s", c.Name) &&
				strings.Contains(c.Name, "redhat-")
			if !ok {
				fmt.Printf("Found clusters: %v\n", r.Clusters)
				return false, nil
			}
		}
		req = *r
		return true, nil
	})
	return req, err
}

type MockConfig struct {
	config *configuration.Config
}

func (c *MockConfig) GetIBMCloudAPIKey() string {
	return c.config.GetIBMCloudAPIKey()
}

func (c *MockConfig) GetIBMCloudApiCallRetrySec() int {
	return 1
}
