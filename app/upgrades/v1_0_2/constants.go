package v_1_0_2

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/hippocrat-dao/hippo-protocol/app/upgrades"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

const (
	UpgradeName = "v1.0.2"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: storetypes.StoreUpgrades{
		Added: []string{wasmtypes.ModuleName},
	},
}
