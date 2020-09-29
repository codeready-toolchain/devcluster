package cluster

import (
	"context"
	"fmt"

	"github.com/codeready-toolchain/devcluster/pkg/log"

	"github.com/codeready-toolchain/devcluster/pkg/mongodb"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func insertRequest(req Request) error {
	_, err := mongodb.ClusterRequests().InsertOne(context.Background(), convertClusterRequestToBSON(req))
	return errors.Wrap(err, "unable to insert request")
}

func getRequest(id string) (*Request, error) {
	res := mongodb.ClusterRequests().FindOne(
		context.Background(),
		bson.D{{"_id", id}},
	)
	if res == nil {
		return nil, errors.New(fmt.Sprintf("unable to find Request with such ID: %s", id))
	}
	var m bson.M
	err := res.Decode(&m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "unable to get cluster request from mongo")
	}
	r := convertBSONToRequest(m)
	return &r, nil
}

func getAllRequests() ([]Request, error) {
	return getRequests(bson.D{})
}

func getRequestsWithStatus(status string) ([]Request, error) {
	return getRequests(bson.D{
		{"status", status},
	})
}

func getRequests(d bson.D) ([]Request, error) {
	requests := make([]Request, 0, 0)
	cursor, err := mongodb.ClusterRequests().Find(context.Background(), d)
	if err != nil {
		return requests, errors.Wrap(err, "unable to load cluster requests from mongo")
	}
	var rqs []bson.M
	if err = cursor.All(context.Background(), &rqs); err != nil {
		return requests, errors.Wrap(err, "unable to load cluster requests from mongo")
	}
	for _, m := range rqs {
		requests = append(requests, convertBSONToRequest(m))
	}
	return requests, err
}

func updateRequestStatus(id, status, error string) error {
	_, err := mongodb.ClusterRequests().UpdateOne(
		context.Background(),
		bson.D{
			{"_id", id},
		},
		bson.D{
			{"$set", bson.D{
				{"status", status},
				{"error", error},
			}},
		},
	)
	return errors.Wrap(err, "unable to update request status")
}

func setRequestStatusToSuccessIfDone(req Request) error {
	clusters, err := getClusters(req.ID)
	if err != nil {
		return err
	}
	if len(clusters) < req.Requested {
		return nil
	}
	for _, c := range clusters {
		if !clusterReady(c) {
			return nil
		}
	}
	log.Infof(nil, "request %s is ready", req.ID)
	return updateRequestStatus(req.ID, StatusReady, "")
}

func replaceRequest(req Request) error {
	opts := options.Replace().SetUpsert(true)
	_, err := mongodb.ClusterRequests().ReplaceOne(
		context.Background(),
		bson.D{
			{"_id", req.ID},
		},
		convertClusterRequestToBSON(req),
		opts,
	)
	return errors.Wrap(err, "unable to replace request")
}

func replaceCluster(c Cluster) error {
	opts := options.Replace().SetUpsert(true)
	_, err := mongodb.Clusters().ReplaceOne(
		context.Background(),
		bson.D{
			{"_id", c.ID},
		},
		convertClusterToBSON(c),
		opts,
	)
	return errors.Wrap(err, "unable to replace cluster")
}

func getCluster(id string) (*Cluster, error) {
	res := mongodb.Clusters().FindOne(
		context.Background(),
		bson.D{{"_id", id}},
	)
	if res == nil {
		return nil, errors.New(fmt.Sprintf("unable to find Cluster with such ID: %s", id))
	}
	var m bson.M
	err := res.Decode(&m)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrap(err, "unable to get cluster from mongo")
	}
	c := convertBSONToCluster(m)
	return &c, nil
}

func getClusters(requestID string) ([]Cluster, error) {
	clusters := make([]Cluster, 0, 0)
	cursor, err := mongodb.Clusters().Find(
		context.Background(),
		bson.D{{"request_id", requestID}},
	)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return clusters, nil
		}
		return clusters, errors.Wrap(err, "unable to load clusters from mongo")
	}
	var cls []bson.M
	if err = cursor.All(context.Background(), &cls); err != nil {
		return clusters, errors.Wrap(err, "unable to load clusters from mongo")
	}
	for _, m := range cls {
		clusters = append(clusters, convertBSONToCluster(m))
	}
	return clusters, err
}

// getUserWithoutCluster returns the first found user with no cluster_id set
// returns an error if no user found
func getUserWithoutCluster() (*User, error) {
	return getUserByClusterID("")
}

