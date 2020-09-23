package ibmcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/codeready-toolchain/devcluster/pkg/auth"
	devclustererr "github.com/codeready-toolchain/devcluster/pkg/errors"
	"github.com/codeready-toolchain/devcluster/pkg/rest"

	"github.com/pkg/errors"
)

const (
	apiRegion = "us-south"
)

type CloudDirectoryUser struct {
	ID        string  `json:"id"`
	Username  string  `json:"userName"`
	Emails    []Value `json:"emails"`
	ProfileID string  `json:"profileId"`
	Password  string
}

type Value struct {
	Value string `json:"value"`
}

func (u *CloudDirectoryUser) Email() string {
	if len(u.Emails) > 0 {
		return u.Emails[0].Value
	}
	return ""
}

const CloudDirectoryUserTemplate = `{"active":true, "emails":[{"value":"%s", "primary":true}], "userName":"%s", "password":"%s"}`

// CreateCloudDirectoryUser creates a new cloud directory user with generated username, email, and password.
func (c *Client) CreateCloudDirectoryUser() (*CloudDirectoryUser, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	username := auth.GenerateShortID("dev")
	email := fmt.Sprintf("%s.redhat.com", username)
	password := generatePassword(8)
	body := bytes.NewBuffer([]byte(fmt.Sprintf(CloudDirectoryUserTemplate, email, username, password)))
	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s.appid.cloud.ibm.com/management/v4/%s/cloud_directory/sign_up?shouldCreateProfile=true&language=en", apiRegion, c.config.GetIBMCloudTenantID()), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create cloud directory user")
	}
	defer rest.CloseResponse(res)
	bodyString := rest.ReadBody(res.Body)
	if res.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("unable to create cloud directory user. Response status: %s. Response body: %s", res.Status, bodyString)
	}

	var userObj CloudDirectoryUser
	err = json.Unmarshal([]byte(bodyString), &userObj)
	if err != nil {
		return nil, errors.Wrapf(err, "error when unmarshal json with cloud directory user: %s ", bodyString)
	}
	userObj.Password = password // Set the generated password before returning the user struct
	return &userObj, nil
}

// DeleteCloudDirectoryUser deletes the cloud directory user with the given ID.
func (c *Client) DeleteCloudDirectoryUser(id string) error {
	token, err := c.Token()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://%s.appid.cloud.ibm.com/management/v4/%s/cloud_directory/remove/%s", apiRegion, c.config.GetIBMCloudTenantID(), id), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "unable to delete cloud directory user")
	}
	defer rest.CloseResponse(res)
	bodyString := rest.ReadBody(res.Body)
	if res.StatusCode != http.StatusNoContent {
		return errors.Errorf("unable to delete cloud directory user. Response status: %s. Response body: %s", res.Status, bodyString)
	}

	return nil
}

type IAMUserResult struct {
	TotalResults int       `json:"total_results"`
	Limit        int       `json:"limit"`
	Resources    []IAMUser `json:"resources"`
}

type IAMUser struct {
	ID     string `json:"id"`
	IAMID  string `json:"iam_id"`
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// GetIAMUserByUserID fetches the AIM user with the corresponding user_id. Returns a Not Found Error if the user is not found.
func (c *Client) GetIAMUserByUserID(userID string) (*IAMUser, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://user-management.cloud.ibm.com/v2/accounts/%s/users", c.config.GetIBMCloudAccountID()), nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("user_id", userID)
	req.URL.RawQuery = params.Encode()
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get IAM users")
	}
	defer rest.CloseResponse(res)
	bodyString := rest.ReadBody(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unable to get IAM users. Response status: %s. Response body: %s", res.Status, bodyString)
	}

	var iamRes IAMUserResult
	err = json.Unmarshal([]byte(bodyString), &iamRes)
	if err != nil {
		return nil, errors.Wrapf(err, "error when unmarshal json with IAM users: %s ", bodyString)
	}
	if len(iamRes.Resources) == 0 {
		return nil, devclustererr.NewNotFoundError(fmt.Sprintf("IAM user with user_id=%s not found", userID), bodyString)
	}
	if len(iamRes.Resources) > 1 {
		return nil, devclustererr.NewInternalServerError(fmt.Sprintf("too many IAM users with user_id=%s", userID), bodyString)
	}

	return &iamRes.Resources[0], nil
}

// DeleteIAMUser deletes the AIM user with the corresponding id (IAMUser.ID).
func (c *Client) DeleteIAMUser(id string) error {
	token, err := c.Token()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://user-management.cloud.ibm.com/v2/accounts/%s/users/%s", c.config.GetIBMCloudAccountID(), id), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "unable to delete IAM user")
	}
	defer rest.CloseResponse(res)
	bodyString := rest.ReadBody(res.Body)
	if res.StatusCode != http.StatusNoContent {
		return errors.Errorf("unable to delete IAM user. Response status: %s. Response body: %s", res.Status, bodyString)
	}

	return nil
}

const AccessPolicyTemplate = `
{
   "type": "access",
   "subjects": [{"attributes": [{"name": "iam_id", "value": "%s"}]}],
   "roles": [
       {"role_id": "crn:v1:bluemix:public:iam::::role:Viewer"},
       {"role_id": "crn:v1:bluemix:public:iam::::serviceRole:Manager"}
   ],
   "resources": [
       {
           "attributes": [
               {"name": "accountId", "operator": "stringEquals", "value": "%s"},
               {"name": "serviceName", "value": "containers-kubernetes"},
               {"name": "serviceInstance", "value": "%s"}
           ]
       }
   ]
}`

type AccessPolicy struct {
	TotalResults int       `json:"total_results"`
	Limit        int       `json:"limit"`
	Resources    []IAMUser `json:"resources"`
}

// CreateAccessPolicy creates an access policy for the cluster and assigns it to the user.
// Returns the id of the created policy.
func (c *Client) CreateAccessPolicy(accountID, iamID, clusterID string) (string, error) {
	token, err := c.Token()
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer([]byte(fmt.Sprintf(AccessPolicyTemplate, iamID, accountID, clusterID)))
	req, err := http.NewRequest("POST", "https://iam.cloud.ibm.com/v1/policies", body)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "unable to create access policy")
	}
	defer rest.CloseResponse(res)
	bodyString := rest.ReadBody(res.Body)
	if res.StatusCode != http.StatusCreated {
		return "", errors.Errorf("unable to create access policy. Response status: %s. Response body: %s", res.Status, bodyString)
	}

	var idObj ID
	err = json.Unmarshal([]byte(bodyString), &idObj)
	if err != nil {
		return "", errors.Wrapf(err, "error when unmarshal json with access policy ID %s ", bodyString)
	}
	return idObj.ID, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func generatePassword(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
