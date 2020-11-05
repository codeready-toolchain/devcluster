package ibmcloud

import (
	"fmt"
	"testing"

	"github.com/codeready-toolchain/devcluster/pkg/errors"
	"github.com/codeready-toolchain/devcluster/pkg/log"
	"github.com/codeready-toolchain/devcluster/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

type TestClusterSuite struct {
	test.UnitTestSuite
	mockConfig *MockConfig
}

func TestRunClusterSuite(t *testing.T) {
	suite.Run(t, &TestClusterSuite{UnitTestSuite: test.UnitTestSuite{}, mockConfig: &MockConfig{}})
}

func (s *TestClusterSuite) TestGetZones() {
	cl := newClient(s.T(), s.mockConfig)
	s.T().Run("OK", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://containers.cloud.ibm.com").
			Get("global/v1/locations").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			Persist().
			Reply(200).
			BodyString(`[
{"id":"sng01","name": "sng01","kind": "dc","display_name": "Singapore 01"},
{"id":"lon06","name": "lon06","kind": "dc","display_name": "London 06"},
{"id":"hou","name": "hou","kind": "metro","geography": "na","display_name": "Houston"}]`)

		zones, err := cl.GetZones()
		require.NoError(t, err)
		assert.Equal(t, []Location{
			{
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
		}, zones)
	})

	s.T().Run("Error getting zones", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://containers.cloud.ibm.com").
			Get("global/v1/locations").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			Persist().
			Reply(504).
			BodyString(`something went wrong`).
			SetHeader("X-Request-Id", "1234509876")

		_, err := cl.GetZones()
		require.EqualError(t, err, "unable to get zones. x-request-id: 1234509876, Response status: 504 Gateway Timeout. Response body: something went wrong")
	})
}

func (s *TestClusterSuite) TestGetCluster() {
	cl := newClient(s.T(), s.mockConfig)
	s.T().Run("OK", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://containers.cloud.ibm.com").
			Get("global/v2/getCluster").
			MatchParam("cluster", "some-id").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			Persist().
			Reply(200).
			BodyString(`{"id": "some-id", "name": "some-name"}`)

		cluster, err := cl.GetCluster("some-id")
		require.NoError(t, err)
		assert.Equal(t, &Cluster{
			Name: "some-name",
			ID:   "some-id",
		}, cluster)
	})

	s.T().Run("Not Found", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://containers.cloud.ibm.com").
			Get("global/v2/getCluster").
			MatchParam("cluster", "some-id").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			Persist().
			Reply(404)

		_, err := cl.GetCluster("some-id")
		require.Error(t, err)
		assert.Equal(t, errors.NewNotFoundError("cluster some-id not found", ""), err)
	})

	s.T().Run("Error getting cluster", func(t *testing.T) {
		gock.New("https://containers.cloud.ibm.com").
			Get("global/v2/getCluster").
			MatchParam("cluster", "some-id").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			Persist().
			Reply(401).
			BodyString(`something went wrong`).
			SetHeader("X-Request-Id", "1234509876")

		_, err := cl.GetCluster("some-id")
		require.EqualError(t, err, "unable to get cluster. x-request-id: 1234509876, Response status: 401 Unauthorized. Response body: something went wrong")
	})
}

func (s *TestClusterSuite) TestCreateCluster() {
	cl := newClient(s.T(), s.mockConfig)
	s.T().Run("Vlan is available", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://containers.cloud.ibm.com").
			Get(fmt.Sprintf("global/v1/datacenters/%s/vlans", "zone-1")).
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			Persist().
			Reply(200).
			BodyString(`[{"id":"12345","type": "private"},{"id":"54321","type":"public"}]`)
		gock.New("https://containers.cloud.ibm.com").
			Post("global/v1/clusters").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			JSON(fmt.Sprintf(ClusterConfigTemplate, "zone-1", "john", "54321", "12345", false)).
			Persist().
			Reply(201).
			BodyString(`{"id": "some-id"}`).
			SetHeader("X-Request-Id", "10293")

		id, err := cl.CreateCluster("john", "zone-1", false)
		require.NoError(t, err)
		assert.Equal(t, &IBMCloudClusterRequest{
			ClusterID: "some-id",
			RequestID: "10293",
		}, id)
	})

	s.T().Run("Vlan is not available", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://containers.cloud.ibm.com").
			Get(fmt.Sprintf("global/v1/datacenters/%s/vlans", "zone-1")).
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			Persist().
			Reply(200).
			BodyString(`[]`)
		gock.New("https://containers.cloud.ibm.com").
			Post("global/v1/clusters").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			JSON(fmt.Sprintf(ClusterConfigTemplate, "zone-1", "john", "", "", true)).
			Persist().
			Reply(201).
			BodyString(`{"id": "some-id"}`)

		id, err := cl.CreateCluster("john", "zone-1", true)
		require.NoError(t, err)
		assert.Equal(t, &IBMCloudClusterRequest{
			ClusterID: "some-id",
		}, id)
	})

	s.T().Run("Error when creating a cluster", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://containers.cloud.ibm.com").
			Get(fmt.Sprintf("global/v1/datacenters/%s/vlans", "zone-1")).
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			Persist().
			Reply(200).
			BodyString(`[{"id":"12345","type": "private"},{"id":"54321","type":"public"}]`)
		gock.New("https://containers.cloud.ibm.com").
			Post("global/v1/clusters").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			JSON(fmt.Sprintf(ClusterConfigTemplate, "zone-1", "john", "54321", "12345", false)).
			Persist().
			Reply(500).
			SetHeader("X-Request-Id", "1234509876").
			BodyString(`oopsie woopsie`)

		_, err := cl.CreateCluster("john", "zone-1", false)
		require.EqualError(t, err, "unable to create cluster. x-request-id: 1234509876, Response status: 500 Internal Server Error. Response body: oopsie woopsie")
	})
}

