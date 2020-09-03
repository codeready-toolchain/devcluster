package ibmcloud_test

import (
	"testing"
	"time"

	"github.com/alexeykazakov/devcluster/pkg/ibmcloud"

	"gopkg.in/h2non/gock.v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://iam.cloud.ibm.com").
			Post("identity/token").
			Persist().
			Reply(200).
			BodyString(`{"access_token": "abc","refresh_token":"qwerty","ims_user_id": 4778951,"token_type": "Bearer","expires_in": 3600,"expiration": 1598983098,"scope":"ibm openid"}`)

		token, expiry, err := ibmcloud.Token("secretkey")
		require.NoError(t, err)
		assert.Equal(t, "abc", token)
		assert.True(t, expiry.After(time.Now()) && expiry.Before(time.Now().Add(time.Duration(3600)*time.Second)))
	})
}
