package cmd_test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	testcmd "github.com/hippocrat-dao/hippo-protocol/hippod/cmd"
)

func mustAccAddress(bytes []byte) sdk.AccAddress {
	addr := sdk.AccAddress(bytes)
	return addr
}

func createGenesisFile(t *testing.T, home string, accounts []authtypes.GenesisAccount, balances []banktypes.Balance) {
	cdc, _ := newTestCodecAndRegistry()

	packedAccounts, err := authtypes.PackAccounts(accounts)
	require.NoError(t, err)

	authState := authtypes.GenesisState{
		Accounts: packedAccounts,
		Params:   authtypes.DefaultParams(),
	}
	authJSON, err := cdc.MarshalJSON(&authState)
	require.NoError(t, err)

	bankState := banktypes.DefaultGenesisState()
	bankState.Balances = balances
	bankJSON, err := cdc.MarshalJSON(bankState)
	require.NoError(t, err)

	appState := map[string]json.RawMessage{
		authtypes.ModuleName: authJSON,
		banktypes.ModuleName: bankJSON,
	}

	appStateJSON, err := json.Marshal(appState)
	require.NoError(t, err)

	genDoc := &tmtypes.GenesisDoc{
		ChainID:         "test-chain",
		AppState:        appStateJSON,
		ConsensusParams: tmtypes.DefaultConsensusParams(),
	}

	configPath := filepath.Join(home, "config")
	require.NoError(t, os.MkdirAll(configPath, 0755))
	genFile := filepath.Join(configPath, "genesis.json")
	require.NoError(t, genDoc.SaveAs(genFile))

	srvCfgPath := filepath.Join(configPath, "config.toml")
	if _, err := os.Stat(srvCfgPath); os.IsNotExist(err) {
		require.NoError(t, os.WriteFile(srvCfgPath, []byte("[api]\n"), 0644))
	}
}

func newTestCodecAndRegistry() (*codec.ProtoCodec, codectypes.InterfaceRegistry) {
	registry := codectypes.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(registry)
	authvesting.RegisterInterfaces(registry)
	banktypes.RegisterInterfaces(registry)
	return codec.NewProtoCodec(registry), registry
}

func injectClientCtx(cmd *cobra.Command, home string) {
	cdc, registry := newTestCodecAndRegistry()
	clientCtx := client.Context{}.
		WithHomeDir(home).
		WithCodec(cdc).
		WithInterfaceRegistry(registry)

	cmd.SetContext(context.WithValue(context.Background(), client.ClientContextKey, &clientCtx))
}

func TestAddGenesisAccountCmd(t *testing.T) {
	home := t.TempDir()
	srvCtx := server.NewDefaultContext()
	srvCtx.Config.RootDir = home

	addr := mustAccAddress([]byte("testaddr000000000000000000000000000001")).String()
	coins := "100stake"

	t.Run("New account", func(t *testing.T) {
		createGenesisFile(t, home, []authtypes.GenesisAccount{}, []banktypes.Balance{})
		cmd := testcmd.AddGenesisAccountCmd(home)
		args := []string{addr, coins, "--home", home}
		require.NoError(t, cmd.ParseFlags(args))
		injectClientCtx(cmd, home)
		viper.Set(flags.FlagHome, home)
		require.NoError(t, cmd.RunE(cmd, args[:2]))
	})

	t.Run("Append to existing", func(t *testing.T) {
		acc := authtypes.NewBaseAccount(mustAccAddress([]byte("testaddrappend")), nil, 0, 0)
		balance := banktypes.Balance{
			Address: acc.GetAddress().String(),
			Coins:   sdk.NewCoins(sdk.NewInt64Coin("stake", 50)),
		}
		createGenesisFile(t, home, []authtypes.GenesisAccount{acc}, []banktypes.Balance{balance})

		cmd := testcmd.AddGenesisAccountCmd(home)
		args := []string{acc.GetAddress().String(), coins, "--home", home, "--append"}

		require.NoError(t, cmd.ParseFlags(args))
		viper.Set(flags.FlagHome, home)
		injectClientCtx(cmd, home)

		require.NoError(t, cmd.RunE(cmd, args[:2]))
	})

	t.Run("Invalid Bech32", func(t *testing.T) {
		cmd := testcmd.AddGenesisAccountCmd(home)
		args := []string{"invalid_addr_%%%", coins, "--home", home}
		require.NoError(t, cmd.ParseFlags(args))
		viper.Set(flags.FlagHome, home)
		injectClientCtx(cmd, home)
		require.Error(t, cmd.RunE(cmd, args[:2]))
	})

	t.Run("Invalid coins format", func(t *testing.T) {
		cmd := testcmd.AddGenesisAccountCmd(home)
		args := []string{addr, "invalidcoin!", "--home", home}
		require.NoError(t, cmd.ParseFlags(args))
		viper.Set(flags.FlagHome, home)
		injectClientCtx(cmd, home)
		require.Error(t, cmd.RunE(cmd, args[:2]))
	})

	t.Run("Vesting > balance", func(t *testing.T) {
		createGenesisFile(t, home, []authtypes.GenesisAccount{}, []banktypes.Balance{})
		cmd := testcmd.AddGenesisAccountCmd(home)
		args := []string{
			addr, coins,
			"--home", home,
			"--vesting-amount", "150stake",
			"--vesting-end-time", fmt.Sprintf("%d", time.Now().Add(time.Hour).Unix()),
		}
		require.NoError(t, cmd.ParseFlags(args))
		viper.Set(flags.FlagHome, home)
		injectClientCtx(cmd, home)
		require.ErrorContains(t, cmd.RunE(cmd, args[:2]), "vesting amount cannot be greater")
	})

	t.Run("Continuous vesting", func(t *testing.T) {
		createGenesisFile(t, home, []authtypes.GenesisAccount{}, []banktypes.Balance{})
		cmd := testcmd.AddGenesisAccountCmd(home)
		args := []string{
			addr, coins,
			"--home", home,
			"--vesting-amount", "50stake",
			"--vesting-start-time", fmt.Sprintf("%d", time.Now().Unix()),
			"--vesting-end-time", fmt.Sprintf("%d", time.Now().Add(time.Hour).Unix()),
		}
		require.NoError(t, cmd.ParseFlags(args))
		viper.Set(flags.FlagHome, home)
		injectClientCtx(cmd, home)
		require.NoError(t, cmd.RunE(cmd, args[:2]))
	})

	t.Run("Delayed vesting", func(t *testing.T) {
		createGenesisFile(t, home, []authtypes.GenesisAccount{}, []banktypes.Balance{})
		cmd := testcmd.AddGenesisAccountCmd(home)
		args := []string{
			addr, coins,
			"--home", home,
			"--vesting-amount", "50stake",
			"--vesting-end-time", fmt.Sprintf("%d", time.Now().Add(time.Hour).Unix()),
		}
		require.NoError(t, cmd.ParseFlags(args))
		viper.Set(flags.FlagHome, home)
		injectClientCtx(cmd, home)
		require.NoError(t, cmd.RunE(cmd, args[:2]))
	})

	t.Run("Invalid vesting params", func(t *testing.T) {
		cmd := testcmd.AddGenesisAccountCmd(home)
		args := []string{
			addr, coins,
			"--home", home,
			"--vesting-amount", "10stake",
		}
		require.NoError(t, cmd.ParseFlags(args))
		viper.Set(flags.FlagHome, home)
		injectClientCtx(cmd, home)
		require.ErrorContains(t, cmd.RunE(cmd, args[:2]), "invalid vesting parameters")
	})
}
