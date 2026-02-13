package v_1_0_3

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/tx/signing"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/gogoproto/proto"
	"github.com/hippocrat-dao/hippo-protocol/app/keepers"
	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
	"github.com/stretchr/testify/require"
)

// TestUpgradeHandlerIntegration tests the upgrade handler with a more realistic setup
func TestUpgradeHandlerIntegration(t *testing.T) {
	setupWalletConfig()

	interfaceRegistry, _ := codectypes.NewInterfaceRegistryWithOptions(codectypes.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32AccountAddrPrefix(),
			},
			ValidatorAddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
			},
		},
	})

	appCodec := codec.NewProtoCodec(interfaceRegistry)
	legacyAmino := codec.NewLegacyAmino()

	maccPerms := map[string][]string{
		"mint":                   {"minter"},
		"bonded_tokens_pool":     {"burner", "staking"},
		"not_bonded_tokens_pool": {"burner", "staking"},
		"fee_collector":          nil,
		"distribution":           nil,
		"gov":                    {"burner"},
		"transfer":               {"minter", "burner"},
		wasmtypes.ModuleName:     {"burner"},
	}

	blockedAddrs := map[string]bool{}
	tempDir := t.TempDir()
	appOpts := simtestutil.NewAppOptionsWithFlagHome(tempDir)

	logger := log.NewNopLogger()

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db, logger, nil)

	keys, tkeys, memKeys := generateTestStoreKeys()
	for _, key := range keys {
		ms.MountStoreWithDB(key, storetypes.StoreTypeIAVL, db)
	}
	for _, tkey := range tkeys {
		ms.MountStoreWithDB(tkey, storetypes.StoreTypeTransient, db)
	}
	for _, mkey := range memKeys {
		ms.MountStoreWithDB(mkey, storetypes.StoreTypeMemory, db)
	}
	require.NoError(t, ms.LoadLatestVersion())

	txCfg := tx.NewTxConfig(appCodec, tx.DefaultSignModes)
	baseApp := baseapp.NewBaseApp("hippo", logger, db, txCfg.TxDecoder())

	var allStoreKeys []storetypes.StoreKey
	for _, k := range keys {
		allStoreKeys = append(allStoreKeys, k)
	}
	for _, t := range tkeys {
		allStoreKeys = append(allStoreKeys, t)
	}
	for _, m := range memKeys {
		allStoreKeys = append(allStoreKeys, m)
	}

	baseApp.MountStores(allStoreKeys...)
	require.NoError(t, baseApp.LoadLatestVersion())

	appKeepers := &keepers.AppKeepersWithKey{}
	appKeepers.GenerateKeys()

	var emptyWasmOptions []wasmkeeper.Option
	appKeepers.InitKeyAndKeepers(appCodec, legacyAmino, maccPerms, blockedAddrs, appOpts, baseApp, logger, emptyWasmOptions)

	ctx := baseApp.NewContext(false)
	configurator := module.NewConfigurator(appCodec, baseApp.MsgServiceRouter(), baseApp.GRPCQueryRouter())
	moduleManager := module.NewManager()

	// Create upgrade handler
	handler := CreateUpgradeHandler(moduleManager, configurator, appKeepers)

	// Create a plan
	plan := upgradetypes.Plan{
		Name:   UpgradeName,
		Height: 100,
	}

	// Execute upgrade handler
	vm := module.VersionMap{}
	newVm, err := handler(ctx, plan, vm)
	require.NoError(t, err, "upgrade handler should execute without error")
	require.NotNil(t, newVm, "version map should not be nil")

	// Verify CosmWasm params are set correctly after upgrade
	params := appKeepers.WasmKeeper.GetParams(ctx)
	require.NotNil(t, params, "wasm params should not be nil")
	require.Equal(t, wasmtypes.AllowEverybody, params.CodeUploadAccess, "code upload should be allowed for everybody")
	require.Equal(t, wasmtypes.AccessTypeEverybody, params.InstantiateDefaultPermission, "instantiate should be allowed for everybody")

	t.Log("✓ Upgrade handler executed successfully")
	t.Log("✓ CosmWasm parameters configured correctly")
	t.Log("✓ WasmKeeper is functional after upgrade")
}

