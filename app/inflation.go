package app

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// InflationCalculationFn defines the function required to calculate inflation rate during
// BeginBlock. It receives the minter and params stored in the keeper, along with the current
// bondedRatio and returns the newly calculated inflation rate.
// It can be used to specify a custom inflation calculation logic, instead of relying on the
// default logic provided by the sdk.
func CustomInflationCalculationFn(ctx sdk.Context, minter minttypes.Minter, params minttypes.Params, bondedRatio sdk.Dec) sdk.Dec {
	//	targetSupply <- genesisSupply
	//	targetInflatedToken <- firstYearInflatedToken
	//	currentYear <- 1 + floor(currentBlockHeight / BlocksPerYear)
	//
	//	for i <-1 to currentYear do
	//		if i % 2 = 1 and i != 1 then
	//			targetInflatedToken <- targetInflatedToken / 2
	//		end
	//		targetSupply <- targetSupply + targetInflatedToken
	//	end
	//
	//	inflation <- targetInflatedToken / (targetSupply - targetInflatedToken)

	genesisSupply := int64(1_084_734_273)
	firstYearInflatedToken := int64(271_183_568)

	targetSupply := genesisSupply
	targetInflatedToken := firstYearInflatedToken
	currentYear := 1 + (ctx.BlockHeight() / int64(params.BlocksPerYear))

	for i := int64(1); i <= currentYear; i++ {
		if i%2 == 1 && i != 1 {
			targetInflatedToken /= 2
		}
		targetSupply += targetInflatedToken
	}

	inflation := math.LegacyNewDec(targetInflatedToken).Quo(math.LegacyNewDec(targetSupply).Sub(math.LegacyNewDec(targetInflatedToken)))

	if inflation.GT(params.InflationMax) {
		inflation = params.InflationMax
	}
	if inflation.LT(params.InflationMin) {
		inflation = params.InflationMin
	}

	ctx.Logger().Info("INFLATION:::" + inflation.String())

	return inflation
}
