package ibmcloud

import (
	"fmt"
	"testing"

	devclustererr "github.com/codeready-toolchain/devcluster/pkg/errors"
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

func (s *TestUserSuite) TestCloudDirectoryUser() {
	cl := newClient(s.T(), s.mockConfig)
	s.T().Run("Create OK", func(t *testing.T) {
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

	s.T().Run("Delete OK", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://us-south.appid.cloud.ibm.com").
			Delete("management/v4/9876543210/cloud_directory/remove/13579").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			Persist().
			Reply(204)

		err := cl.DeleteCloudDirectoryUser("13579")
		require.NoError(t, err)
	})
}

const iamSingleUserExample = `
{
    "total_results": 1,
    "limit": 100,
    "resources": [
        {
            "id": "75fb826aff7b418a9c79915c90258bb8",
            "iam_id": "5drS2pDaiG-a3659d3e-24e5-4973-b55a-68177850cae9",
            "user_id": "dev-cluster-user-5193",
            "state": "ACTIVE",
            "email": "devcluster-user-1@redhat.com"
        }
    ]
}`

const iamNoUsersExample = `
{
    "total_results": 0,
    "limit": 100,
    "resources": []
}`

const iamMultipleUsersExample = `
{
    "total_results": 2,
    "limit": 100,
    "resources": [
        {
            "id": "75fb826aff7b418a9c79915c90258bb8",
            "iam_id": "5drS2pDaiG-a3659d3e-24e5-4973-b55a-68177850cae9",
            "user_id": "dev-cluster-user-5193",
            "state": "ACTIVE",
            "email": "devcluster-user-1@redhat.com"
        },
        {
            "id": "86fb826aff7b418a9c79915c90258bb8",
            "iam_id": "1wfS2pDaiG-a3659d3e-24e5-4973-b55a-68177850cae9",
            "user_id": "dev-cluster-user-5193",
            "state": "ACTIVE",
            "email": "devcluster-user-1@redhat.com"
        }
    ]
}`

func (s *TestUserSuite) TestIAMUser() {
	cl := newClient(s.T(), s.mockConfig)
	s.T().Run("Get IAM User by user_id OK", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://user-management.cloud.ibm.com").
			Get(fmt.Sprintf("v2/accounts/%s/users", s.mockConfig.GetIBMCloudAccountID())).
			MatchParam("user_id", "dev-cluster-user-5193").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			MatchHeader("Accept", "application/json").
			Persist().
			Reply(200).
			BodyString(iamSingleUserExample)

		user, err := cl.GetIAMUserByUserID("dev-cluster-user-5193")
		require.NoError(t, err)
		assert.Equal(t, "75fb826aff7b418a9c79915c90258bb8", user.ID)
		assert.Equal(t, "dev-cluster-user-5193", user.UserID)
		assert.Equal(t, "5drS2pDaiG-a3659d3e-24e5-4973-b55a-68177850cae9", user.IAMID)
		assert.Equal(t, "devcluster-user-1@redhat.com", user.Email)
	})

	s.T().Run("Get IAM User by user_id not found", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://user-management.cloud.ibm.com").
			Get(fmt.Sprintf("v2/accounts/%s/users", s.mockConfig.GetIBMCloudAccountID())).
			MatchParam("user_id", "dev-cluster-user-5193").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			MatchHeader("Accept", "application/json").
			Persist().
			Reply(200).
			BodyString(iamNoUsersExample)

		_, err := cl.GetIAMUserByUserID("dev-cluster-user-5193")
		assert.Equal(t, devclustererr.NewNotFoundError("IAM user with user_id=dev-cluster-user-5193 not found", iamNoUsersExample), err)
	})

	s.T().Run("Too many IAM users causes Not Found", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://user-management.cloud.ibm.com").
			Get(fmt.Sprintf("v2/accounts/%s/users", s.mockConfig.GetIBMCloudAccountID())).
			MatchParam("user_id", "dev-cluster-user-5193").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			MatchHeader("Accept", "application/json").
			Persist().
			Reply(200).
			BodyString(iamMultipleUsersExample)

		_, err := cl.GetIAMUserByUserID("dev-cluster-user-5193")
		assert.Equal(t, devclustererr.NewInternalServerError("too many IAM users with user_id=dev-cluster-user-5193", iamMultipleUsersExample), err)
	})

	s.T().Run("Delete IAM User OK", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://user-management.cloud.ibm.com").
			Delete(fmt.Sprintf("v2/accounts/%s/users/08531", s.mockConfig.GetIBMCloudAccountID())).
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			Persist().
			Reply(204)

		err := cl.DeleteIAMUser("08531")
		require.NoError(t, err)
	})
}

func (s *TestUserSuite) TestAccessPolicy() {
	cl := newClient(s.T(), s.mockConfig)
	s.T().Run("Create OK", func(t *testing.T) {
		defer gock.OffAll()

		gock.New("https://iam.cloud.ibm.com").
			Post("v1/policies").
			MatchHeader("Authorization", "Bearer "+cl.token.AccessToken).
			MatchHeader("Accept", "application/json").
			MatchHeader("Content-Type", "application/json").
			JSON(fmt.Sprintf(AccessPolicyTemplate, "12345", "54321", "135790")).
			Persist().
			Reply(201).
			BodyString(`{"id": "some-id"}`)

		id, err := cl.CreateAccessPolicy("54321", "12345", "135790")
		require.NoError(t, err)
		assert.Equal(t, "some-id", id)
	})
}
