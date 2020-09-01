package configuration_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/alexeykazakov/devcluster/pkg/configuration"
	"github.com/alexeykazakov/devcluster/test"
	. "github.com/codeready-toolchain/toolchain-common/pkg/test"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestConfigurationSuite struct {
	test.UnitTestSuite
}

func TestRunConfigurationSuite(t *testing.T) {
	suite.Run(t, &TestConfigurationSuite{test.UnitTestSuite{}})
}

// getDefaultConfiguration returns a configuration registry without anything but
// defaults set. Remember that environment variables can overwrite defaults, so
// please ensure to properly unset envionment variables using
// UnsetEnvVarAndRestore().
func (s *TestConfigurationSuite) getDefaultConfiguration() *configuration.Config {
	config := configuration.New()
	return config
}

func (s *TestConfigurationSuite) TestNew() {
	s.Run("default configuration", func() {
		reg := configuration.New()
		require.NotNil(s.T(), reg)
	})
}

func (s *TestConfigurationSuite) TestGetHTTPAddress() {
	key := configuration.EnvPrefix + "_" + "HTTP_ADDRESS"

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), configuration.DefaultHTTPAddress, config.GetHTTPAddress())
	})

	s.Run("env overwrite", func() {
		u, err := uuid.NewV4()
		require.NoError(s.T(), err)
		newVal := u.String()
		err = os.Setenv(key, newVal)
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), newVal, config.GetHTTPAddress())
	})
}

func (s *TestConfigurationSuite) TestGetLogLevel() {
	key := configuration.EnvPrefix + "_" + "LOG_LEVEL"
	resetFunc := UnsetEnvVarAndRestore(s.T(), key)
	defer resetFunc()

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), configuration.DefaultLogLevel, config.GetLogLevel())
	})

	s.Run("env overwrite", func() {
		u, err := uuid.NewV4()
		require.NoError(s.T(), err)
		newVal := u.String()
		err = os.Setenv(key, newVal)
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), newVal, config.GetLogLevel())
	})
}

func (s *TestConfigurationSuite) TestIsLogJSON() {
	key := configuration.EnvPrefix + "_" + "LOG_JSON"
	resetFunc := UnsetEnvVarAndRestore(s.T(), key)
	defer resetFunc()

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), configuration.DefaultLogJSON, config.IsLogJSON())
	})

	s.Run("env overwrite", func() {
		newVal := !configuration.DefaultLogJSON
		err := os.Setenv(key, strconv.FormatBool(newVal))
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), newVal, config.IsLogJSON())
	})
}

func (s *TestConfigurationSuite) TestGetGracefulTimeout() {
	key := configuration.EnvPrefix + "_" + "GRACEFUL_TIMEOUT"

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), configuration.DefaultGracefulTimeout, config.GetGracefulTimeout())
	})

	s.Run("env overwrite", func() {
		newVal := 666 * time.Second
		err := os.Setenv(key, newVal.String())
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), newVal, config.GetGracefulTimeout())
	})
}

func (s *TestConfigurationSuite) TestGetHTTPWriteTimeout() {
	key := configuration.EnvPrefix + "_" + "HTTP_WRITE_TIMEOUT"

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), configuration.DefaultHTTPWriteTimeout, config.GetHTTPWriteTimeout())
	})

	s.Run("env overwrite", func() {
		newVal := 666 * time.Second
		err := os.Setenv(key, newVal.String())
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), newVal, config.GetHTTPWriteTimeout())
	})
}

func (s *TestConfigurationSuite) TestGetHTTPReadTimeout() {
	key := configuration.EnvPrefix + "_" + "HTTP_READ_TIMEOUT"

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), configuration.DefaultHTTPReadTimeout, config.GetHTTPReadTimeout())
	})

	s.Run("env overwrite", func() {
		newVal := 777 * time.Second
		err := os.Setenv(key, newVal.String())
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), newVal, config.GetHTTPReadTimeout())
	})
}

func (s *TestConfigurationSuite) TestGetHTTPIdleTimeout() {
	key := configuration.EnvPrefix + "_" + "HTTP_IDLE_TIMEOUT"

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), configuration.DefaultHTTPIdleTimeout, config.GetHTTPIdleTimeout())
	})

	s.Run("env overwrite", func() {
		newVal := 888 * time.Second
		err := os.Setenv(key, newVal.String())
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), newVal, config.GetHTTPIdleTimeout())
	})
}

