package test

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"hash/fnv"
	"os"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/codeready-toolchain/devcluster/pkg/configuration"
	"github.com/codeready-toolchain/devcluster/pkg/log"
	"github.com/codeready-toolchain/devcluster/pkg/mongodb"

	"github.com/stretchr/testify/suite"
)

// IntegrationTestSuite is the base test suite for integration tests.
type IntegrationTestSuite struct {
	suite.Suite
	Config          *configuration.Config
	mongoDisconnect func()
}

// SetupSuite sets the suite up and sets testmode.
func (s *IntegrationTestSuite) SetupSuite() {
	// create logger and registry
	log.Init("devcluster-testing")

	s.Config = configuration.New()

	// set environment to integration-tests
	s.Config.GetViperInstance().Set("environment", configuration.IntegrationTestsEnvironment)

	// Init the default mongo connection.
	// DEVCLUSTER_MONGODB_CONNECTION_STRING env var is expected to be set to a test MongoDB
	cs := os.Getenv("DEVCLUSTER_MONGODB_CONNECTION_STRING")
	if cs == "" {
		panic("DEVCLUSTER_MONGODB_CONNECTION_STRING env var is not set. It must be set to a test MongoDB")
	}

	// Random db name
	dbname:=fmt.Sprintf("devclustertest-%d", hash(uuid.NewV4().String()))
	s.Config.GetViperInstance().Set("mongodb.database", dbname)
	disconnect, err := mongodb.InitDefaultClient(s.Config)
	if err != nil {
		fmt.Printf("Connection string:" + s.Config.GetMongodbConnectionString())
		panic(err)
	}
	s.mongoDisconnect = disconnect
	s.cleanupDatabase()
}

func hash(s string) uint32 {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		log.Error(nil, err, "unable to generate cluster name")
	}
	return h.Sum32()
}

// TearDownSuite tears down the test suite.
func (s *IntegrationTestSuite) TearDownSuite() {
	s.cleanupDatabase()
	s.mongoDisconnect()
}

func (s *IntegrationTestSuite) cleanupDatabase() {
	names, err := mongodb.Devcluster().ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		panic(err)
	}
	for _, name := range names {
		err = mongodb.Devcluster().Collection(name).Drop(context.Background())
		if err != nil {
			panic(err)
		}
	}
}
