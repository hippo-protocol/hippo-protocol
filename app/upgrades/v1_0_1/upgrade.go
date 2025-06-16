package v_1_0_1

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
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

		if vm[capabilitytypes.ModuleName] == 0 {
			vm[capabilitytypes.ModuleName] = 1
		}
		if err := configurator.RegisterMigration(capabilitytypes.ModuleName, 1, func(ctx sdk.Context) error {
			return nil // skip migration for capability module
		}); err != nil {
			return vm, err
		}

		vm, err := mm.RunMigrations(ctx, configurator, vm) // Run migrations for all modules
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Upgrade v1.0.1 complete")
		return vm, nil
	}
}
