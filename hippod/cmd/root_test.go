package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cosmossdk.io/log"
	cmbtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/hippocrat-dao/hippo-protocol/app"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestNewRootCmd(t *testing.T) {
	home := "/tmp/hippo-test"
	_ = os.MkdirAll(home+"/config", 0755)
	_ = os.WriteFile(home+"/config/genesis.json", []byte(`{"genesis_time":"2023-01-01T00:00:00Z","chain_id":"test-chain","app_state":{}}`), 0644)

	app.DefaultNodeHome = home

	rc := NewRootCmd()
	require.NotNil(t, rc)
	require.Equal(t, "hippod", rc.Use)

	rc.SetArgs([]string{"config"})
	buf := new(bytes.Buffer)
	rc.SetOut(buf)
	rc.SetErr(buf)
	err := rc.Execute()
	require.NoError(t, err)

	subCmds := []string{"debug", "config", "completion", "status", "genesis", "query", "tx", "keys"}
	foundCommands := map[string]bool{}
	for _, c := range rc.Commands() {
		for _, name := range subCmds {
			if c.Use == name || strings.HasPrefix(c.Use, name+" ") {
				foundCommands[name] = true
			}
		}
	}

	for _, name := range subCmds {
		require.True(t, foundCommands[name], "expected subcommand %s not found", name)
	}
}

func TestInitCometBFTConfig(t *testing.T) {
	config := initCometBFTConfig()
	defaultConfig := cmbtcfg.DefaultConfig()
	require.Equal(t, config, defaultConfig)
}

func TestInitAppConfig(t *testing.T) {
	defaultConfig, _ := initAppConfig()

	require.Equal(t, serverconfig.DefaultConfigTemplate, defaultConfig)
	// add test for min gas price
}

func TestAppExport(t *testing.T) {
	exportedApp, err := appExport(log.NewNopLogger(), dbm.NewMemDB(), nil, 0, true, nil, simtestutil.NewAppOptionsWithFlagHome(app.DefaultNodeHome), nil)
	require.Error(t, err)
	require.NotNil(t, exportedApp)

	exportedApp, err = appExport(log.NewNopLogger(), dbm.NewMemDB(), nil, 0, true, nil, simtestutil.NewAppOptionsWithFlagHome(""), nil)
	require.Error(t, err)
	require.NotNil(t, exportedApp)

	exportedApp, err = appExport(log.NewNopLogger(), dbm.NewMemDB(), nil, -1, true, nil, simtestutil.NewAppOptionsWithFlagHome(app.DefaultNodeHome), nil)
	require.Error(t, err)
	require.NotNil(t, exportedApp)
}

type mockAppOptions struct {
	options map[string]interface{}
}

func (m mockAppOptions) Get(key string) interface{} {
	if val, ok := m.options[key]; ok {
		return val
	}
	return nil
}

func setupGenesisFile(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir() // return a new temp dir
	configDir := filepath.Join(tmpDir, "config")
	err := os.Mkdir(configDir, 0755) // create a new dir
	require.NoError(t, err)

	genesisPath := filepath.Join(configDir, "genesis.json")
	err = os.WriteFile(genesisPath, []byte(`{"chain_id":"test-chain"}`), 0644) // read minimum genesis file
	require.NoError(t, err)

	return tmpDir
}

func TestNewApp(t *testing.T) {
	logger := log.Logger(log.NewNopLogger())
	db := dbm.NewMemDB()
	traceStore := new(bytes.Buffer)

	tmpHome := setupGenesisFile(t)

	appOpts := mockAppOptions{
		options: map[string]interface{}{
			"home":                     tmpHome,
			server.FlagPruning:         "nothing",    // or "default" / "everything" / "nothing"
			server.FlagMinGasPrices:    "0.001uatom", // minimum gas fees
			server.FlagHaltHeight:      uint64(0),    // no automatic halt
			server.FlagHaltTime:        uint64(0),
			server.FlagInterBlockCache: true,
			server.FlagIndexEvents:     []string{"tx.height", "tx.hash"},
			server.FlagIAVLCacheSize:   781250, // size of the IAVL cache
		},
	}

	appInstance := newApp(logger, db, traceStore, appOpts)
	require.NotNil(t, appInstance, "Should not be nil")
}

type fakeAppOptions map[string]interface{}

func (f fakeAppOptions) Get(key string) interface{} {
	return f[key]
}

func TestAddModuleInitFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "start"}
	require.NotPanics(t, func() {
		addModuleInitFlags(cmd)
	})
}

func TestGenesisCommand(t *testing.T) {
	cdc := makeTestEncodingConfig()
	txConfig := tx.NewTxConfig(cdc, tx.DefaultSignModes)
	basicManager := module.NewBasicManager(
		genutil.AppModuleBasic{},
	)
	customCmd := &cobra.Command{Use: "custom"}

	genCmd := genesisCommand(txConfig, basicManager, customCmd)
	require.NotNil(t, genCmd)
	found := false
	for _, c := range genCmd.Commands() {
		if c.Use == "custom" {
			found = true
		}
	}
	require.True(t, found, "custom command not found in genesis command")
}

func TestQueryCommand(t *testing.T) {
	qc := queryCommand()
	require.NotNil(t, qc)
	require.Equal(t, "query", qc.Use)
	require.Greater(t, len(qc.Commands()), 0)
}

func TestTxCommand(t *testing.T) {
	tc := txCommand()
	require.NotNil(t, tc)
	require.Equal(t, "tx", tc.Use)
	require.Greater(t, len(tc.Commands()), 0)
}

// makeMinimalAppOptions returns minimal valid app options to prevent nil panic
func makeMinimalAppOptions() fakeAppOptions {
	_ = os.MkdirAll("/tmp/hippo-test/config", 0755)
	_ = os.WriteFile("/tmp/hippo-test/config/genesis.json", []byte(`{"genesis_time":"2023-01-01T00:00:00Z","chain_id":"test-chain","app_state":{}}`), 0644)

	return fakeAppOptions{
		"home":               "/tmp/hippo-test",
		"trace":              false,
		"inv-check-period":   uint(1),
		"pruning":            "default",
		"minimum-gas-prices": "0stake",
	}
}

func TestOverwriteFlagDefaults(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	child := &cobra.Command{Use: "child"}
	cmd.AddCommand(child)

	cmd.PersistentFlags().String("chain-id", "", "")
	cmd.PersistentFlags().String("keyring-backend", "", "")
	cmd.Flags().String("chain-id", "", "")
	cmd.Flags().String("keyring-backend", "", "")

	require.NotPanics(t, func() {
		overwriteFlagDefaults(cmd, map[string]string{
			"chain-id":        "hippo-test",
			"keyring-backend": "test",
		})
	})
}

func TestAppExport_InvalidHome(t *testing.T) {
	_, err := appExport(log.NewNopLogger(), nil, nil, -1, true, nil, fakeAppOptions{"home": nil}, nil)
	require.Error(t, err)
	require.Equal(t, "application home not set", err.Error())
}

func TestAppExport_InvalidViper(t *testing.T) {
	_, err := appExport(log.NewNopLogger(), nil, nil, -1, true, nil, fakeAppOptions{"home": "home"}, nil)
	require.Error(t, err)
	require.Equal(t, "appOpts is not viper.Viper", err.Error())
}
