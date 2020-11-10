package ibmcloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	devclustererr "github.com/codeready-toolchain/devcluster/pkg/errors"
	"github.com/codeready-toolchain/devcluster/pkg/log"
	"github.com/codeready-toolchain/devcluster/pkg/rest"

	"github.com/pkg/errors"
)

type Configuration interface {
	GetIBMCloudAPIKey() string
	GetIBMCloudApiCallRetrySec() int
	GetIBMCloudApiCallTimeoutSec() int
	GetIBMCloudAccountID() string
	GetIBMCloudTenantID() string
	GetIBMCloudIDPName() string
}

type ICClient interface {
	GetVlans(zone string) ([]Vlan, error)
	GetZones() ([]Location, error)
	CreateCluster(name, zone string, noSubnet bool) (*IBMCloudClusterRequest, error)
	GetCluster(id string) (*Cluster, error)
	DeleteCluster(id string) error
	CreateCloudDirectoryUser(username string) (*CloudDirectoryUser, error)
	UpdateCloudDirectoryUserPassword(id string) (*CloudDirectoryUser, error)
	GetIAMUserByUserID(userID string) (*IAMUser, error)
	CreateAccessPolicy(accountID, userID, clusterID string) (string, error)
	DeleteAccessPolicy(id string) error
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

type Vlan struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// vlanIDByType returns the vlan ID for the given type. Returns "" if there is no such type.
func vlanIDByType(vlans []Vlan, t string) string {
	for _, vlan := range vlans {
		if vlan.Type == t {
			return vlan.ID
		}
	}
	return ""
}

func responseErr(res *http.Response, message, respBody string) error {
	id := extractRequestID(res)
	return errors.Errorf("%s. x-request-id: %s, Response status: %s. Response body: %s", message, id, res.Status, respBody)
}

func extractRequestID(res *http.Response) string {
	var id string
	ids := res.Header["X-Request-Id"]
	if len(ids) == 0 {
		ids = res.Header["x-request-id"]
	}
	if len(ids) > 0 {
		id = ids[0]
	}
	return id
}

// GetVlans fetches the list of vlans available in the zone
func (c *Client) GetVlans(zone string) ([]Vlan, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://containers.cloud.ibm.com/global/v1/datacenters/%s/vlans", zone), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get vlans")
	}
	defer rest.CloseResponse(res)
	bodyString := rest.ReadBody(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, responseErr(res, "unable to get vlans", bodyString)
	}

	var vlans []Vlan
	err = json.Unmarshal([]byte(bodyString), &vlans)
	if err != nil {
		return nil, errors.Wrapf(err, "error when unmarshal json with vlans %s ", bodyString)
	}
	return vlans, nil
}

type Location struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	DisplayName string `json:"display_name"`
}

// GetZones fetches the list of zones (data centers)
func (c *Client) GetZones() ([]Location, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", "https://containers.cloud.ibm.com/global/v1/locations", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get zones")
	}
	defer rest.CloseResponse(res)
	bodyString := rest.ReadBody(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, responseErr(res, "unable to get zones", bodyString)
	}

	var locations []Location
	err = json.Unmarshal([]byte(bodyString), &locations)
	if err != nil {
		return nil, errors.Wrapf(err, "error when unmarshal json with zones %s ", bodyString)
	}
	// Return only locations with kind == "dc" (data center)
	zones := make([]Location, 0, 0)
	for _, z := range locations {
		if z.Kind == "dc" {
			zones = append(zones, z)
		}
	}
	// Sort by Display Name
	sort.SliceStable(zones, func(i, j int) bool {
		return zones[i].DisplayName < zones[j].DisplayName
	})
	return zones, nil
}

type ID struct {
	ID string `json:"id"`
}

const ClusterConfigTemplate = `
{
  "dataCenter": "%s",
  "disableAutoUpdate": true,
  "machineType": "b3c.4x16",
  "masterVersion": "4.5_openshift",
  "name": "%s",
  "publicVlan": "%s",
  "privateVlan": "%s",
  "noSubnet": %t,
  "workerNum": 2
}`

