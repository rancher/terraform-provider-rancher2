package tests

import (
	"testing"

	cleanup "github.com/rancher/tfp-automation/functions/cleanup"
	test "github.com/rancher/tfp-automation/functions/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BuildModuleTestSuite struct {
	suite.Suite
}

func (r *BuildModuleTestSuite) TestBuildModule() (bool, error) {
	r.T().Parallel()

	defer cleanup.CleanupConfigTF(r.T())

	result, err := test.BuildModule(r.T())
	require.NoError(r.T(), err)
	assert.Equal(r.T(), true, result)

	return result, nil
}

func TestBuildModuleTestSuite(t *testing.T) {
	suite.Run(t, new(BuildModuleTestSuite))
}