func (s *TestConfigurationSuite) TestGetHTTPCompressResponses() {
	key := configuration.EnvPrefix + "_" + "HTTP_COMPRESS"

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), configuration.DefaultHTTPCompressResponses, config.GetHTTPCompressResponses())
	})

	s.Run("env overwrite", func() {
		newVal := !configuration.DefaultHTTPCompressResponses
		err := os.Setenv(key, strconv.FormatBool(newVal))
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), newVal, config.GetHTTPCompressResponses())
	})
}

func (s *TestConfigurationSuite) TestGetEnvironmentAndTestingMode() {
	key := fmt.Sprintf("%s_ENVIRONMENT", configuration.EnvPrefix)
	resetFunc := UnsetEnvVarAndRestore(s.T(), key)
	defer resetFunc()

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), "prod", config.GetEnvironment())
		assert.False(s.T(), config.IsTestingMode())
	})

	s.Run("env overwrite", func() {
		err := os.Setenv(key, "TestGetEnvironmentFromEnvVar")
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), "TestGetEnvironmentFromEnvVar", config.GetEnvironment())
		assert.False(s.T(), config.IsTestingMode())
	})

	s.Run("unit-tests env", func() {
		err := os.Setenv(key, "unit-tests")
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.True(s.T(), config.IsTestingMode())
	})
}

func (s *TestConfigurationSuite) TestGetAuthClientConfigRaw() {
	key := configuration.EnvPrefix + "_" + "AUTH_CLIENT_CONFIG_RAW"

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), configuration.DefaultAuthClientConfigRaw, config.GetAuthClientConfigAuthRaw())
	})

	s.Run("env overwrite", func() {
		u, err := uuid.NewV4()
		require.NoError(s.T(), err)
		newVal := u.String()
		err = os.Setenv(key, newVal)
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), newVal, config.GetAuthClientConfigAuthRaw())
	})
}

func (s *TestConfigurationSuite) TestGetAuthClientConfigContentType() {
	key := configuration.EnvPrefix + "_" + "AUTH_CLIENT_CONFIG_CONTENT_TYPE"

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), configuration.DefaultAuthClientConfigContentType, config.GetAuthClientConfigAuthContentType())
	})

	s.Run("env overwrite", func() {
		u, err := uuid.NewV4()
		require.NoError(s.T(), err)
		newVal := u.String()
		err = os.Setenv(key, newVal)
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), newVal, config.GetAuthClientConfigAuthContentType())
	})
}

func (s *TestConfigurationSuite) TestGetAuthClientLibraryURL() {
	key := configuration.EnvPrefix + "_" + "AUTH_CLIENT_LIBRARY_URL"
	resetFunc := UnsetEnvVarAndRestore(s.T(), key)
	defer resetFunc()

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), configuration.DefaultAuthClientLibraryURL, config.GetAuthClientLibraryURL())
	})

	s.Run("env overwrite", func() {
		u, err := uuid.NewV4()
		require.NoError(s.T(), err)
		newVal := u.String()
		err = os.Setenv(key, newVal)
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), newVal, config.GetAuthClientLibraryURL())
	})
}

func (s *TestConfigurationSuite) TestGetAuthClientPublicKeysURL() {
	key := configuration.EnvPrefix + "_" + "AUTH_CLIENT_PUBLIC_KEYS_URL"
	resetFunc := UnsetEnvVarAndRestore(s.T(), key)
	defer resetFunc()

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), configuration.DefaultAuthClientPublicKeysURL, config.GetAuthClientPublicKeysURL())
	})

	s.Run("env overwrite", func() {
		u, err := uuid.NewV4()
		require.NoError(s.T(), err)
		newVal := u.String()
		err = os.Setenv(key, newVal)
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), newVal, config.GetAuthClientPublicKeysURL())
	})
}

func (s *TestConfigurationSuite) TestGetNamespace() {
	key := configuration.EnvPrefix + "_" + "NAMESPACE"
	resetFunc := UnsetEnvVarAndRestore(s.T(), key)
	defer resetFunc()

	s.Run("default", func() {
		resetFunc := UnsetEnvVarAndRestore(s.T(), key)
		defer resetFunc()
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), configuration.DefaultNamespace, config.GetNamespace())
	})

	s.Run("env overwrite", func() {
		u, err := uuid.NewV4()
		require.NoError(s.T(), err)
		newVal := u.String()
		err = os.Setenv(key, newVal)
		require.NoError(s.T(), err)
		config := s.getDefaultConfiguration()
		assert.Equal(s.T(), newVal, config.GetNamespace())
	})
}
