package ibmcloud

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/alexeykazakov/devcluster/pkg/rest"

	"github.com/pkg/errors"
)

// Tokens obtains an access token
// Returns the access token string and the time when the token is going to expire
func Token(apikey string) (string, time.Time, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.PostForm("https://iam.cloud.ibm.com/identity/token", url.Values{
		"grant_type": {"urn:ibm:params:oauth:grant-type:apikey"},
		"apikey":     {apikey},
	})
	if err != nil {
		return "", time.Time{}, err
	}

	defer rest.CloseResponse(res)
	if res.StatusCode != http.StatusOK {
		bodyString := rest.ReadBody(res.Body)
		return "", time.Time{}, errors.Errorf("unable to obtain access token from IBM Cloud. Response status: %s. Response body: %s", res.Status, bodyString)
	}
	tokenSet, err := ReadTokenSet(res)
	if err != nil {
		return "", time.Time{}, err
	}
	if tokenSet.AccessToken == "" {
		return "", time.Time{}, errors.New("unable to obtain access token from IBM Cloud. Access Token is missing in the response")
	}
	expiry := time.Now().Add((time.Duration(tokenSet.ExpiresIn) - 60) * time.Second) // Subtract 1 minute from the expiry
	return tokenSet.AccessToken, expiry, nil
}

// TokenSet represents a set of Access and Refresh tokens
type TokenSet struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// ReadTokenSet extracts json with token data from the response
func ReadTokenSet(res *http.Response) (*TokenSet, error) {
	buf := new(bytes.Buffer)
	io.Copy(buf, res.Body)
	jsonString := strings.TrimSpace(buf.String())
	return ReadTokenSetFromJson(jsonString)
}

// ReadTokenSetFromJson parses json with a token set
func ReadTokenSetFromJson(jsonString string) (*TokenSet, error) {
	var token TokenSet
	err := json.Unmarshal([]byte(jsonString), &token)
	if err != nil {
		return nil, errors.Wrapf(err, "error when unmarshal json with access token %s ", jsonString)
	}
	return &token, nil
}
