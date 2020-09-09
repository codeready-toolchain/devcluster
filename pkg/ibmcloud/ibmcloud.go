package ibmcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/codeready-toolchain/devcluster/pkg/rest"

	"github.com/pkg/errors"
)

type Configuration interface {
	GetIBMCloudAPIKey() string
	GetIBMCloudApiCallRetrySec() int
}

type ICClient interface {
	CreateCluster(name string) (string, error)
	GetCluster(id string) (*Cluster, error)
}

type Client struct {
	config   Configuration
	token    *TokenSet
	tokenMux sync.RWMutex
}

func NewClient(config Configuration) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) GetToken() TokenSet {
	defer c.tokenMux.RUnlock()
	c.tokenMux.RLock()
	return *c.token
}

// Token returns IBM Cloud Token.
// If the token is expired or not obtained yet it will obtain a new one.
func (c *Client) Token() (TokenSet, error) {
	c.tokenMux.RLock()
	if tokenExpired(c.token) {
		c.tokenMux.RUnlock()
		c.tokenMux.Lock()
		defer c.tokenMux.Unlock()
		if tokenExpired(c.token) {
			var err error
			c.token, err = c.obtainNewToken()
			if err != nil {
				return TokenSet{}, err
			}
		}
		return *c.token, nil
	}
	defer c.tokenMux.RUnlock()
	return *c.token, nil
}

// tokenExpired return false if the token is not nil and good for at least one more minute
func tokenExpired(token *TokenSet) bool {
	return token == nil || time.Now().After(time.Unix(token.Expiration-60, 0))
}

const ClusterConfigTemplate = `
{
  "dataCenter": "wdc04",
  "disableAutoUpdate": true,
  "machineType": "b3c.4x16",
  "masterVersion": "4.4_openshift",
  "name": "%s",
  "publicVlan": "2940148",
  "privateVlan": "2940150",
  "workerNum": 2
}`

type ID struct {
	ID string `json:"id"`
}

// CreateCluster creates a cluster
// Returns the cluster ID
func (c *Client) CreateCluster(name string) (string, error) {
	token, err := c.Token()
	if err != nil {
		return "", err
	}
	body := bytes.NewBuffer([]byte(fmt.Sprintf(ClusterConfigTemplate, name)))
	req, err := http.NewRequest("POST", "https://containers.cloud.ibm.com/global/v1/clusters", body)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "unable to create cluster")
	}
	defer rest.CloseResponse(res)
	bodyString := rest.ReadBody(res.Body)
	if res.StatusCode != http.StatusCreated {
		return "", errors.Errorf("unable to create cluster. Response status: %s. Response body: %s", res.Status, bodyString)
	}

	var idObj ID
	err = json.Unmarshal([]byte(bodyString), &idObj)
	if err != nil {
		return "", errors.Wrapf(err, "error when unmarshal json with cluster ID %s ", bodyString)
	}
	return idObj.ID, nil
}

type Cluster struct {
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Region            string  `json:"region"`
	CreatedDate       string  `json:"createdDate"`
	MasterKubeVersion string  `json:"masterKubeVersion"`
	WorkerCount       int     `json:"workerCount"`
	Location          string  `json:"location"`
	Datacenter        string  `json:"datacenter"`
	State             string  `json:"state"`
	Type              string  `json:"type"`
	Crn               string  `json:"crn"`
	Ingress           Ingress `json:"ingress"`
}

type Ingress struct {
	Hostname string `json:"hostname"`
}

// GetCluster fetches the cluster with the given ID/name
func (c *Client) GetCluster(id string) (*Cluster, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://containers.cloud.ibm.com/global/v2/getCluster?cluster=%s", id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get cluster")
	}
	defer rest.CloseResponse(res)
	bodyString := rest.ReadBody(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unable to get cluster. Response status: %s. Response body: %s", res.Status, bodyString)
	}

	var cluster Cluster
	err = json.Unmarshal([]byte(bodyString), &cluster)
	if err != nil {
		return nil, errors.Wrapf(err, "error when unmarshal json with cluster %s ", bodyString)
	}
	return &cluster, nil
}

// obtainNewToken obtains an access token
// Returns the access token string and the time when the token is going to expire
func (c *Client) obtainNewToken() (*TokenSet, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.PostForm("https://iam.cloud.ibm.com/identity/token", url.Values{
		"grant_type": {"urn:ibm:params:oauth:grant-type:apikey"},
		"apikey":     {c.config.GetIBMCloudAPIKey()},
	})
	if err != nil {
		return nil, err
	}

	defer rest.CloseResponse(res)
	if res.StatusCode != http.StatusOK {
		bodyString := rest.ReadBody(res.Body)
		return nil, errors.Errorf("unable to obtain access token from IBM Cloud. Response status: %s. Response body: %s", res.Status, bodyString)
	}
	tokenSet, err := readTokenSet(res)
	if err != nil {
		return nil, err
	}
	if tokenSet.AccessToken == "" {
		return nil, errors.New("unable to obtain access token from IBM Cloud. Access Token is missing in the response")
	}
	return tokenSet, nil
}

// TokenSet represents a set of Access and Refresh tokens
type TokenSet struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Expiration   int64  `json:"expiration"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// readTokenSet extracts json with token data from the response
func readTokenSet(res *http.Response) (*TokenSet, error) {
	buf := new(bytes.Buffer)
	io.Copy(buf, res.Body)
	jsonString := strings.TrimSpace(buf.String())
	return readTokenSetFromJson(jsonString)
}

// readTokenSetFromJson parses json with a token set
func readTokenSetFromJson(jsonString string) (*TokenSet, error) {
	var token TokenSet
	err := json.Unmarshal([]byte(jsonString), &token)
	if err != nil {
		return nil, errors.Wrapf(err, "error when unmarshal json with access token %s ", jsonString)
	}
	return &token, nil
}
