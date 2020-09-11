package cluster

import (
	"testing"
	"time"

	"github.com/codeready-toolchain/devcluster/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestServiceSuite struct {
	test.UnitTestSuite
}

func TestRunServiceSuite(t *testing.T) {
	suite.Run(t, &TestServiceSuite{test.UnitTestSuite{}})
}

func (s *TestServiceSuite) TestExpired() {
	s.T().Run("expired", func(t *testing.T) {
		r := Request{
			Created:       time.Now().Add(-2 * time.Hour).Unix(),
			DeleteInHours: 1,
		}
		assert.True(t, expired(r))
	})

	s.T().Run("not expired", func(t *testing.T) {
		r := Request{
			Created:       time.Now().Add(-2 * time.Hour).Unix(),
			DeleteInHours: 3,
		}
		assert.False(t, expired(r))
	})
}
