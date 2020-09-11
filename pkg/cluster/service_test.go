package cluster

import (
	"testing"
	"time"

	"github.com/codeready-toolchain/devcluster/pkg/log"

	"github.com/stretchr/testify/assert"
)

func TestExpired(t *testing.T) {
	log.Init("devcluster-testing")

	t.Run("expired", func(t *testing.T) {
		r := Request{
			Created:       time.Now().Add(-2 * time.Hour).Unix(),
			DeleteInHours: 1,
		}
		assert.True(t, expired(r))
	})

	t.Run("not expired", func(t *testing.T) {
		r := Request{
			Created:       time.Now().Add(-2 * time.Hour).Unix(),
			DeleteInHours: 3,
		}
		assert.False(t, expired(r))
	})
}
