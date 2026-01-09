package app

import (
	"testing"

	"cosmossdk.io/log"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestPrepForZeroHeightGenesis_NotNil(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	ctx := app.NewContextLegacy(true, cmtproto.Header{Height: 1})
	jailAllowedAddrs := []string{}
	app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs)
	assert.NotPanics(t, func() { app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs) }, "prepForZeroHeightGenesis should not panic")
}

func TestPrepForZeroHeightGenesis_WithValidJailAllowed(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	appInstance := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)

	ctx := appInstance.NewContextLegacy(true, cmtproto.Header{Height: 1})
	addrBytes := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x10,
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x20}
	addr := sdk.ValAddress(addrBytes)

	jailAllowedAddrs := []string{addr.String()}

	assert.NotPanics(t, func() {
		appInstance.prepForZeroHeightGenesis(ctx, jailAllowedAddrs)
	}, "prepForZeroHeightGenesis should not panic with valid jailAllowedAddrs")
}

func TestPrepForZeroHeightGenesis_MultipleCalls(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	appInstance := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	ctx := appInstance.NewContextLegacy(true, cmtproto.Header{Height: 10})
	jailAllowedAddrs := []string{}
	assert.NotPanics(t, func() {
		for i := 0; i < 3; i++ {
			appInstance.prepForZeroHeightGenesis(ctx, jailAllowedAddrs)
		}
	}, "prepForZeroHeightGenesis should not panic when called repeatedly")
}

func TestExportAppStateAndValidators_NotPanics(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()), EmptyWasmOptions)
	ctx := app.NewContextLegacy(true, cmtproto.Header{Height: app.LastBlockHeight()})
	genesisState := app.DefaultGenesis()
	app.ModuleManager.InitGenesis(ctx, app.AppCodec(), genesisState)

	assert.NotPanics(t, func() { app.ExportAppStateAndValidators(true, []string{}, []string{}) }, "ExportAppStateAndValidators should not panic")
}
