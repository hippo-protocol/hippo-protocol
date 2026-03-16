package v_2_0_0

import (
	storetypes "cosmossdk.io/store/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/hippocrat-dao/hippo-protocol/app/upgrades"
)

const (
	UpgradeName = "v2.0.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{wasmtypes.ModuleName},
	},
}
