package v_1_0_2

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/hippocrat-dao/hippo-protocol/app/keepers"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	keepers *keepers.AppKeepersWithKey,
) upgradetypes.UpgradeHandler {
	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)
		ctx.Logger().Info("Starting module migrations...")

		vm, err := mm.RunMigrations(ctx, configurator, vm) // Run migrations for all modules
		if err != nil {
			return vm, err
		}
		metadata := banktypes.Metadata{
			Description: "The native staking token of the Hippo Protocol.",
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    "hp",
					Exponent: 18,
					Aliases:  []string{},
				},
				{
					Denom:    "ahp",
					Exponent: 0,
					Aliases:  []string{},
				},
			},
			Base:    "ahp",
			Display: "hp",
			Name:    "Hippo",
			Symbol:  "HP",
			URI:     "",
			URIHash: "",
		}
		keepers.BankKeeper.SetDenomMetaData(ctx, metadata)

		ctx.Logger().Info("Upgrade v1.0.2 complete")
		return vm, nil
	}
}
