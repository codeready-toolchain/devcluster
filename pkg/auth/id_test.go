package auth_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/codeready-toolchain/devcluster/pkg/auth"
	"github.com/codeready-toolchain/devcluster/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestIDSuite struct {
	test.UnitTestSuite
}

func TestRunIDSuite(t *testing.T) {
	suite.Run(t, &TestIDSuite{test.UnitTestSuite{}})
}

func (s *TestIDSuite) TestKeyManagerDefaultKeyManager() {
	s.Run("generate ID with prefix", func() {
		id := auth.GenerateShortID("rhd")
		assert.True(s.T(), strings.HasPrefix(id, "rhd-"))
		assert.True(s.T(), len(id) > 10)
	})
	s.Run("generate ID no prefix", func() {
		id := auth.GenerateShortID("")
		assert.False(s.T(), strings.HasPrefix(id, "-"))
	})
	s.Run("generate ID with date", func() {
		prefix := fmt.Sprintf("rhd-%s-", time.Now().Format("Jan02"))
		id := auth.GenerateShortIDWithDate("rhd")
		assert.True(s.T(), strings.HasPrefix(id, prefix), "expected format: %s, actual ID: %s", prefix+"xxx", id)
	})
}
