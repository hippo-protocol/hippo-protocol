package v_1_0_2

import (
	"github.com/hippocrat-dao/hippo-protocol/app/upgrades"
)

const (
	UpgradeName = "v1.0.2"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
}
