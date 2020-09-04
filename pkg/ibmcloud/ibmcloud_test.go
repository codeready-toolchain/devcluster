package ibmcloud_test

import (
	"testing"

	"github.com/alexeykazakov/devcluster/pkg/ibmcloud"

	"gopkg.in/h2non/gock.v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mockConfig = &MockConfig{}

func TestToken(t *testing.T) {
	cl := ibmcloud.NewClient(mockConfig)
	t.Run("OK", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://iam.cloud.ibm.com").
			Post("identity/token").
			Persist().
			Reply(200).
			BodyString(`{"access_token": "abc","refresh_token":"qwerty","ims_user_id": 4778951,"token_type": "Bearer","expires_in": 3600,"expiration": 1598983098,"scope":"ibm openid"}`)

		token, err := cl.Token()
		require.NoError(t, err)
		require.NotNil(t, token)
		assert.Equal(t, ibmcloud.TokenSet{
			AccessToken:  "abc",
			RefreshToken: "qwerty",
			ExpiresIn:    3600,
			Expiration:   1598983098,
			TokenType:    "Bearer",
		}, token)
	})
}

type MockConfig struct {
}

func (c *MockConfig) GetIBMCloudAPIKey() string {
	return "secretkey"
}
