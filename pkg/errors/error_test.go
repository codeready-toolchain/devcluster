package errors_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	devclustererr "github.com/codeready-toolchain/devcluster/pkg/errors"
	"github.com/codeready-toolchain/devcluster/test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestErrorsSuite struct {
	test.UnitTestSuite
}

func TestRunErrorsSuite(t *testing.T) {
	suite.Run(t, &TestErrorsSuite{test.UnitTestSuite{}})
}

func (s *TestErrorsSuite) TestErrors() {
	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)

	s.Run("check json error payload", func() {
		details := "testing payload"
		errMsg := "testing new error"
		code := http.StatusInternalServerError

		devclustererr.AbortWithError(ctx, code, errors.New(errMsg), details)

		res := devclustererr.Error{}
		err := json.Unmarshal(rr.Body.Bytes(), &res)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), res.Code, http.StatusInternalServerError)
		assert.Equal(s.T(), res.Details, details)
		assert.Equal(s.T(), res.Message, errMsg)
		assert.Equal(s.T(), res.Status, http.StatusText(code))
	})

	s.Run("IsNotFound", func() {
		err := devclustererr.NewNotFoundError("some message", "some details")
		assert.True(s.T(), devclustererr.IsNotFound(err))
		assert.True(s.T(), devclustererr.IsNotFound(*err))

		err = devclustererr.NewNotFoundError("some message", "some details")
		err.Code = http.StatusInternalServerError
		assert.False(s.T(), devclustererr.IsNotFound(err))
		assert.False(s.T(), devclustererr.IsNotFound(*err))

		e := errors.New("some error")
		assert.False(s.T(), devclustererr.IsNotFound(e))

		err = devclustererr.NewInternalServerError("some message", "some details")
		assert.False(s.T(), devclustererr.IsNotFound(err))
	})

	s.Run("IsInternalServerError", func() {
		err := devclustererr.NewInternalServerError("some message", "some details")
		assert.True(s.T(), devclustererr.IsInternalServerError(err))
		assert.True(s.T(), devclustererr.IsInternalServerError(*err))

		err = devclustererr.NewInternalServerError("some message", "some details")
		err.Code = http.StatusNotFound
		assert.False(s.T(), devclustererr.IsInternalServerError(err))
		assert.False(s.T(), devclustererr.IsInternalServerError(*err))

		e := errors.New("some error")
		assert.False(s.T(), devclustererr.IsInternalServerError(e))

		err = devclustererr.NewNotFoundError("some message", "some details")
		assert.False(s.T(), devclustererr.IsInternalServerError(err))
	})
}
