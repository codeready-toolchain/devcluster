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

type TestIBMCloudSuite struct {
	test.UnitTestSuite
	mockConfig *MockConfig
}

func TestRunIBMCloudSuite(t *testing.T) {
	suite.Run(t, &TestIBMCloudSuite{UnitTestSuite: test.UnitTestSuite{}, mockConfig: &MockConfig{}})
}

func (s *TestIBMCloudSuite) TestGetZones() {
	cl := s.newClient(s.T())
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
}

func (s *TestIBMCloudSuite) TestGetCluster() {
	cl := s.newClient(s.T())
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
}

func (s *TestIBMCloudSuite) TestCreateCluster() {
	cl := s.newClient(s.T())
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
			BodyString(`{"id": "some-id"}`)

		id, err := cl.CreateCluster("john", "zone-1", false)
		require.NoError(t, err)
		assert.Equal(t, "some-id", id)
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
		assert.Equal(t, "some-id", id)
	})
}

func (s *TestIBMCloudSuite) TestDeleteCluster() {
	log.Init("devcluster-testing")
	cl := s.newClient(s.T())
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
}

func (s *TestIBMCloudSuite) TestToken() {
	s.T().Run("refresh", func(t *testing.T) {
		cl := s.newClient(t)
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

func (s *TestIBMCloudSuite) newClient(t *testing.T) *Client {
	cl := NewClient(s.mockConfig)
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
