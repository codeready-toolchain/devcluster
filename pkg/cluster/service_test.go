package cluster_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	devclustererr "github.com/codeready-toolchain/devcluster/pkg/errors"

	"github.com/codeready-toolchain/devcluster/pkg/cluster"
	"github.com/codeready-toolchain/devcluster/pkg/configuration"
	"github.com/codeready-toolchain/devcluster/pkg/ibmcloud"
	"github.com/codeready-toolchain/devcluster/pkg/mongodb"
	"github.com/codeready-toolchain/devcluster/test"
	ibmcloudmock "github.com/codeready-toolchain/devcluster/test/ibmcloud"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"k8s.io/apimachinery/pkg/util/wait"
)

type TestIntegrationSuite struct {
	test.IntegrationTestSuite
}

func TestRunDTestIntegrationSuite(t *testing.T) {
	suite.Run(t, &TestIntegrationSuite{test.IntegrationTestSuite{}})
}

func (s *TestIntegrationSuite) TestRequestClusters() {
	s.Run("request is provisioning", func() {
		service, mockClient := s.prepareService()
		s.newUsers(service, 50)
		request1 := s.newRequest(service, 10, 100)
		request2 := s.newRequest(service, 10, 100)

		reqWithClusters1, err := waitForClustersToStartProvisioning(service, request1)
		require.NoError(s.T(), err)
		// Check the cluster were created in ibm cloud
		for _, c := range reqWithClusters1.Clusters {
			_, err := mockClient.GetCluster(c.ID)
			assert.NoError(s.T(), err)
		}

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
}

func (s *TestIntegrationSuite) TestGetZones() {
	service, _ := s.prepareService()
	s.Run("get zones OK", func() {
		zones, err := service.GetZones()
		require.NoError(s.T(), err)
		expected, err := service.IbmCloudClient.GetZones()
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), zones)
		assert.Equal(s.T(), expected, zones)
	})
}

func (s *TestIntegrationSuite) TestDeleteCluster() {
	service, cl := s.prepareService()
	s.newUsers(service, 10)
	s.Run("delete cluster OK", func() {
		// Provision some clusters
		req, reqWithClusters := s.provisionClusters(service, cl, 3, 100)

		// Now delete one
		toDelete := reqWithClusters.Clusters[1]
		err := service.DeleteCluster(toDelete.ID)
		require.NoError(s.T(), err)

		// Check the deleted cluster
		result, err := service.GetRequestWithClusters(req.ID)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), cluster.Cluster{
			ID:        toDelete.ID,
			RequestID: req.ID,
			Name:      toDelete.Name,
			URL:       toDelete.URL,
			Status:    "deleted",
			Error:     "",
		}, result.Clusters[1])

		// Check the cluster was deleted in ibm cloud
		_, err = service.IbmCloudClient.GetCluster(toDelete.ID)
		require.Error(s.T(), err)
		assert.True(s.T(), devclustererr.IsNotFound(err))
	})
}

