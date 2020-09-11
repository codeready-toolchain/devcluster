package static_test

import (
	"testing"

	"github.com/codeready-toolchain/devcluster/test/resource"

	"github.com/codeready-toolchain/devcluster/pkg/static"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatic(t *testing.T) {
	resource.Require(t, resource.UnitTest)

	// Get the static assets.
	hfs := static.Assets
	// Open the default file; note that the
	// actual files and contents are tested elsewhere.
	file, err := hfs.Open("index.html")
	require.NoError(t, err)
	// Check the file stats.
	stat, err := file.Stat()
	require.NoError(t, err)
	assert.Greater(t, stat.Size(), int64(0), "static asset 'index.html' size is zero.")
}
