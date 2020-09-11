package test

import (
	"github.com/codeready-toolchain/devcluster/pkg/configuration"
	"github.com/codeready-toolchain/devcluster/pkg/log"
	"github.com/codeready-toolchain/devcluster/test/resource"

	"github.com/stretchr/testify/suite"
)

// UnitTestSuite is the base test suite for unit tests.
type UnitTestSuite struct {
	suite.Suite
	Config *configuration.Config
}

// SetupSuite sets the suite up and sets testmode.
func (s *UnitTestSuite) SetupSuite() {
	resource.Require(s.T(), resource.UnitTest)

	// create logger and registry
	log.Init("devcluster-testing")

	s.Config = configuration.New()

	// set environment to unit-tests
	s.Config.GetViperInstance().Set("environment", configuration.UnitTestsEnvironment)
}

// TearDownSuite tears down the test suite.
func (s *UnitTestSuite) TearDownSuite() {
	// summon the GC!
	s.Config = nil
}
