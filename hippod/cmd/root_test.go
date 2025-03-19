package cmd

import (
	"testing"

	"cosmossdk.io/log"
	cmbtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/hippocrat-dao/hippo-protocol/app"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestNewRootCmd(t *testing.T) {
	rootCmd := NewRootCmd()

	require.NotNil(t, rootCmd, "rootCmd should not be nil")
	require.IsType(t, &cobra.Command{}, rootCmd, "rootCmd should be of type *cobra.Command")
	require.Equal(t, "hippod", rootCmd.Use, "Command name should be 'hippod'")
	require.NotEmpty(t, rootCmd.Commands(), "rootCmd should have subcommands")
	require.Equal(t, "Hippo App", rootCmd.Short, "Command name should be 'hippod'")

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
}
