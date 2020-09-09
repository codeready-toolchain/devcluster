package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

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
	req, err := service.CreateNewRequest("johnsmith@domain.com", n)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "johnsmith@domain.com", req.RequestedBy)
	assert.Equal(s.T(), n, req.Requested)
	assert.Equal(s.T(), "provisioning", req.Status)
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

	request1 := s.newRequest(service, numberOfClustersPerReq)
	request2 := s.newRequest(service, numberOfClustersPerReq)

	s.Run("request is provisioning", func() {
		err := waitForClustersToStartProvisioning(service, request1)
		require.NoError(s.T(), err)
		err = waitForClustersToStartProvisioning(service, request2)
		require.NoError(s.T(), err)

		s.Run("provisioned", func() {
			// Update all clusters as provisioned in the mock client
			r, err := service.GetRequestWithClusters(request1.ID)
			require.NoError(s.T(), err)
			for _, c := range r.Clusters {
				err := mockClient.UpdateCluster(ibmcloud.Cluster{
					ID:      c.ID,
					State:   "normal",
					Ingress: ibmcloud.Ingress{Hostname: fmt.Sprintf("prefix-%s", c.Name)},
				})
				require.NoError(s.T(), err)
			}

			// Check that the request is now also returned as provisioned
			err = waitForClustersToGetProvisioned(service, request1)
			require.NoError(s.T(), err)
		})
	})
}

var retryInterval = 100 * time.Millisecond
var timeout = 5 * time.Second
var numberOfClustersPerReq = 10

func waitForClustersToStartProvisioning(service *cluster.ClusterService, request cluster.Request) error {
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
		if len(r.Clusters) != numberOfClustersPerReq {
			fmt.Printf("Number of clusters in Request: %d\n", len(r.Clusters))
			return false, nil
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
		return true, nil
	})
	return err
}

func waitForClustersToGetProvisioned(service *cluster.ClusterService, request cluster.Request) error {
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
		if r.Status != "ready" || len(r.Clusters) != numberOfClustersPerReq {
			fmt.Printf("Found request: %v\n", r)
			return false, nil
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
		return true, nil
	})
	return err
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