func (s *TestIntegrationSuite) TestExpiredClusters() {
	service, cl := s.prepareService()
	users := s.newUsers(service, 6)
	s.Run("delete expired clusters OK", func() {
		// 1. Provision two requests. One is expired and the other one is not.
		reqExpired, reqExpiredWithClusters := s.provisionClusters(service, cl, 3, 0)
		req, reqWithClusters := s.provisionClusters(service, cl, 3, 100)

		// 1.1. Verify that all clusters have assigned users
		assertClusterHasAssignedUser := func(c cluster.Cluster) *cluster.User {
			u, err := cluster.GetUserByClusterID(c.ID)
			require.NoError(s.T(), err)
			assert.Equal(s.T(), c.ID, u.ClusterID)
			assert.NotEmpty(s.T(), u.Password)
			assert.NotEmpty(s.T(), u.PolicyID)
			assert.True(s.T(), cl.AccessPolicyExists(u.PolicyID))
			assert.Empty(s.T(), c.User.Recycled) // The user has not been recycled yet
			return u
		}
		foundUsers := make(map[string]*string, 0)
		var policiesToBeDeleted = make([]string, 0, 3)
		usersToBeRecycled := make([]*cluster.User, 0, 3)
		for _, c := range reqExpiredWithClusters.Clusters {
			u := assertClusterHasAssignedUser(c)
			foundUsers[u.ID] = &u.ID
			policiesToBeDeleted = append(policiesToBeDeleted, u.PolicyID)
			usersToBeRecycled = append(usersToBeRecycled, u)
		}
		for _, c := range reqWithClusters.Clusters {
			u := assertClusterHasAssignedUser(c)
			foundUsers[u.ID] = &u.ID
		}
		assert.Len(s.T(), foundUsers, len(users))
		for _, u := range users {
			assert.NotNil(s.T(), foundUsers[u.ID])
		}

		// 2. Start deleting clusters.
		beforeDeleting := time.Now().Unix()
		service.StartDeletingExpiredClusters(1)

		// 3. Check the expired one is deleted and the other one is not.
		deletedReq, err := waitForRequest(service, reqExpired, requestExpired, clustersDeleted, usersRecycled)
		require.NoError(s.T(), err)
		// Check the cluster were deleted from ibm cloud
		for _, c := range deletedReq.Clusters {
			_, err := service.IbmCloudClient.GetCluster(c.ID)
			require.Error(s.T(), err)
			assert.True(s.T(), devclustererr.IsNotFound(err))
		}
		// And the expired clusters do not have assigned users
		for _, c := range reqExpiredWithClusters.Clusters {
			_, err := cluster.GetUserByClusterID(c.ID)
			require.EqualError(s.T(), err, fmt.Sprintf("404 Not Found: no User with cluster_id %s found: mongo: no documents in result", c.ID))
		}
		// And all the users from the expired clusters are recycled
		currentUsers, err := service.Users()
		require.NoError(s.T(), err)
		for _, u := range usersToBeRecycled {
			for _, cu := range currentUsers {
				if cu.ID == u.ID {
					assert.True(s.T(), cu.Recycled >= beforeDeleting && cu.Recycled <= time.Now().Unix()) // Recycle timestamp is set
				}
			}
		}
		// All the polices have been deleted
		for _, policy := range policiesToBeDeleted {
			assert.False(s.T(), cl.AccessPolicyExists(policy))
		}

		// the other one is still ready and the clusters still exist in IC
		_, err = waitForClustersToGetProvisioned(service, req)
		require.NoError(s.T(), err)
		for _, c := range reqWithClusters.Clusters {
			_, err := service.IbmCloudClient.GetCluster(c.ID)
			require.NoError(s.T(), err)
		}
		// the clusters are still assigned
		foundUsers = make(map[string]*string, 0)
		for _, c := range reqWithClusters.Clusters {
			u := assertClusterHasAssignedUser(c)
			foundUsers[u.ID] = &u.ID
		}
		assert.Len(s.T(), foundUsers, 3)

		s.Run("re-use recycled users", func() {
			// Add one new user with the recycle timestamp not set so it should be used first before the recycled ones
			newUsers, err := service.CreateUsers(1, 1000)
			require.NoError(s.T(), err)
			usersToBeAssigned := append([]*cluster.User{&newUsers[0]}, usersToBeRecycled...)

			// Provision new clusters which should use the new user and the recycled ones which were returned to the pull after the first request expired
			_, reqWithClusters := s.provisionClusters(service, cl, 4, 100)

			// Verify that all the clusters use the recycled users
			for i, c := range reqWithClusters.Clusters {
				assert.Equal(s.T(), usersToBeAssigned[i].ID, c.User.ID)
				assert.Equal(s.T(), c.User.ClusterID, c.ID)
				assert.Equal(s.T(), usersToBeAssigned[i].Email, c.User.Email)
				assert.Equal(s.T(), usersToBeAssigned[i].CloudDirectID, c.User.CloudDirectID)
				assert.NotEqual(s.T(), usersToBeAssigned[i].PolicyID, c.User.PolicyID) // different policy
				assert.NotEmpty(s.T(), c.User.PolicyID)
				if i == 0 {
					// new user
					assert.Empty(s.T(), c.User.Recycled)
					assert.NotEmpty(s.T(), c.User.Password)
				} else {
					// recycled user
					assert.True(s.T(), c.User.Recycled >= beforeDeleting && c.User.Recycled <= time.Now().Unix()) // Recycle timestamp is set
					assert.NotEqual(s.T(), usersToBeAssigned[i].Password, c.User.Password)                        // different password
				}
			}
		})
	})
}

