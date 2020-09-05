package test

import (
	"context"
	"testing"

	"github.com/alexeykazakov/devcluster/pkg/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestIntegrationSuite struct {
	IntegrationTestSuite
}

func TestRunDTestIntegrationSuite(t *testing.T) {
	suite.Run(t, &TestIntegrationSuite{IntegrationTestSuite{}})
}

func (s *TestIntegrationSuite) TestFoo() {

	opts := options.Replace().SetUpsert(true)
	_, err := mongodb.Clusters().ReplaceOne(
		context.Background(),
		bson.D{
			{"_id", "123"},
		},
		bson.D{
			{"_id", "123"},
			{"name", "test"},
		},
		opts,
	)
	assert.NoError(s.T(), err)

	s.Run("get before init", func() {
		mongodb.Clusters()
		require.True(s.T(), true)
		assert.True(s.T(), true)
	})
}
