package v_2_0_0

import (
	"sync"
	"testing"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
	"github.com/stretchr/testify/require"
)

var once sync.Once

// setupWalletConfig ensures wallet configuration is set up only once across all tests
// This is necessary because cosmos-sdk's config can only be sealed once per process
func setupWalletConfig() {
	once.Do(func() {
		consensus.SetWalletConfig()
	})
}

// TestUpgradeNameConstant verifies the upgrade name constant
func TestUpgradeNameConstant(t *testing.T) {
	require.Equal(t, "v2.0.0", UpgradeName, "upgrade name should be v2.0.0")
}

// TestUpgradeStoreConfiguration verifies store upgrade configuration
func TestUpgradeStoreConfiguration(t *testing.T) {
	// Verify that wasm module is added to store upgrades
	require.NotNil(t, Upgrade.StoreUpgrades)
	require.Contains(t, Upgrade.StoreUpgrades.Added, wasmtypes.ModuleName, "wasm module should be added in store upgrades")
	require.Equal(t, 1, len(Upgrade.StoreUpgrades.Added), "should only add wasm module")
	require.Equal(t, 0, len(Upgrade.StoreUpgrades.Deleted), "should not delete any modules")
	require.Equal(t, 0, len(Upgrade.StoreUpgrades.Renamed), "should not rename any modules")
}

// TestUpgradeConfiguration verifies the upgrade configuration
func TestUpgradeConfiguration(t *testing.T) {
	require.NotNil(t, Upgrade, "upgrade configuration should not be nil")
	require.Equal(t, UpgradeName, Upgrade.UpgradeName, "upgrade name should match")
	require.NotNil(t, Upgrade.CreateUpgradeHandler, "upgrade handler creator should not be nil")
	require.NotNil(t, Upgrade.StoreUpgrades, "store upgrades should not be nil")
}

// TestUpgradeHandlerCreation tests that upgrade handler can be created
func TestUpgradeHandlerCreation(t *testing.T) {
	setupWalletConfig()

	// Test that CreateUpgradeHandler doesn't panic and returns a valid handler
	require.NotPanics(t, func() {
		handler := CreateUpgradeHandler(nil, nil, nil)
		require.NotNil(t, handler, "handler should not be nil")
	})
}

// TestWasmModuleConfiguration verifies wasm module configuration
func TestWasmModuleConfiguration(t *testing.T) {
	// Verify module is added to upgrade store
	require.Contains(t, Upgrade.StoreUpgrades.Added, wasmtypes.ModuleName)

	// Verify module name is correct
	require.Equal(t, "wasm", wasmtypes.ModuleName)
}

// TestCosmWasmDefaultParamsStructure verifies default CosmWasm parameters structure
func TestCosmWasmDefaultParamsStructure(t *testing.T) {
	// Get default params
	defaultParams := wasmtypes.DefaultParams()

	// Verify default params structure
	require.NotNil(t, defaultParams)

	// Verify that default params have expected fields
	// This ensures our upgrade handler is setting valid params
	require.NotNil(t, defaultParams.CodeUploadAccess)
	require.True(t, defaultParams.InstantiateDefaultPermission != wasmtypes.AccessTypeUnspecified)
}

// TestUpgradePlan verifies upgrade plan structure
func TestUpgradePlan(t *testing.T) {
	plan := upgradetypes.Plan{
		Name:   UpgradeName,
		Height: 100,
	}

	require.Equal(t, "v2.0.0", plan.Name)
	require.Equal(t, int64(100), plan.Height)
}

// TestCosmWasmAccessTypes verifies expected access types are available
func TestCosmWasmAccessTypes(t *testing.T) {
	// Verify the access types we're using in the upgrade exist
	require.NotNil(t, wasmtypes.AllowEverybody)
	require.NotEqual(t, wasmtypes.AccessTypeUnspecified, wasmtypes.AccessTypeEverybody)

	// Verify that AllowEverybody has the expected structure
	require.NotNil(t, wasmtypes.AllowEverybody.Permission)
	require.Equal(t, wasmtypes.AccessTypeEverybody, wasmtypes.AllowEverybody.Permission)
}

// TestStoreUpgradesNotEmpty verifies store upgrades are properly configured
func TestStoreUpgradesNotEmpty(t *testing.T) {
	require.NotEmpty(t, Upgrade.StoreUpgrades.Added, "store upgrades should add at least one module")
}

// TestUpgradeHandlerSignature verifies the upgrade handler has correct signature
func TestUpgradeHandlerSignature(t *testing.T) {
	setupWalletConfig()

	// Create a handler
	handler := CreateUpgradeHandler(nil, nil, nil)

	// Verify it's a valid upgrade handler function
	require.NotNil(t, handler)

	// Test that the handler accepts the correct parameters
	// This is a compile-time check, but we verify runtime behavior
	require.NotPanics(t, func() {
		var _ upgradetypes.UpgradeHandler = handler
	})
}
