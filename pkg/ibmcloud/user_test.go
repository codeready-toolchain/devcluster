package ibmcloud

import (
	"testing"

	"github.com/codeready-toolchain/devcluster/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
)

type TestUserSuite struct {
	test.UnitTestSuite
	mockConfig *MockConfig
}

func TestRunUserSuite(t *testing.T) {
	suite.Run(t, &TestUserSuite{UnitTestSuite: test.UnitTestSuite{}, mockConfig: &MockConfig{}})
}

const cloudDirectoryUserExample = `
{
    "userName": "myUsername",
    "emails": [
        {
            "value": "user@domain.com",
            "primary": true
        }
    ],
    "profileId": "128a7445-1371-4cb2-8656-3e4590df132e",
    "schemas": [
        "urn:ietf:params:scim:schemas:core:2.0:User"
    ],
    "id": "1029a9cb-7b8a-4d18-ba91-fa4830ae0860"
}`

func (s *TestUserSuite) TestCreateCloudDirectoryUser() {
	cl := newClient(s.T(), s.mockConfig)
	s.T().Run("OK", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://us-south.appid.cloud.ibm.com").
			Post("management/v4/9876543210/cloud_directory/sign_up").
			MatchParam("shouldCreateProfile", "true").
			MatchParam("language", "en").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			MatchHeader("Accept", "application/json").
			MatchHeader("Content-Type", "application/json").
			Persist().
			Reply(201).
			BodyString(cloudDirectoryUserExample)

		user, err := cl.CreateCloudDirectoryUser()
		require.NoError(t, err)
		assert.Equal(t, "1029a9cb-7b8a-4d18-ba91-fa4830ae0860", user.ID)
		assert.Equal(t, "myUsername", user.Username)
		assert.Equal(t, "user@domain.com", user.Email())
		assert.Equal(t, "128a7445-1371-4cb2-8656-3e4590df132e", user.ProfileID)
		assert.NotEmpty(t, user.Password)
	})
}
