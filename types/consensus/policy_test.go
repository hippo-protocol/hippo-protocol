package consensus_test

import (
	"testing"
	"time"

	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
	"github.com/stretchr/testify/require"
)

func TestConsensusPolicyConstants(t *testing.T) {
	// Check that blocks have expected sizes
	require.Equal(t, 4194304, consensus.MaxBlockSize, "MaxBlockSize should be 4MB")
	require.Equal(t, 100000000, consensus.MaxBlockGas, "MaxBlockGas should be 100 million")

	// Check minimum gas price
	expectedMinGasPrices := "5000000000000" + consensus.DefaultHippoDenom
	require.Equal(t, expectedMinGasPrices, consensus.MinGasPrices, "MinGasPrices calculated incorrectly")

	// Check block time parameters
	require.Equal(t, 6, consensus.BlockTimeSec, "BlockTimeSec should be 6 seconds")
	require.Equal(t, 60*60*24*7*3*time.Second, consensus.UnbondingPeriod, "UnbondingPeriod calculated incorrectly")

	// Check staking parameters
	require.Equal(t, 22, consensus.MaxValidators, "MaxValidators should be 22")
	require.Equal(t, 10, consensus.MinCommissionRate, "MinCommissionRate should be 10")

	// Check inflation parameters
	require.Equal(t, 25, consensus.Minter, "Minter should be 25")
	require.Equal(t, 25, consensus.InflationRateChange, "InflationRateChange should be 25")
	require.Equal(t, 0, consensus.InflationMin, "InflationMin should be 0")
	require.Equal(t, 25, consensus.InflationMax, "InflationMax should be 25")
	require.Equal(t, uint64(60*60*24*365)/uint64(consensus.BlockTimeSec), consensus.BlocksPerYear, "BlocksPerYear calculated incorrectly")

	// Check distribution parameters
	require.Equal(t, 92, consensus.CommunityTax, "CommunityTax should be 92")

	// Check governance parameters
	require.Equal(t, 50000, consensus.MinDepositTokens, "MinDepositTokens should be 50,000")
	require.Equal(t, 60*60*24*14*time.Second, consensus.MaxDepositPeriod, "MaxDepositPeriod calculated incorrectly")
	require.Equal(t, 60*60*24*14*time.Second, consensus.VotingPeriod, "VotingPeriod calculated incorrectly")

	// Check slashing parameters
	require.Equal(t, 10000, consensus.SignedBlocksWindow, "SignedBlocksWindow should be 10,000")
	require.Equal(t, 75, consensus.MinSignedPerWindow, "MinSignedPerWindow should be 75")
	require.Equal(t, 5, consensus.SlashFractionDoubleSign, "SlashFractionDoubleSign should be 5")
	require.Equal(t, 0, consensus.SlashFractionDowntime, "SlashFractionDowntime should be 0")

	// Check evidence parameters
	expectedMaxAgeDuration := consensus.UnbondingPeriod * 30 / 21
	require.Equal(t, expectedMaxAgeDuration, consensus.MaxAgeDuration, "MaxAgeDuration calculated incorrectly")

	expectedMaxAgeNumBlocks := consensus.BlocksPerYear * 30 / 365
	require.Equal(t, expectedMaxAgeNumBlocks, consensus.MaxAgeNumBlocks, "MaxAgeNumBlocks calculated incorrectly")
}
