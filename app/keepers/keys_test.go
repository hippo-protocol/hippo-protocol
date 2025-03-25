package keepers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hippocrat-dao/hippo-protocol/app/keepers"

	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
)

func TestGenerateKeys(t *testing.T) {
	appKeepers := &keepers.AppKeepersWithKey{}

	appKeepers.GenerateKeys()

	expectedKeys := []string{
		authtypes.StoreKey,
		banktypes.StoreKey,
		stakingtypes.StoreKey,
		minttypes.StoreKey,
		distrtypes.StoreKey,
		slashingtypes.StoreKey,
		govtypes.StoreKey,
		paramstypes.StoreKey,
		consensusparamtypes.StoreKey,
		upgradetypes.StoreKey,
		feegrant.StoreKey,
		evidencetypes.StoreKey,
		capabilitytypes.StoreKey,
		authzkeeper.StoreKey,
		group.StoreKey,
		ibcexported.StoreKey,
		ibctransfertypes.StoreKey,
	}

	for _, key := range expectedKeys {
		require.NotNil(t, appKeepers.GetKey(key), "Key should not be nil: %s", key)
	}

	expectedTKey := appKeepers.GetTransientStoreKey()[paramstypes.TStoreKey]
	expectedMemKey := appKeepers.GetMemoryStoreKey()[capabilitytypes.MemStoreKey]
	expectedKVKey := appKeepers.GetKVStoreKey()[banktypes.StoreKey]
	actualTKey := appKeepers.GetTKey(paramstypes.TStoreKey)
	actualMemKey := appKeepers.GetMemKey(capabilitytypes.MemStoreKey)
	actualKVKey := appKeepers.GetKey(banktypes.StoreKey)

	require.NotNil(t, expectedTKey, "TransientStoreKey should not be nil: params")
	require.NotNil(t, expectedMemKey, "MemoryStoreKey should not be nil: capability")
	require.NotNil(t, expectedKVKey, "KVStoreKey should not be nil: bank")

	require.NotNil(t, actualTKey, "TransientStoreKey for params should not be nil")
	require.Equal(t, expectedTKey, actualTKey, "GetTKey should return the correct transient store key for params")

	require.NotNil(t, actualMemKey, "MemoryStoreKey for capability should not be nil")
	require.Equal(t, expectedMemKey, actualMemKey, "GetMemKey should return the correct memory store key for capability")

	require.NotNil(t, actualKVKey, "KVStoreKey for bank should not be nil")
	require.Equal(t, expectedKVKey, actualKVKey, "GetKVKey should return the correct KV store key for bank")
}