type IBMCloudClusterRequest struct {
	ClusterID   string
	RequestID   string
	PublicVlan  string
	PrivateVlan string
}

// CreateCluster creates a cluster
// Returns the cluster ID
func (c *Client) CreateCluster(name, zone string, noSubnet bool) (*IBMCloudClusterRequest, error) {
	token, err := c.Token()
	if err != nil {
		return nil, err
	}

	// Get vlans
	vlans, err := c.GetVlans(zone)
	if err != nil {
		return nil, err
	}
	private := vlanIDByType(vlans, "private")
	if private == "" {
		log.Infof(nil, "WARNING: no private vlan found for zone %s. New vlan will be created", zone)
	}
	public := vlanIDByType(vlans, "public")
	if public == "" {
		log.Infof(nil, "WARNING: no public vlan found for zone %s. New vlan will be created", zone)
	}

	body := bytes.NewBuffer([]byte(fmt.Sprintf(ClusterConfigTemplate, zone, name, public, private, noSubnet)))
	req, err := http.NewRequest("POST", "https://containers.cloud.ibm.com/global/v1/clusters", body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create cluster")
	}
	defer rest.CloseResponse(res)
	bodyString := rest.ReadBody(res.Body)
	if res.StatusCode != http.StatusCreated {
		return nil, responseErr(res, "unable to create cluster", bodyString)
	}

	var idObj ID
	err = json.Unmarshal([]byte(bodyString), &idObj)
	if err != nil {
		return nil, responseErr(res, "error when unmarshal json with cluster ID", bodyString)
	}

	if public == "" || private == "" {
		// VLANs were just created. Obtain them so we can store them in the cluster object in the DB
		vlans, err := c.GetVlans(zone)
		if err != nil {
			return nil, err
		}
		private = vlanIDByType(vlans, "private")
		if private == "" {
			log.Infof(nil, "WARNING: no private vlan found for zone %s even after creating a cluster in that zone", zone)
		}
		public = vlanIDByType(vlans, "public")
		if public == "" {
			log.Infof(nil, "WARNING: no public vlan found for zone %s even after creating a cluster in that zone", zone)
		}
	}
	return &IBMCloudClusterRequest{
		ClusterID:   idObj.ID,
		RequestID:   extractRequestID(res),
		PublicVlan:  public,
		PrivateVlan: private,
	}, nil
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
	MasterURL         string  `json:"masterURL"`
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
	if res.StatusCode == http.StatusNotFound {
		return nil, devclustererr.NewNotFoundError(fmt.Sprintf("cluster %s not found", id), bodyString)
	}
	if res.StatusCode != http.StatusOK {
		return nil, responseErr(res, "unable to get cluster", bodyString)
	}

	var cluster Cluster
	err = json.Unmarshal([]byte(bodyString), &cluster)
	if err != nil {
		return nil, errors.Wrapf(err, "error when unmarshal json with cluster %s ", bodyString)
	}
	return &cluster, nil
}

// DeleteCluster deletes the cluster with the given ID/name
func (c *Client) DeleteCluster(id string) error {
	token, err := c.Token()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://containers.cloud.ibm.com/global/v1/clusters/%s?deleteResources=true", id), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "unable to delete cluster")
	}
	defer rest.CloseResponse(res)
	bodyString := rest.ReadBody(res.Body)
	if res.StatusCode == http.StatusNotFound {
		return devclustererr.NewNotFoundError(fmt.Sprintf("cluster %s not found", id), "")
	}
	if res.StatusCode != http.StatusNoContent {
		return responseErr(res, "unable to delete cluster", bodyString)
	}
	return nil
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
		return nil, responseErr(res, "unable to obtain access token from IBM Cloud", bodyString)
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
	_, err := io.Copy(buf, res.Body)
	if err != nil {
		return nil, err
	}
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