func (s *TestIntegrationSuite) TestUsers() {
	s.Run("request new users OK", func() {
		mockClient := ibmcloudmock.NewMockIBMCloudClient()
		service := &cluster.ClusterService{
			IbmCloudClient: mockClient,
			Config: &MockConfig{
				config: s.Config,
			},
		}

		assertUsers := func(users []cluster.User, err error) {
			require.NoError(s.T(), err)
			require.Len(s.T(), users, 3)
			for i := 0; i < 3; i++ {
				assert.Equal(s.T(), fmt.Sprintf("rh-dev-%d", 1001+i), users[i].ID)
				assert.NotEmpty(s.T(), users[i].Email)
				assert.NotEmpty(s.T(), users[i].Password)
				assert.NotEmpty(s.T(), users[i].CloudDirectID)
				assert.Empty(s.T(), users[i].PolicyID)
				assert.Empty(s.T(), users[i].ClusterID)
				assert.Empty(s.T(), users[i].Recycled)
			}
		}

		// Request 3 new users and assert the result
		assertUsers(service.CreateUsers(3, 1000))

		s.Run("get users", func() {
			// assert the available users
			assertUsers(service.Users())
		})
	})
}

func (s *TestIntegrationSuite) newRequest(service *cluster.ClusterService, n int, deleteIn int) cluster.Request {
	req, err := service.CreateNewRequest("johnsmith@domain.com", n, "lon06", deleteIn, false)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "johnsmith@domain.com", req.RequestedBy)
	assert.Equal(s.T(), n, req.Requested)
	assert.Equal(s.T(), "provisioning", req.Status)
	assert.Equal(s.T(), "lon06", req.Zone)

	return req
}

func (s *TestIntegrationSuite) newUsers(service *cluster.ClusterService, n int) []cluster.User {
	users, err := service.CreateUsers(n, 0)
	require.NoError(s.T(), err)
	return users
}

func (s *TestIntegrationSuite) prepareService() (*cluster.ClusterService, *ibmcloudmock.MockIBMCloudClient) {
	mockClient := ibmcloudmock.NewMockIBMCloudClient()
	service := &cluster.ClusterService{
		IbmCloudClient: mockClient,
		Config: &MockConfig{
			config: s.Config,
		},
	}
	return service, mockClient
}

func (s *TestIntegrationSuite) provisionClusters(service *cluster.ClusterService, client *ibmcloudmock.MockIBMCloudClient, n, deleteIn int) (cluster.Request, cluster.RequestWithClusters) {
	req := s.newRequest(service, n, deleteIn)
	_, err := waitForClustersToStartProvisioning(service, req)
	require.NoError(s.T(), err)
	s.markClustersAsProvisioned(service, client, req)
	r, err := waitForClustersToGetProvisioned(service, req)
	require.NoError(s.T(), err)

	return req, r
}