// getUserByClusterID returns the first found user with the given cluster_id
// returns an error if no user found
func getUserByClusterID(clusterID string) (*User, error) {
	res := mongodb.Users().FindOne(
		context.Background(),
		bson.D{{"cluster_id", clusterID}},
	)
	if res == nil {
		return nil, errors.New(fmt.Sprintf("unable to find User with cluster_id: %s", clusterID))
	}
	var m bson.M
	err := res.Decode(&m)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("no User with cluster_id %s found", clusterID))
	}
	u := convertBSONToUser(m)
	return &u, nil
}

func getAllUsers() ([]User, error) {
	return getUsers(bson.D{})
}

func getUsers(d bson.D) ([]User, error) {
	users := make([]User, 0, 0)
	cursor, err := mongodb.Users().Find(context.Background(), d)
	if err != nil {
		return users, errors.Wrap(err, "unable to load users from mongo")
	}
	var usrs []bson.M
	if err = cursor.All(context.Background(), &usrs); err != nil {
		return users, errors.Wrap(err, "unable to load users from mongo")
	}
	for _, m := range usrs {
		users = append(users, convertBSONToUser(m))
	}
	return users, err
}

func replaceUser(u User) error {
	opts := options.Replace().SetUpsert(true)
	_, err := mongodb.Users().ReplaceOne(
		context.Background(),
		bson.D{
			{"_id", u.ID},
		},
		convertUserToBSON(u),
		opts,
	)
	return errors.Wrap(err, "unable to replace user")
}

func insertUser(u User) error {
	_, err := mongodb.Users().InsertOne(context.Background(), convertUserToBSON(u))
	return errors.Wrap(err, "unable to insert user")
}

func convertBSONToRequest(m bson.M) Request {
	return Request{
		ID:            string(fmt.Sprintf("%v", m["_id"])),
		RequestedBy:   string(fmt.Sprintf("%v", m["requested_by"])),
		Created:       int64(m["created"].(int64)),
		Error:         string(fmt.Sprintf("%v", m["error"])),
		Requested:     int(m["requested"].(int32)),
		Status:        string(fmt.Sprintf("%v", m["status"])),
		Zone:          string(fmt.Sprintf("%v", m["zone"])),
		DeleteInHours: int(m["delete_in_hours"].(int32)),
		NoSubnet:      m["no_subnet"].(bool),
	}
}

func convertClusterRequestToBSON(req Request) bson.D {
	return bson.D{
		{"_id", req.ID},
		{"status", req.Status},
		{"requested", req.Requested},
		{"error", req.Error},
		{"created", req.Created},
		{"requested_by", req.RequestedBy},
		{"zone", req.Zone},
		{"delete_in_hours", req.DeleteInHours},
		{"no_subnet", req.NoSubnet},
	}
}

func convertBSONToCluster(m bson.M) Cluster {
	return Cluster{
		ID:        string(fmt.Sprintf("%v", m["_id"])),
		RequestID: string(fmt.Sprintf("%v", m["request_id"])),
		URL:       string(fmt.Sprintf("%v", m["url"])),
		Error:     string(fmt.Sprintf("%v", m["error"])),
		Name:      string(fmt.Sprintf("%v", m["name"])),
		Status:    string(fmt.Sprintf("%v", m["status"])),
	}
}

func convertClusterToBSON(c Cluster) bson.D {
	return bson.D{
		{"_id", c.ID},
		{"status", c.Status},
		{"name", c.Name},
		{"error", c.Error},
		{"url", c.URL},
		{"request_id", c.RequestID},
	}
}

func convertBSONToUser(m bson.M) User {
	return User{
		ID:            string(fmt.Sprintf("%v", m["_id"])),
		CloudDirectID: string(fmt.Sprintf("%v", m["cloud_direct_id"])),
		Email:         string(fmt.Sprintf("%v", m["email"])),
		Password:      string(fmt.Sprintf("%v", m["password"])),
		ClusterID:     string(fmt.Sprintf("%v", m["cluster_id"])),
		PolicyID:      string(fmt.Sprintf("%v", m["policy_id"])),
	}
}

func convertUserToBSON(u User) bson.D {
	return bson.D{
		{"_id", u.ID},
		{"cloud_direct_id", u.CloudDirectID},
		{"email", u.Email},
		{"password", u.Password},
		{"cluster_id", u.ClusterID},
		{"policy_id", u.PolicyID},
	}
}
