package controller

import (
	"testing"

	"github.com/codeready-toolchain/devcluster/pkg/ibmcloud"
	"github.com/codeready-toolchain/devcluster/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestClusterReqSuite struct {
	test.UnitTestSuite
}

func TestRunClusterReqSuite(t *testing.T) {
	suite.Run(t, &TestClusterReqSuite{test.UnitTestSuite{}})
}

func (s *TestClusterReqSuite) TestFilterZones() {
	r := &ClusterRequest{}
	expectedZones := []ibmcloud.Location{
		dc("wdc04"), dc("wdc06"), dc("wdc07"),
		dc("che01"),
		dc("fra02"), dc("fra04"), dc("fra05"),
		dc("ams03"),
	}
	extraZones := []ibmcloud.Location{
		dc("syd01"), dc("mil01"),
	}
	result := r.filterZones(nil, append(extraZones, expectedZones...))
	assert.Equal(s.T(), expectedZones, result)
}

func dc(name string) ibmcloud.Location {
	return ibmcloud.Location{
		ID:          name,
		Name:        name,
		Kind:        "dc",
		DisplayName: "DC-" + name,
	}
}
