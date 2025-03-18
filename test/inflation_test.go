package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/hippocrat-dao/hippo-protocol/app"
	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
)

func CalcCustomInflation(sdkCtx types.Context) math.LegacyDec {
	// In CustomInflationCalculationFn, we do not use the minter, so we can pass any value
	minter := minttypes.Minter{Inflation: math.LegacyNewDec(25).Quo(math.LegacyNewDec(100)), AnnualProvisions: math.LegacyNewDec(1)}
	params := minttypes.Params{
		BlocksPerYear: consensus.BlocksPerYear,
		InflationMin:  math.LegacyNewDec(consensus.InflationMin).Quo(math.LegacyNewDec(100)),
		InflationMax:  math.LegacyNewDec(consensus.InflationMax).Quo(math.LegacyNewDec(100)),
	}
	// In CustomInflationCalculationFn, we do not use the bondedRatio, so we can pass any value
	bondedRatio := math.LegacyNewDec(1)

	inflation := app.CustomInflationCalculationFn(sdkCtx, minter, params, bondedRatio)
	return inflation
}

func TestInflation(t *testing.T) {
	FirstYearInflatedToken := app.FirstYearInflatedToken
	GenesisSupply := app.GenesisSupply

	// Tolerance is 0.01% for errors resulting from floating point arithmetic
	tolerance, _ := math.LegacyNewDecFromStr("0.0001")

	// Test Inflation rate at the beginning of the chain
	expectedInflationRate := math.LegacyNewDec(FirstYearInflatedToken).Quo(math.LegacyNewDec(GenesisSupply)) // ~= 25%
	inflation := CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: 0}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be around 25% when block height is 0")

	// Test Inflation rate after 6 months
	expectedInflationRate = math.LegacyNewDec(FirstYearInflatedToken).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken/2)) // ~= 22%
	inflation = CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: int64(consensus.BlocksPerYear / 2)}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be calculated correctly after 6 months")

	// Test Inflation rate after 1 year
	expectedInflationRate = math.LegacyNewDec(FirstYearInflatedToken).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken)) // ~= 20%
	inflation = CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: int64(consensus.BlocksPerYear)}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be calculated correctly after 1 year")

	// Test Inflation rate after 2 years
	expectedInflationRate = math.LegacyNewDec(FirstYearInflatedToken / 2).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken*2)) // ~= 8.33%
	inflation = CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: int64(consensus.BlocksPerYear * 2)}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be calculated correctly after 2 years")

	// Test Inflation rate after 3 years
	expectedInflationRate = math.LegacyNewDec(FirstYearInflatedToken / 2).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken*2 + FirstYearInflatedToken/2)) // ~= 7.69%
	inflation = CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: int64(consensus.BlocksPerYear * 3)}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be calculated correctly after 3 years")

	// Test Inflation rate after 4 years
	expectedInflationRate = math.LegacyNewDec(FirstYearInflatedToken / 4).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken*2 + (FirstYearInflatedToken/2)*2)) // ~= 3.57%
	inflation = CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: int64(consensus.BlocksPerYear * 4)}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be calculated correctly after 4 years")

	// Test Inflation rate after 5 years
	expectedInflationRate = math.LegacyNewDec(FirstYearInflatedToken / 4).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken*2 + (FirstYearInflatedToken/2)*2 + (FirstYearInflatedToken / 4))) // ~= 3.45%
	inflation = CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: int64(consensus.BlocksPerYear * 5)}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be calculated correctly after 5 years")

	// Test Inflation rate after 6 years
	expectedInflationRate = math.LegacyNewDec(FirstYearInflatedToken / 8).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken*2 + (FirstYearInflatedToken/2)*2 + (FirstYearInflatedToken/4)*2)) // ~= 1.67%
	inflation = CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: int64(consensus.BlocksPerYear * 6)}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be calculated correctly after 6 years")

	// Test Inflation rate after 7 years
	expectedInflationRate = math.LegacyNewDec(FirstYearInflatedToken / 8).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken*2 + (FirstYearInflatedToken/2)*2 + (FirstYearInflatedToken/4)*2 + (FirstYearInflatedToken / 8))) // ~= 1.64%
	inflation = CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: int64(consensus.BlocksPerYear * 7)}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be calculated correctly after 7 years")

	// Test Inflation rate after 8 years
	expectedInflationRate = math.LegacyNewDec(FirstYearInflatedToken / 16).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken*2 + (FirstYearInflatedToken/2)*2 + (FirstYearInflatedToken/4)*2 + (FirstYearInflatedToken/8)*2)) // ~= 0.81%
	inflation = CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: int64(consensus.BlocksPerYear * 8)}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be calculated correctly after 8 years")
}
