package app

import (
	"testing"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
	"github.com/stretchr/testify/assert"
)

type AppOptionsMap map[string]interface{}

func (m AppOptionsMap) Get(key string) interface{} {
	v, ok := m[key]
	if !ok {
		return interface{}(nil)
	}

	return v
}

func NewAppOptionsWithFlagHome(homePath string) servertypes.AppOptions {
	return AppOptionsMap{
		flags.FlagHome: homePath,
	}
}

func TestNewApp(t *testing.T) {
	consensus.SetWalletConfig()
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	assert.NotNil(t, app.Name(), "app name should not be nil")
}

func TestAutoCli(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	assert.NotNil(t, app.AutoCliOpts(), "AutoCliOpts should not return nil")
}

func TestUpgradeHandlers(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	assert.NotPanics(t, func() { app.setupUpgradeHandlers() }, "setupUpgradeHandlers should not panic")
}

func TestUpgradeStoreLoaders(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	assert.NotPanics(t, func() { app.setupUpgradeStoreLoaders() }, "setupUpgradeStoreLoaders should not panic")

}

func TestBlockedAddresses(t *testing.T) {
	blockedAddresses := BlockedAddresses()
	assert.NotNil(t, blockedAddresses, "BlockedAddrs should not return nil")
}

func TestGetMaccPerms(t *testing.T) {
	maccPerms := GetMaccPerms()
	assert.NotNil(t, maccPerms, "GetMaccPerms should not return nil")
}

func TestConfigurator(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	assert.NotNil(t, app.Configurator(), "Configurator should not return nil")

}

func TestLegacyAmino(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	assert.NotNil(t, app.LegacyAmino(), "LegacyAmino should not return nil")
}

func TestAppCodec(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	assert.NotNil(t, app.AppCodec(), "AppCodec should not return nil")
}
func TestInterfaceRegistry(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	assert.NotNil(t, app.InterfaceRegistry(), "InterfaceRegistry should not return nil")
}

func TestTxConfig(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	assert.NotNil(t, app.TxConfig(), "TxConfig should not return nil")
}

func TestDefaultGenesis(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	assert.NotNil(t, app.DefaultGenesis(), "DefaultGenesis should not return nil")
}

func TestSimulationManager(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	assert.NotNil(t, app.SimulationManager(), "SimulationManager should not return nil")
}

func TestPreblocker(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	ctx := app.NewContextLegacy(true, cmtproto.Header{Height: app.LastBlockHeight()})
	res, _ := app.PreBlocker(ctx, nil)
	assert.NotNil(t, res, "PreBlocker should not return nil")
}

func TestBeginBlocker(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	ctx := app.NewContextLegacy(true, cmtproto.Header{Height: app.LastBlockHeight()})
	res, _ := app.BeginBlocker(ctx)
	assert.NotNil(t, res, "BeginBlocker should not return nil")
}

func TestEndBlocker(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	ctx := app.NewContextLegacy(true, cmtproto.Header{Height: app.LastBlockHeight()})
	res, _ := app.EndBlocker(ctx)
	assert.NotNil(t, res, "EndBlocker should not return nil")
}

func TestRegisterNodeService(t *testing.T) {
	clientCtx := client.Context{}
	cfg := config.Config{}
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	app.RegisterNodeService(clientCtx, cfg)
	assert.NotNil(t, app, "RegisterNodeService should not return nil")
}

func TestRegisterTendermintService(t *testing.T) {
	clientCtx := client.Context{}
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	app.RegisterTendermintService(clientCtx)
	assert.NotNil(t, app, "RegisterTendermintService should not return nil")
}

func TestRegisterTxService(t *testing.T) {
	clientCtx := client.Context{}
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	app.RegisterTxService(clientCtx)
	assert.NotNil(t, app, "RegisterTxService should not return nil")
}

func TestLoadHeight(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	app.LoadHeight(1)
	assert.NotNil(t, app, "LoadHeight should not return nil")
}

func TestInitChainer_InvalidJSON(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	ctx := app.NewContext(true)

	req := &abci.RequestInitChain{
		AppStateBytes: []byte("invalid json"),
	}

	assert.Panics(t, func() { app.InitChainer(ctx, req) }, "InitChainer should panic with invalid JSON")
}

func TestRegisterAPIRoutes(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)

	clientCtx := client.Context{}
	apiServer := api.New(clientCtx, logger, nil)
	apiConfig := config.APIConfig{}

	assert.NotNil(t, app, "App should not be nil before calling RegisterAPIRoutes")
	assert.NotNil(t, apiServer, "API Server should be initialized")

	app.RegisterAPIRoutes(apiServer, apiConfig)

	assert.NotNil(t, apiServer.GRPCGatewayRouter, "GRPCGatewayRouter should be initialized after RegisterAPIRoutes")
}
