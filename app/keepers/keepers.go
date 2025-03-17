package keepers

import (
	storetypes "cosmossdk.io/store/types"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"

	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
	"github.com/spf13/cast"
)

type AppKeepersWithKey struct {
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             govkeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	AuthzKeeper           authzkeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	FeeGrantKeeper        feegrantkeeper.Keeper
	GroupKeeper           groupkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper

	IBCKeeper      *ibckeeper.Keeper        // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	TransferKeeper ibctransferkeeper.Keeper // for cross-chain fungible token transfers

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey
}

func (appKeepers *AppKeepersWithKey) InitKeyAndKeepers(
	appCodec codec.Codec,
	legacyAmino *codec.LegacyAmino,
	maccPerms map[string][]string,
	blockedAddrs map[string]bool,
	appOpts servertypes.AppOptions,
	bApp *baseapp.BaseApp,
	logger log.Logger,
) {
	appKeepers.GenerateKeys()

	appKeepers.ParamsKeeper = initParamsKeeper(appCodec, legacyAmino, appKeepers.keys[paramstypes.StoreKey], appKeepers.tkeys[paramstypes.TStoreKey])

	// From here, We makes keepers.
	// runtime.NewKVStoreService wrapper should be used for some keepers https://docs.cosmos.network/v0.50/build/migrations/upgrading#module-wiring
	// e.g. keys[consensusparamtypes.StoreKey] -> runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey])
	//
	// Also new parameter runtime.EventService{} added https://github.com/cosmos/cosmos-sdk/pull/15547/files#diff-8d1ca8086ee74e8f0490825ba21e7435be4753922192ff691311483aa3e71a0a

	// set the BaseApp's parameter store
	appKeepers.ConsensusParamsKeeper = consensusparamkeeper.NewKeeper(appCodec,
		runtime.NewKVStoreService(appKeepers.keys[consensusparamtypes.StoreKey]), authtypes.NewModuleAddress(govtypes.ModuleName).String(), runtime.EventService{})

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, appKeepers.keys[capabilitytypes.StoreKey], appKeepers.memKeys[capabilitytypes.MemStoreKey])

	// grant capabilities for the ibc and ibc-transfer modules
	appKeepers.ScopedIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	appKeepers.ScopedTransferKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)

	// add keepers
	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(appCodec, runtime.NewKVStoreService(appKeepers.keys[authtypes.StoreKey]), authtypes.ProtoBaseAccount, maccPerms, authcodec.NewBech32Codec(consensus.AddrPrefix), consensus.AddrPrefix, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(appKeepers.keys[banktypes.StoreKey]),
		appKeepers.AccountKeeper,
		blockedAddrs,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		logger,
	)

	appKeepers.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[stakingtypes.StoreKey]), appKeepers.AccountKeeper, appKeepers.BankKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String(), authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()), authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)
	appKeepers.MintKeeper = mintkeeper.NewKeeper(appCodec, runtime.NewKVStoreService(appKeepers.keys[minttypes.StoreKey]), appKeepers.StakingKeeper, appKeepers.AccountKeeper, appKeepers.BankKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	appKeepers.DistrKeeper = distrkeeper.NewKeeper(appCodec, runtime.NewKVStoreService(appKeepers.keys[distrtypes.StoreKey]), appKeepers.AccountKeeper, appKeepers.BankKeeper, appKeepers.StakingKeeper, authtypes.FeeCollectorName, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	appKeepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, legacyAmino, runtime.NewKVStoreService(appKeepers.keys[slashingtypes.StoreKey]), appKeepers.StakingKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(appCodec, runtime.NewKVStoreService(appKeepers.keys[feegrant.StoreKey]), appKeepers.AccountKeeper)

	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(runtime.NewKVStoreService(appKeepers.keys[authzkeeper.StoreKey]), appCodec, bApp.MsgServiceRouter(), appKeepers.AccountKeeper)

	groupConfig := group.DefaultConfig()
	/*
		Example of setting group params:
		groupConfig.MaxMetadataLen = 1000
	*/
	appKeepers.GroupKeeper = groupkeeper.NewKeeper(appKeepers.keys[group.StoreKey], appCodec, bApp.MsgServiceRouter(), appKeepers.AccountKeeper, groupConfig)

	// get skipUpgradeHeights from the app options
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	// set the governance module account as the authority for conducting upgrades
	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, runtime.NewKVStoreService(appKeepers.keys[upgradetypes.StoreKey]), appCodec, homePath, bApp, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	// appKeepers.AccountKeeper.AddressCodec(): https://github.com/cosmos/cosmos-sdk/pull/15825/files#diff-8d1ca8086ee74e8f0490825ba21e7435be4753922192ff691311483aa3e71a0a
	// runtime.ProvideCometInfoService(): https://github.com/cosmos/cosmos-sdk/pull/15850
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[evidencetypes.StoreKey]), appKeepers.StakingKeeper, appKeepers.SlashingKeeper, appKeepers.AccountKeeper.AddressCodec(), runtime.ProvideCometInfoService(),
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	appKeepers.EvidenceKeeper = *evidenceKeeper

	// Create IBC Keeper
	// IBC Keepers need new parameter authority, that is authtypes.NewModuleAddress(govtypes.ModuleName).String()
	// https://ibc.cosmos.network/v8/migrations/v7-to-v8/#authority
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, appKeepers.keys[ibcexported.StoreKey], appKeepers.GetSubspace(ibcexported.ModuleName), appKeepers.StakingKeeper, appKeepers.UpgradeKeeper, appKeepers.ScopedIBCKeeper, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	govConfig := govtypes.DefaultConfig()
	govConfig.MaxMetadataLen = 10200
	appKeepers.GovKeeper = *govkeeper.NewKeeper(
		appCodec, runtime.NewKVStoreService(appKeepers.keys[govtypes.StoreKey]), appKeepers.AccountKeeper, appKeepers.BankKeeper,
		appKeepers.StakingKeeper, appKeepers.DistrKeeper, bApp.MsgServiceRouter(), govConfig, authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Create Transfer Keepers
	appKeepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, appKeepers.keys[ibctransfertypes.StoreKey], appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper, appKeepers.BankKeeper, appKeepers.ScopedTransferKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Set router

	// Register the proposal types
	// Deprecated: Avoid adding new handlers, instead use the new proposal flow
	// by granting the governance module the right to execute the message.
	// See: https://docs.cosmos.network/main/modules/gov#proposal-messages
	govRouter := govv1beta1.NewRouter()
	// Legacy 3 handler(upgrade, ibc)for govv1beta removed.
	// https://github.com/cosmos/cosmos-sdk/pull/16845
	// https://github.com/cosmos/ibc-go/pull/4620/files
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper))

	// Set legacy router for backwards compatibility with gov v1beta1
	appKeepers.GovKeeper.SetLegacyRouter(govRouter)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transfer.NewIBCModule(appKeepers.TransferKeeper))

	// Setting Router will finalize all routes by sealing router
	// No more routes can be added
	appKeepers.IBCKeeper.SetRouter(ibcRouter)
}

func initParamsKeeper(appCodec codec.Codec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)

	return paramsKeeper
}

func (appKeepers *AppKeepersWithKey) SetupHooks() {
	appKeepers.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(appKeepers.DistrKeeper.Hooks(), appKeepers.SlashingKeeper.Hooks()),
	)
}
