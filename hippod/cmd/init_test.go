package cmd

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	auth "github.com/cosmos/cosmos-sdk/x/auth"
	bank "github.com/cosmos/cosmos-sdk/x/bank"
	mint "github.com/cosmos/cosmos-sdk/x/mint"
	staking "github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/types/module"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	distribution "github.com/cosmos/cosmos-sdk/x/distribution"
	gov "github.com/cosmos/cosmos-sdk/x/gov"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing"

	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func makeTestEncodingConfig() codec.Codec {
	interfaceRegistry := types.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(interfaceRegistry)
	authvesting.RegisterInterfaces(interfaceRegistry)
	banktypes.RegisterInterfaces(interfaceRegistry)
	govv1types.RegisterInterfaces(interfaceRegistry)
	minttypes.RegisterInterfaces(interfaceRegistry)
	stakingtypes.RegisterInterfaces(interfaceRegistry)
	slashingtypes.RegisterInterfaces(interfaceRegistry)
	distrtypes.RegisterInterfaces(interfaceRegistry)

	return codec.NewProtoCodec(interfaceRegistry)
}

func TestInitCmd_Basic(t *testing.T) {
	home := t.TempDir()
	defaultNodeHome := home
	viper.Set("home", home)

	basicManager := module.NewBasicManager(
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		gov.AppModuleBasic{},
		slashing.AppModuleBasic{},
		distribution.AppModuleBasic{},
	)

	command := InitCmd(basicManager, defaultNodeHome)
	command.SetArgs([]string{"hippo-moniker", "--home", home, "--overwrite"})

	configPath := filepath.Join(home, "config")
	require.NoError(t, os.MkdirAll(configPath, 0755))

	srvCtx := server.NewDefaultContext()
	srvCtx.Config.RootDir = home
	clientCtx := client.Context{}.WithHomeDir(home).WithCodec(makeTestEncodingConfig())

	command.SetContext(context.WithValue(
		context.WithValue(context.Background(), server.ServerContextKey, srvCtx),
		client.ClientContextKey, &clientCtx,
	))

	require.NoError(t, command.Execute())

	expectedFiles := []string{
		"genesis.json",
		"node_key.json",
		"priv_validator_key.json",
	}

	for _, f := range expectedFiles {
		_, err := os.Stat(filepath.Join(configPath, f))
		require.NoError(t, err, "%s should exist", f)
	}
}

func TestInitCmd_RecoverMnemonic_Invalid(t *testing.T) {
	home := t.TempDir()
	defaultNodeHome := home
	basicManager := module.NewBasicManager()

	command := InitCmd(basicManager, defaultNodeHome)
	command.SetArgs([]string{"hippo-moniker", "--home", home, "--recover"})

	r, w, err := os.Pipe()
	require.NoError(t, err)
	command.SetIn(r)
	_, err = w.WriteString("invalid mnemonic phrase\n")
	require.NoError(t, err)
	w.Close()

	srvCtx := server.NewDefaultContext()
	srvCtx.Config.RootDir = home
	clientCtx := client.Context{}.WithHomeDir(home).WithCodec(makeTestEncodingConfig())
	command.SetContext(context.WithValue(
		context.WithValue(context.Background(), server.ServerContextKey, srvCtx),
		client.ClientContextKey, &clientCtx,
	))

	err = command.Execute()
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid mnemonic")
}

func TestOverrideGenesis_Valid(t *testing.T) {
	cdc := makeTestEncodingConfig()
	appState := map[string]json.RawMessage{
		stakingtypes.ModuleName:  cdc.MustMarshalJSON(stakingtypes.DefaultGenesisState()),
		minttypes.ModuleName:     cdc.MustMarshalJSON(minttypes.DefaultGenesisState()),
		distrtypes.ModuleName:    cdc.MustMarshalJSON(distrtypes.DefaultGenesisState()),
		slashingtypes.ModuleName: cdc.MustMarshalJSON(slashingtypes.DefaultGenesisState()),
		govtypes.ModuleName:      cdc.MustMarshalJSON(govv1types.DefaultGenesisState()),
	}
	genDoc := &tmtypes.GenesisDoc{
		ConsensusParams: tmtypes.DefaultConsensusParams(),
	}

	_, err := overrideGenesis(cdc, genDoc, appState)
	require.NoError(t, err)
}

func TestDisplayInfo(t *testing.T) {
	info := printInfo{
		Moniker:    "hippo",
		ChainID:    "hippo-1",
		NodeID:     "node123",
		GenTxsDir:  "",
		AppMessage: json.RawMessage(`{"foo":"bar"}`),
	}

	err := displayInfo(info)
	require.NoError(t, err)
}
