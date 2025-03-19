package consensus_test

import (
	"testing"

	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
	"github.com/stretchr/testify/require"
)

func TestHippoDenoms(t *testing.T) {
	// Check that the denomination symbols are set correctly
	require.Equal(t, "ahp", consensus.AlphaHippoDenom, "AlphaHippoDenom should be 'ahp'")
	require.Equal(t, "uhp", consensus.MicroHippoDenom, "MicroHippoDenom should be 'uhp'")
	require.Equal(t, "mhp", consensus.MilliHippoDenom, "MilliHippoDenom should be 'mhp'")
	require.Equal(t, "chp", consensus.CentiHippoDenom, "CentiHippoDenom should be 'chp'")

	// Check the precision for each denomination
	require.Equal(t, int64(18), consensus.AlphaHippoPrecision, "AlphaHippoPrecision should be 18")
	require.Equal(t, int64(6), consensus.MicroHippoPrecision, "MicroHippoPrecision should be 6")
	require.Equal(t, int64(3), consensus.MilliHippoPrecision, "MilliHippoPrecision should be 3")
	require.Equal(t, int64(2), consensus.CentiHippoPrecision, "CentiHippoPrecision should be 2")

	// Check the default values
	require.Equal(t, consensus.AlphaHippoDenom, consensus.DefaultHippoDenom, "DefaultHippoDenom should be equal to AlphaHippoDenom")
	require.Equal(t, consensus.AlphaHippoPrecision, consensus.DefaultHippoPrecision, "DefaultHippoPrecision should be equal to AlphaHippoPrecision")
}
