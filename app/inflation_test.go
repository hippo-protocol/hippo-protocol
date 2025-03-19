package app

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/assert"

	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
)

// Helper function to calculate expected inflation
func calculateExpectedInflation(blockHeight, blocksPerYear int64) math.LegacyDec {
	targetSupply := GenesisSupply
	targetInflatedToken := FirstYearInflatedToken
	currentYear := 1 + (blockHeight / blocksPerYear)

	for i := int64(1); i <= currentYear; i++ {
		if i%2 == 1 && i != 1 {
			targetInflatedToken /= 2
		}
		targetSupply += targetInflatedToken
	}

	currentYearMinedBlock := blockHeight - ((currentYear - 1) * blocksPerYear)
	equalizer := math.LegacyOneDec().Sub(
		math.LegacyNewDec(currentYearMinedBlock - 1).QuoTruncate(math.LegacyNewDec(blocksPerYear)),
	)

	numerator := math.LegacyNewDec(targetInflatedToken)
	denominator := math.LegacyNewDec(targetSupply).Sub(
		math.LegacyNewDec(targetInflatedToken).MulTruncate(equalizer),
	)

	return numerator.QuoTruncate(denominator)
}

func TestInflation(t *testing.T) {
	testCases := []struct {
		name        string
		blockHeight int64
	}{
		{name: "genesis block", blockHeight: 0},
		{name: "6 months", blockHeight: int64(consensus.BlocksPerYear / 2)},
		{name: "1 year", blockHeight: int64(consensus.BlocksPerYear)},
		{name: "2 years", blockHeight: int64(consensus.BlocksPerYear * 2)},
		{name: "3 years", blockHeight: int64(consensus.BlocksPerYear * 3)},
		{name: "4 years", blockHeight: int64(consensus.BlocksPerYear * 4)},
		{name: "5 years", blockHeight: int64(consensus.BlocksPerYear * 5)},
		{name: "6 years", blockHeight: int64(consensus.BlocksPerYear * 6)},
		{name: "7 years", blockHeight: int64(consensus.BlocksPerYear * 7)},
		{name: "8 years", blockHeight: int64(consensus.BlocksPerYear * 8)},
	}

	params := minttypes.Params{
		BlocksPerYear: consensus.BlocksPerYear,
		InflationMin:  math.LegacyNewDec(consensus.InflationMin).Quo(math.LegacyNewDec(100)),
		InflationMax:  math.LegacyNewDec(consensus.InflationMax).Quo(math.LegacyNewDec(100)),
	}

	tolerance := math.LegacyNewDecWithPrec(1, 4) // 0.01% tolerance

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup context
			ctx := types.NewContext(nil, cmtproto.Header{Height: tc.blockHeight}, false, log.NewTestLogger(t))

			// Calculate actual and expected inflation
			actual := CustomInflationCalculationFn(
				ctx,
				minttypes.Minter{}, // Minter is not used in calculation
				params,
				math.LegacyOneDec(), // Bonded ratio not used
			)

			expected := calculateExpectedInflation(tc.blockHeight, int64(params.BlocksPerYear))

			// Apply inflation bounds
			if expected.GT(params.InflationMax) {
				expected = params.InflationMax
			}
			if expected.LT(params.InflationMin) {
				expected = params.InflationMin
			}

			// Assert with tolerance
			assert.True(t,
				actual.Sub(expected).Abs().LTE(tolerance),
				"At block height %d: Expected %s, got %s (tolerance: %s)",
				tc.blockHeight,
				expected.String(),
				actual.String(),
				tolerance.String(),
			)
		})
	}
}
