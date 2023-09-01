package tests

import (
	"testing"

	cleanup "github.com/rancher/tfp-automation/functions/cleanup"
	test "github.com/rancher/tfp-automation/functions/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ProvisionTestSuite struct {
	suite.Suite
}

func (r *ProvisionTestSuite) TestProvision() (bool, error) {
	r.T().Parallel()

	terraformOptions, result, err := test.Setup(r.T())
	require.NoError(r.T(), err)
	assert.Equal(r.T(), true, result)

	defer cleanup.Cleanup(r.T(), terraformOptions)

	_, error1 := test.Provision(r.T(), terraformOptions)
	require.NoError(r.T(), error1)
	assert.Equal(r.T(), true, result)

	return result, nil
}

func TestProvisionTestSuite(t *testing.T) {
	suite.Run(t, new(ProvisionTestSuite))
}
