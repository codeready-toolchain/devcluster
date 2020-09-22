package ibmcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/codeready-toolchain/devcluster/pkg/auth"
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

var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func generatePassword(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
