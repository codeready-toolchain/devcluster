package cluster

import (
	"context"
	"fmt"

	"github.com/alexeykazakov/devcluster/pkg/mongodb"

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

func getRequests() ([]Request, error) {
	requests := make([]Request, 0, 0)
	cursor, err := mongodb.ClusterRequests().Find(context.Background(), bson.D{})
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
	return updateRequestStatus(req.ID, "ready", "")
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
	_, err := mongodb.ClusterRequests().ReplaceOne(
		context.Background(),
		bson.D{
			{"_id", c.ID},
		},
		convertClusterToBSON(c),
		opts,
	)
	return errors.Wrap(err, "unable to replace cluster")
}

func insertCluster(c Cluster) error {
	_, err := mongodb.Clusters().InsertOne(context.Background(), convertClusterToBSON(c))
	return errors.Wrap(err, "unable to insert cluster")
}

func getCluster(id string) (*Cluster, error) {
	res := mongodb.Clusters().FindOne(
		context.Background(),
		bson.D{{"_id", id}},
	)
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

func convertBSONToRequest(m bson.M) Request {
	return Request{
		ID:          string(fmt.Sprintf("%v", m["_id"])),
		RequestedBy: string(fmt.Sprintf("%v", m["requested_by"])),
		Created:     string(fmt.Sprintf("%v", m["created"])),
		Error:       string(fmt.Sprintf("%v", m["error"])),
		Requested:   int(m["requested"].(int32)),
		Status:      string(fmt.Sprintf("%v", m["status"])),
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
