package ibmcloud

import (
	"fmt"
	"testing"

	"gopkg.in/h2non/gock.v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mockConfig = &MockConfig{}

func TestCreateCluster(t *testing.T) {
	cl := newClient(t)
	t.Run("OK", func(t *testing.T) {
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
}

func TestGetCluster(t *testing.T) {
	cl := newClient(t)
	t.Run("OK", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://containers.cloud.ibm.com").
			Post("global/v1/clusters").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			JSON(fmt.Sprintf(ClusterConfigTemplate, "john")).
			Persist().
			Reply(201).
			BodyString(`{"id": "some-id"}`)

		id, err := cl.CreateCluster("john")
		require.NoError(t, err)
		assert.Equal(t, "some-id", id)
	})
}

func TestToken(t *testing.T) {
	t.Run("refresh", func(t *testing.T) {
		cl := newClient(t)
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

func newClient(t *testing.T) *Client {
	cl := NewClient(mockConfig)
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
