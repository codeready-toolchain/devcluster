package test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/codeready-toolchain/devcluster/pkg/ibmcloud"

	"github.com/codeready-toolchain/devcluster/pkg/cluster"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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
	service := &cluster.ClusterService{IbmCloudClient: mockClient}

	request1 := s.newRequest(service, 10)
	request2 := s.newRequest(service, 10)

	s.Run("request is provisioning", func() {
		s.waitForClustersToStartProvisioning(service, request1)
		s.waitForClustersToStartProvisioning(service, request2)

		s.Run("provisioned", func() {
			// Update all clusters as provisioned in the mock client
			r, err := service.GetRequestWithClusters(request1.ID)
			require.NoError(s.T(), err)
			for _, c := range r.Clusters {
				err := mockClient.UpdateCluster(ibmcloud.Cluster{
					ID:      c.ID,
					State:   "normal",
					Ingress: ibmcloud.Ingress{Hostname: fmt.Sprintf("https://console.%s", c.Name)},
				})
				require.NoError(s.T(), err)
			}

			// Check that the request is now also returned as provisioned
			s.waitForClustersToGetProvisioned(service, request1)
		})
	})
}

func (s *TestIntegrationSuite) waitForClustersToStartProvisioning(service *cluster.ClusterService, request cluster.Request) {
	waitFor(func() bool {
		r, err := service.GetRequestWithClusters(request.ID)
		require.NoError(s.T(), err)
		if len(r.Clusters) != 10 {
			return false
		}
		for _, c := range r.Clusters {
			assert.Equal(s.T(), "deploying", c.Status)
			assert.Equal(s.T(), request.ID, c.RequestID)
			assert.Equal(s.T(), "", c.Error)
			assert.Equal(s.T(), "", c.URL)
			assert.Contains(s.T(), c.Name, "redhat-")
		}
		return true
	}, time.Second)
}

func (s *TestIntegrationSuite) waitForClustersToGetProvisioned(service *cluster.ClusterService, request cluster.Request) {
	waitFor(func() bool {
		r, err := service.GetRequestWithClusters(request.ID)
		require.NoError(s.T(), err)
		if r.Status != "ready" || len(r.Clusters) != 10 {
			return false
		}
		for _, c := range r.Clusters {
			assert.Equal(s.T(), "normal", c.Status)
			assert.Equal(s.T(), request.ID, c.RequestID)
			assert.Equal(s.T(), "", c.Error)
			assert.Equal(s.T(), fmt.Sprintf("https://console.%s", c.Name), c.URL)
			assert.Contains(s.T(), c.Name, "redhat-")
		}
		return true
	}, time.Second)
}

func worker(f func() bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		if f() {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func waitFor(f func() bool, timeout time.Duration) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go worker(f, &wg)

	fmt.Printf("Wait for waitgroup (up to %s)\n", timeout)
	if waitTimeout(&wg, timeout) {
		fmt.Println("Timed out waiting for wait group")
	} else {
		fmt.Println("Wait group finished")
	}
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
