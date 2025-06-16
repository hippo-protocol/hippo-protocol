package v_1_0_1

import (
	"github.com/hippocrat-dao/hippo-protocol/app/upgrades"
)

const (
	UpgradeName = "v1.0.1"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
}