func (s *TestClusterSuite) TestDeleteCluster() {
	log.Init("devcluster-testing")
	cl := newClient(s.T(), s.mockConfig)
	s.T().Run("OK", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://containers.cloud.ibm.com").
			Delete("/global/v1/clusters/some-id").
			MatchParam("deleteResources", "true").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			Persist().
			Reply(204)

		err := cl.DeleteCluster("some-id")
		require.NoError(t, err)
	})

	s.T().Run("Error when deleting a cluster", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://containers.cloud.ibm.com").
			Delete("/global/v1/clusters/some-id").
			MatchParam("deleteResources", "true").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			Persist().
			Reply(502).
			SetHeader("X-Request-Id", "1234509876").
			BodyString(`error deleting cluster`)

		err := cl.DeleteCluster("some-id")
		require.EqualError(t, err, "unable to delete cluster. x-request-id: 1234509876, Response status: 502 Bad Gateway. Response body: error deleting cluster")
	})
}

func (s *TestClusterSuite) TestToken() {
	s.T().Run("refresh", func(t *testing.T) {
		cl := newClient(s.T(), s.mockConfig)
		setNewToken(t, cl, 1598983098) // Expired

		var newExpiration int64
		newExpiration = 1798983098
		defer gock.OffAll()
		gock.New("https://iam.cloud.ibm.com").
			Post("identity/token").
			Persist().
			Reply(200).
			BodyString(fmt.Sprintf(`{"access_token": "new-token","refresh_token":"qwerty","ims_user_id": 4778951,"token_type": "Bearer","expires_in": 3600,"expiration": %d,"scope":"ibm openid"}`, newExpiration))

		token, err := cl.Token()
		require.NoError(t, err)
		require.NotNil(t, token)
		assert.Equal(t, TokenSet{
			AccessToken:  "new-token",
			RefreshToken: "qwerty",
			ExpiresIn:    3600,
			Expiration:   newExpiration,
			TokenType:    "Bearer",
		}, token)
	})
}

func newClient(t *testing.T, c Configuration) *Client {
	cl := NewClient(c)
	setNewToken(t, cl, 1998983098)
	return cl
}

func setNewToken(t *testing.T, cl *Client, expiration int64) {
	cl.token = nil
	defer gock.OffAll()
	gock.New("https://iam.cloud.ibm.com").
		Post("identity/token").
		Persist().
		Reply(200).
		BodyString(fmt.Sprintf(`{"access_token": "abc","refresh_token":"qwerty","ims_user_id": 4778951,"token_type": "Bearer","expires_in": 3600,"expiration": %d,"scope":"ibm openid"}`, expiration))

	token, err := cl.Token()
	require.NoError(t, err)
	require.NotNil(t, token)
	assert.Equal(t, TokenSet{
		AccessToken:  "abc",
		RefreshToken: "qwerty",
		ExpiresIn:    3600,
		Expiration:   expiration,
		TokenType:    "Bearer",
	}, token)
}

type MockConfig struct {
}

func (c *MockConfig) GetIBMCloudAPIKey() string {
	return "secretkey"
}

func (c *MockConfig) GetIBMCloudApiCallRetrySec() int {
	return 1
}

func (c *MockConfig) GetIBMCloudApiCallTimeoutSec() int {
	return 100
}

func (c *MockConfig) GetIBMCloudAccountID() string {
	return "0123456789"
}

func (c *MockConfig) GetIBMCloudTenantID() string {
	return "9876543210"
}

func (c *MockConfig) GetIBMCloudIDPName() string {
	return "devcluster"
}
