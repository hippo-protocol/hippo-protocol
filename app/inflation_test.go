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

func CalcCustomInflation(sdkCtx types.Context) math.LegacyDec {
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
	return inflation
}

func TestInflation(t *testing.T) {
	tolerance, _ := math.LegacyNewDecFromStr("0.0001")

	// Test Inflation rate at the beginning of the chain
	expectedInflationRate := math.LegacyNewDec(FirstYearInflatedToken).Quo(math.LegacyNewDec(GenesisSupply)) // ~= 25%
	inflation := CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: 0}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be around 25% when block height is 0")

	// Test Inflation rate after 6 months
	expectedInflationRate = math.LegacyNewDec(FirstYearInflatedToken).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken/2)) // ~= 22%
	inflation = CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: int64(consensus.BlocksPerYear / 2)}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be calculated correctly after 6 months")

	// Test Inflation rate after 2 year
	expectedInflationRate = math.LegacyNewDec(FirstYearInflatedToken / 2).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken*2)) // ~= 8.3%
	inflation = CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: int64(consensus.BlocksPerYear * 2)}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be calculated correctly after 2 years")

	// Test Inflation rate after 5 years
	expectedInflationRate = math.LegacyNewDec(FirstYearInflatedToken / 4).Quo(math.LegacyNewDec(GenesisSupply + FirstYearInflatedToken*2 + (FirstYearInflatedToken/2)*2 + (FirstYearInflatedToken / 4))) // ~= 3.4%
	inflation = CalcCustomInflation(types.NewContext(nil, cmtproto.Header{Height: int64(consensus.BlocksPerYear * 5)}, false, log.NewTestLogger(t)))
	assert.True(t, inflation.GT(expectedInflationRate.Sub(tolerance)) && inflation.LT(expectedInflationRate.Add(tolerance)), "inflation rate should be calculated correctly after 5 years")
}