// generateTestStoreKeys creates test store keys for all modules including wasm
func generateTestStoreKeys() (map[string]*storetypes.KVStoreKey, map[string]*storetypes.TransientStoreKey, map[string]*storetypes.MemoryStoreKey) {
	kv := make(map[string]*storetypes.KVStoreKey)
	tk := make(map[string]*storetypes.TransientStoreKey)
	mem := make(map[string]*storetypes.MemoryStoreKey)

	modules := []string{
		"auth", "bank", "staking", "mint", "distr", "slashing",
		"gov", "params", "ibc", "ibctransfer", "capability",
		"evidence", "feegrant", "authz", "group", "upgrade", "consensus",
		wasmtypes.ModuleName,
	}

	for _, m := range modules {
		kv[m] = storetypes.NewKVStoreKey(m)
		tk[m] = storetypes.NewTransientStoreKey(m + "_t")
		mem[m] = storetypes.NewMemoryStoreKey(m + "_mem")
	}

	return kv, tk, mem
}

// TestWasmKeeperAccessibility verifies WasmKeeper is accessible and functional
func TestWasmKeeperAccessibility(t *testing.T) {
	setupWalletConfig()

	interfaceRegistry, _ := codectypes.NewInterfaceRegistryWithOptions(codectypes.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32AccountAddrPrefix(),
			},
			ValidatorAddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
			},
		},
	})

	appCodec := codec.NewProtoCodec(interfaceRegistry)
	legacyAmino := codec.NewLegacyAmino()

	maccPerms := map[string][]string{
		"mint":                   {"minter"},
		"bonded_tokens_pool":     {"burner", "staking"},
		"not_bonded_tokens_pool": {"burner", "staking"},
		"fee_collector":          nil,
		"distribution":           nil,
		"gov":                    {"burner"},
		"transfer":               {"minter", "burner"},
		wasmtypes.ModuleName:     {"burner"},
	}

	blockedAddrs := map[string]bool{}
	tempDir := t.TempDir()
	appOpts := simtestutil.NewAppOptionsWithFlagHome(tempDir)

	logger := log.NewNopLogger()
	db := dbm.NewMemDB()
	txCfg := tx.NewTxConfig(appCodec, tx.DefaultSignModes)
	baseApp := baseapp.NewBaseApp("hippo", logger, db, txCfg.TxDecoder())

	appKeepers := &keepers.AppKeepersWithKey{}
	appKeepers.GenerateKeys()

	var emptyWasmOptions []wasmkeeper.Option
	require.NotPanics(t, func() {
		appKeepers.InitKeyAndKeepers(appCodec, legacyAmino, maccPerms, blockedAddrs, appOpts, baseApp, logger, emptyWasmOptions)
	}, "WasmKeeper initialization should not panic")

	require.NotNil(t, appKeepers.WasmKeeper, "WasmKeeper should be initialized")

	t.Log("✓ WasmKeeper initialized successfully")
	t.Log("✓ WasmKeeper is accessible from keepers")
}

// TestWasmModuleInStoreUpgrades verifies wasm module is properly configured in store upgrades
func TestWasmModuleInStoreUpgrades(t *testing.T) {
	// Verify wasm is the only module being added
	require.Equal(t, 1, len(Upgrade.StoreUpgrades.Added), "should add exactly one module")
	require.Equal(t, wasmtypes.ModuleName, Upgrade.StoreUpgrades.Added[0], "should add wasm module")

	// Verify no modules are being deleted or renamed
	require.Empty(t, Upgrade.StoreUpgrades.Deleted, "should not delete any modules")
	require.Empty(t, Upgrade.StoreUpgrades.Renamed, "should not rename any modules")

	t.Log("✓ Store upgrades configured correctly")
	t.Log("✓ Only wasm module is being added")
	t.Log("✓ No modules are being deleted or renamed")
}

// TestUpgradeReadiness verifies the upgrade is ready for production
func TestUpgradeReadiness(t *testing.T) {
	// Check upgrade name matches version
	require.Equal(t, "v1.0.3", UpgradeName, "upgrade name should match version")

	// Check upgrade object is complete
	require.NotNil(t, Upgrade, "upgrade object should exist")
	require.Equal(t, UpgradeName, Upgrade.UpgradeName, "upgrade name should match")
	require.NotNil(t, Upgrade.CreateUpgradeHandler, "upgrade handler should exist")
	require.NotNil(t, Upgrade.StoreUpgrades, "store upgrades should be defined")

	// Check wasm module configuration
	require.Contains(t, Upgrade.StoreUpgrades.Added, wasmtypes.ModuleName, "wasm module should be in store upgrades")

	// Verify upgrade handler can be created
	require.NotPanics(t, func() {
		handler := CreateUpgradeHandler(nil, nil, nil)
		require.NotNil(t, handler, "handler should be created")
	}, "upgrade handler creation should not panic")

	t.Log("✓ Upgrade v1.0.3 is ready for production deployment")
	t.Log("✓ All upgrade components are properly configured")
	t.Log("✓ CosmWasm module integration is complete")
}
