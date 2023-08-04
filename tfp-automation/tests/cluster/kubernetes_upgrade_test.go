package tests

import (
	"testing"

	cleanup "github.com/rancher/tfp-automation/functions/cleanup"
	test "github.com/rancher/tfp-automation/functions/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type KubernetesUpgradeTestSuite struct {
	suite.Suite
}

func (r *KubernetesUpgradeTestSuite) TestKubernetesUpgrade() (bool, error) {
	r.T().Parallel()

	terraformOptions, result, err := test.Setup(r.T())
	require.NoError(r.T(), err)
	assert.Equal(r.T(), true, result)

	defer cleanup.Cleanup(r.T(), terraformOptions)

	client, err := test.Provision(r.T(), terraformOptions)
	require.NoError(r.T(), err)

	err = test.KubernetesUpgrade(r.T(), terraformOptions, client)
	require.NoError(r.T(), err)

	return result, nil
}

func TestKubernetesUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(KubernetesUpgradeTestSuite))
}