func (s *TestIntegrationSuite) markClustersAsProvisioned(service *cluster.ClusterService, client *ibmcloudmock.MockIBMCloudClient, request cluster.Request) {
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

type RequestCriterion func(req *cluster.RequestWithClusters) (bool, error)

func requestReady(req *cluster.RequestWithClusters) (bool, error) {
	return req.Status == "ready", nil
}

func requestExpired(req *cluster.RequestWithClusters) (bool, error) {
	return req.Status == "expired", nil
}

func clustersDeploying(req *cluster.RequestWithClusters) (bool, error) {
	for _, c := range req.Clusters {
		ok := c.Status == "deploying" &&
			c.RequestID == req.ID &&
			c.Error == "" &&
			c.URL == "" &&
			strings.Contains(c.Name, "redhat-")
		if !ok {
			fmt.Printf("Found clusters: %v\n", req.Clusters)
			return false, nil
		}
	}
	return true, nil
}

func clustersDeleted(req *cluster.RequestWithClusters) (bool, error) {
	for _, c := range req.Clusters {
		ok := c.Status == "deleted" &&
			c.RequestID == req.ID &&
			c.Error == "" &&
			strings.Contains(c.Name, "redhat-")
		if !ok {
			fmt.Printf("Found clusters: %v\n", req.Clusters)
			return false, nil
		}
	}
	return true, nil
}

func clustersReady(req *cluster.RequestWithClusters) (bool, error) {
	for _, c := range req.Clusters {
		ok := c.Status == "normal" &&
			c.RequestID == req.ID &&
			c.Error == "" &&
			c.URL == fmt.Sprintf("https://console-openshift-console.prefix-%s", c.Name) &&
			strings.Contains(c.Name, "redhat-")
		if !ok {
			fmt.Printf("Found clusters: %v\n", req.Clusters)
			return false, nil
		}
	}
	return true, nil
}

func usersAssigned(req *cluster.RequestWithClusters) (bool, error) {
	for _, c := range req.Clusters {
		if c.Status != "deleted" {
			user, err := cluster.GetUserByClusterID(c.ID)
			if err != nil {
				return false, err
			}
			if user.Password == "" ||
				user.PolicyID == "" ||
				user.CloudDirectID == "" ||
				user.Email == "" ||
				*user != c.User {
				return false, errors.New(fmt.Sprintf("unexpected cluster user assigned to cluster. Expected: %v; Actual: %v", user, c.User))
			}
		}
	}
	return true, nil
}

func usersRecycled(req *cluster.RequestWithClusters) (bool, error) {
	for _, c := range req.Clusters {
		if c.Status == "deleted" {
			user, err := cluster.GetUserByClusterID(c.ID)
			if err == nil {
				return false, errors.New(fmt.Sprintf("cluster has an assigned user: %v", user))
			}
			if err.Error() != fmt.Sprintf("404 Not Found: no User with cluster_id %s found: mongo: no documents in result", c.ID) {
				return false, errors.New(fmt.Sprintf("unexpected error: %s", err.Error()))
			}
			empty := cluster.User{}
			if c.User != empty {
				return false, errors.New(fmt.Sprintf("cluster has an assigned user: %v", c.User))
			}
		}
	}
	return true, nil
}

func waitForClustersToStartProvisioning(service *cluster.ClusterService, request cluster.Request) (cluster.RequestWithClusters, error) {
	fmt.Println("Wait for clusters to start provisioning")
	return waitForRequest(service, request, clustersDeploying, usersAssigned)
}

func waitForClustersToGetProvisioned(service *cluster.ClusterService, request cluster.Request) (cluster.RequestWithClusters, error) {
	fmt.Println("Wait for clusters to get provisioned")
	return waitForRequest(service, request, requestReady, clustersReady, usersAssigned)
}

func waitForRequest(service *cluster.ClusterService, request cluster.Request, criteria ...RequestCriterion) (cluster.RequestWithClusters, error) {
	var req cluster.RequestWithClusters
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
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
		for _, match := range criteria {
			ok, err := match(r)
			if err != nil || !ok {
				fmt.Printf("Found request: %v\n", r)
				return ok, err
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

func (c *MockConfig) GetIBMCloudAccountID() string {
	return "0123456789"
}

func (c *MockConfig) GetIBMCloudTenantID() string {
	return "9876543210"
}
