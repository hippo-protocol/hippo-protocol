package v_1_0_3

import (
	storetypes "cosmossdk.io/store/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/hippocrat-dao/hippo-protocol/app/upgrades"
)

const (
	UpgradeName = "v1.0.3"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{wasmtypes.ModuleName},
	},
}
