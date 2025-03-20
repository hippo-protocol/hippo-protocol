package app

import (
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	"cosmossdk.io/log"
    cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/stretchr/testify/assert"
)

func TestPrepForZeroHeightGenesis_NotNil(t *testing.T) {
	db := dbm.NewMemDB()
    logger := log.NewTestLogger(t)
    app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()))
	ctx := app.NewContextLegacy(true, cmtproto.Header{Height: 1})
	jailAllowedAddrs := []string{}
    app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs)
	assert.NotPanics(t, func() { app.prepForZeroHeightGenesis(ctx, jailAllowedAddrs) }, "prepForZeroHeightGenesis should not panic")
}

func TestExportAppStateAndValidators_NotPanics(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := New(logger, db, nil, true, NewAppOptionsWithFlagHome(t.TempDir()))
	ctx := app.NewContextLegacy(true, cmtproto.Header{Height: app.LastBlockHeight()})
	genesisState := app.DefaultGenesis()
	app.ModuleManager.InitGenesis(ctx, app.AppCodec(), genesisState)

	assert.NotPanics(t, func() { app.ExportAppStateAndValidators(true, []string{}, []string{}) }, "ExportAppStateAndValidators should not panic")
}
