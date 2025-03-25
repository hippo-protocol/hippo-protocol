package keepers

import (
	"testing"

	"cosmossdk.io/log"

	store "cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/tx/signing"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/address"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/gogoproto/proto"
	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
	"github.com/stretchr/testify/require"
)

// helper function to generate dummy keys
func generateTestStoreKeys() (map[string]*storetypes.KVStoreKey, map[string]*storetypes.TransientStoreKey, map[string]*storetypes.MemoryStoreKey) {
	kv := make(map[string]*storetypes.KVStoreKey)
	tk := make(map[string]*storetypes.TransientStoreKey)
	mem := make(map[string]*storetypes.MemoryStoreKey)

	// Add only needed keys (for simplicity), you can expand this
	modules := []string{
		"auth", "bank", "staking", "mint", "distr", "slashing",
		"gov", "params", "ibc", "ibctransfer", "capability",
		"evidence", "feegrant", "authz", "group", "upgrade", "consensus",
	}

	for _, m := range modules {
		kv[m] = storetypes.NewKVStoreKey(m)
		tk[m] = storetypes.NewTransientStoreKey(m + "_t")
		mem[m] = storetypes.NewMemoryStoreKey(m + "_mem")
	}

	return kv, tk, mem
}

func TestInitKeyAndKeepers(t *testing.T) {
	consensus.SetWalletConfig()

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
		"gov":                    nil,
		"transfer":               {"minter", "burner"},
	}

	blockedAddrs := map[string]bool{}
	appOpts := server.NewDefaultContext().Viper

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

	appKeepers := &AppKeepersWithKey{
		keys:    keys,
		tkeys:   tkeys,
		memKeys: memKeys,
	}

	require.NotPanics(t, func() {
		appKeepers.InitKeyAndKeepers(appCodec, legacyAmino, maccPerms, blockedAddrs, appOpts, baseApp, logger)
	})

	// Simple verification of important keepers
	require.NotNil(t, appKeepers.AccountKeeper)
	require.NotNil(t, appKeepers.BankKeeper)
	require.NotNil(t, appKeepers.StakingKeeper)
	require.NotNil(t, appKeepers.ParamsKeeper)
	require.NotNil(t, appKeepers.IBCKeeper)
	require.NotNil(t, appKeepers.CapabilityKeeper)
	require.NotNil(t, appKeepers.UpgradeKeeper)
}

func TestSetupHooks(t *testing.T) {
	appCodec := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	legacyAmino := codec.NewLegacyAmino()
	maccPerms := map[string][]string{
		"mint":                   {"minter"},
		"bonded_tokens_pool":     {"burner", "staking"},
		"not_bonded_tokens_pool": {"burner", "staking"},
		"fee_collector":          nil,
		"distribution":           nil,
		"gov":                    nil,
		"transfer":               {"minter", "burner"},
	}
	blockedAddrs := map[string]bool{}
	appOpts := server.NewDefaultContext().Viper
	logger := log.NewNopLogger()

	db := dbm.NewMemDB()
	keys, tkeys, memKeys := generateTestStoreKeys()
	ms := store.NewCommitMultiStore(db, logger, nil)
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

	allStoreKeys := make([]storetypes.StoreKey, 0, len(keys)+len(tkeys)+len(memKeys))
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

	appKeepers := &AppKeepersWithKey{
		keys:    keys,
		tkeys:   tkeys,
		memKeys: memKeys,
	}
	appKeepers.InitKeyAndKeepers(appCodec, legacyAmino, maccPerms, blockedAddrs, appOpts, baseApp, logger)

	require.NotPanics(t, func() {
		appKeepers.SetupHooks()
	})

	// Additional check: the staking hooks should not be nil
	require.NotNil(t, appKeepers.StakingKeeper.Hooks())
}
