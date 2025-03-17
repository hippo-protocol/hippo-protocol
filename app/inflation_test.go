package app

import (
	"testing"

	"cosmossdk.io/math"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
	"github.com/stretchr/testify/assert"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/types"
)

// Test Inflation rate at the beginning of the chain
func TestInflationInitial(t *testing.T) {
	correctInflationRate := math.LegacyNewDec(FirstYearInflatedToken).Quo(math.LegacyNewDec(GenesisSupply)) // ~= 25%
	t.Log("correctInflationRate", correctInflationRate)

	sdkCtx := types.NewContext(nil, cmtproto.Header{}, false, log.NewTestLogger(t))
	// In CustomInflationCalculationFn, we do not use the minter, so we can pass any value
	minter := minttypes.Minter{Inflation: math.LegacyNewDec(10), AnnualProvisions: math.LegacyNewDec(10)}
	params := minttypes.Params{
		BlocksPerYear: consensus.BlocksPerYear,
		InflationMin:  math.LegacyNewDec(consensus.InflationMin).Quo(math.LegacyNewDec(100)),
		InflationMax:  math.LegacyNewDec(consensus.InflationMax).Quo(math.LegacyNewDec(100)),
	}
	// In CustomInflationCalculationFn, we do not use the bondedRatio, so we can pass any value
	bondedRatio := math.LegacyNewDec(1)

	inflation := CustomInflationCalculationFn(sdkCtx, minter, params, bondedRatio)
	t.Log("calculated Inflation", inflation)

	tolerance, _ := math.LegacyNewDecFromStr("0.0001")  // 0.01%
	minInflation := correctInflationRate.Sub(tolerance) // ~= 24.9%
	maxInflation := correctInflationRate.Add(tolerance) // ~= 25.1%
	assert.True(t, inflation.GT(minInflation) && inflation.LT(maxInflation), "inflation rate should be around 25% when block height is 0")
}

// Test Inflation rate after 6 months
func TestInflationAfter6Month(t *testing.T) {
	blockHeight := consensus.BlocksPerYear / 2
	// After 6 months, FirstYearInflatedToken/2 tokens are inflated
	correctInflationRate := math.LegacyNewDec(FirstYearInflatedToken).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken/2)) // ~= 22%
	t.Log("correctInflationRate", correctInflationRate)

	sdkCtx := types.NewContext(nil, cmtproto.Header{Height: int64(blockHeight)}, false, log.NewTestLogger(t))
	// In CustomInflationCalculationFn, we do not use the minter, so we can pass any value
	minter := minttypes.Minter{Inflation: math.LegacyNewDec(10), AnnualProvisions: math.LegacyNewDec(10)}
	params := minttypes.Params{
		BlocksPerYear: consensus.BlocksPerYear,
		InflationMin:  math.LegacyNewDec(consensus.InflationMin).Quo(math.LegacyNewDec(100)),
		InflationMax:  math.LegacyNewDec(consensus.InflationMax).Quo(math.LegacyNewDec(100)),
	}
	// In CustomInflationCalculationFn, we do not use the bondedRatio, so we can pass any value
	bondedRatio := math.LegacyNewDec(1)

	inflation := CustomInflationCalculationFn(sdkCtx, minter, params, bondedRatio)
	t.Log("calculated Inflation", inflation)

	tolerance, _ := math.LegacyNewDecFromStr("0.0001") // 0.01%
	minInflation := correctInflationRate.Sub(tolerance)
	maxInflation := correctInflationRate.Add(tolerance)
	assert.True(t, inflation.GT(minInflation) && inflation.LT(maxInflation), "inflation rate should be calculated correctly after 6 months")
}

// Test Inflation rate after 2 year
func TestInflationAfter2Year(t *testing.T) {
	blockHeight := consensus.BlocksPerYear * 2
	// After 2 years, FirstYearInflatedToken * 2 tokens are inflated
	// And the targetInflatedToken is halved every 2 year, now it is halved 1 time
	correctInflationRate := math.LegacyNewDec(FirstYearInflatedToken / 2).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken*2)) // ~= 8.3%
	t.Log("correctInflationRate", correctInflationRate)

	sdkCtx := types.NewContext(nil, cmtproto.Header{Height: int64(blockHeight)}, false, log.NewTestLogger(t))
	// In CustomInflationCalculationFn, we do not use the minter, so we can pass any value
	minter := minttypes.Minter{Inflation: math.LegacyNewDec(10), AnnualProvisions: math.LegacyNewDec(10)}
	params := minttypes.Params{
		BlocksPerYear: consensus.BlocksPerYear,
		InflationMin:  math.LegacyNewDec(consensus.InflationMin).Quo(math.LegacyNewDec(100)),
		InflationMax:  math.LegacyNewDec(consensus.InflationMax).Quo(math.LegacyNewDec(100)),
	}
	// In CustomInflationCalculationFn, we do not use the bondedRatio, so we can pass any value
	bondedRatio := math.LegacyNewDec(1)

	inflation := CustomInflationCalculationFn(sdkCtx, minter, params, bondedRatio)
	t.Log("calculated Inflation", inflation)

	tolerance, _ := math.LegacyNewDecFromStr("0.0001") // 0.01%
	minInflation := correctInflationRate.Sub(tolerance)
	maxInflation := correctInflationRate.Add(tolerance)
	assert.True(t, inflation.GT(minInflation) && inflation.LT(maxInflation), "inflation rate should be calculated correctly after 2 years")

}

// Test Inflation rate after 5 years
func TestInflationAfter5Year(t *testing.T) {
	blockHeight := consensus.BlocksPerYear * 5
	// After 5 years, FirstYearInflatedToken * 5 tokens are inflated
	// And the targetInflatedToken is halved every 2 year, now it is halved 2 times

	correctInflationRate := math.LegacyNewDec(FirstYearInflatedToken / 4).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken*2 + (FirstYearInflatedToken/2)*2 + (FirstYearInflatedToken / 4))) // ~= 3.4%
	t.Log("correctInflationRate", correctInflationRate)

	sdkCtx := types.NewContext(nil, cmtproto.Header{Height: int64(blockHeight)}, false, log.NewTestLogger(t))
	// In CustomInflationCalculationFn, we do not use the minter, so we can pass any value
	minter := minttypes.Minter{Inflation: math.LegacyNewDec(10), AnnualProvisions: math.LegacyNewDec(10)}
	params := minttypes.Params{
		BlocksPerYear: consensus.BlocksPerYear,
		InflationMin:  math.LegacyNewDec(consensus.InflationMin).Quo(math.LegacyNewDec(100)),
		InflationMax:  math.LegacyNewDec(consensus.InflationMax).Quo(math.LegacyNewDec(100)),
	}
	// In CustomInflationCalculationFn, we do not use the bondedRatio, so we can pass any value
	bondedRatio := math.LegacyNewDec(1)

	inflation := CustomInflationCalculationFn(sdkCtx, minter, params, bondedRatio)
	t.Log("calculated Inflation", inflation)

	tolerance, _ := math.LegacyNewDecFromStr("0.0001") // 0.01%
	minInflation := correctInflationRate.Sub(tolerance)
	maxInflation := correctInflationRate.Add(tolerance)
	assert.True(t, inflation.GT(minInflation) && inflation.LT(maxInflation), "inflation rate should be calculated correctly after 5 years")
}
