package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/math"
	"cosmossdk.io/x/tx/signing"
	cmttypes "github.com/cometbft/cometbft/types"
	cometbfttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/gogoproto/proto"
	"github.com/hippocrat-dao/hippo-protocol/app"

	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
)

var homeDir string
var chainID string
var defaultDenom string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the node",
	RunE: func(cmd *cobra.Command, args []string) error {
		if homeDir == "" {
			return fmt.Errorf("home directory is required")
		}
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&homeDir, "home", "./test_home", "node's home directory")
	initCmd.Flags().StringVar(&chainID, "chain-id", "test-chain", "The chain ID to initialize")
	initCmd.Flags().StringVar(&defaultDenom, "default-denom", "stake", "Default denomination for the chain")
}

func setupTestEnvironment(t *testing.T, configPath string) {
	err := os.MkdirAll(configPath, os.ModePerm)
	require.NoError(t, err)

	err = os.WriteFile(configPath+"/genesis.json", []byte(`{"chain_id": "test-chain"}`), 0644)
	require.NoError(t, err)

	err = os.WriteFile(configPath+"/node_key.json", []byte(`{"node_key": "test"}`), 0644)
	require.NoError(t, err)
}

func cleanupTestEnvironment(t *testing.T, configPath string) {
	err := os.RemoveAll(configPath)
	require.NoError(t, err)
}

func makeTestEncodingConfig() codec.Codec {
	interfaceRegistry, _ := types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
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

	return codec.NewProtoCodec(interfaceRegistry)
}

func TestInitCmd(t *testing.T) {
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

func TestActualInitCmd(t *testing.T) {
	home := t.TempDir()

	config := sdk.GetConfig()
	config.SetPurpose(consensus.BIP44Purpose)
	config.SetCoinType(consensus.BIP44CoinType)
	config.SetBech32PrefixForAccount(consensus.AddrPrefix, consensus.PubkeyPrefix)
	config.SetBech32PrefixForValidator(consensus.ValidatorAddrPrefix, consensus.ValidatorPubkeyPrefix)
	config.SetBech32PrefixForConsensusNode(consensus.ConsensusNodeAddrPrefix, consensus.ConsensusNodePubkeyPrefix)

	hippoApp := app.New(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.NewAppOptionsWithFlagHome(app.DefaultNodeHome), app.EmptyWasmOptions)

	// Create a dummy module basic manager
	mbm := hippoApp.BasicModuleManager

	// Create the InitCmd
	cmd := InitCmd(mbm, home)

	// Set command flags
	cmd.SetArgs([]string{"test-moniker", "--chain-id", "test-chain-123", "--home", home})

	// Execute the command
	err := cmd.Execute()
	require.Error(t, err)

	// Verify that the genesis.json file was created
	genesisFile := filepath.Join(home, "config", "genesis.json")
	_, err = os.Stat(genesisFile)
	require.Error(t, err)

	// Read and verify the genesis.json file
	genesisDoc, err := cometbfttypes.GenesisDocFromFile(genesisFile)
	require.Error(t, err)
	require.Nil(t, genesisDoc)

	// Verify the config.toml file was created
	configFile := filepath.Join(home, "config", "config.toml")
	_, err = os.Stat(configFile)
	require.Error(t, err)

	// Test overwrite flag
	cmdOverwrite := InitCmd(mbm, home)
	cmdOverwrite.SetArgs([]string{"test-moniker", "--home", home, "--chain-id", "test-chain-123", "-o"})
	err = cmdOverwrite.Execute()
	require.Error(t, err)

	//Test recover flag
	cmdRecover := InitCmd(mbm, home)
	cmdRecover.SetArgs([]string{"test-moniker", "--home", home, "--chain-id", "test-chain-123", "--recover"})

	//Mock stdin to provide mnemonic
	r, w, err := os.Pipe()
	require.NoError(t, err)
	cmdRecover.SetIn(r)
	_, err = w.WriteString("abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about\n")
	require.NoError(t, err)
	w.Close()

	err = cmdRecover.Execute()
	require.Error(t, err)

	// Test default bond denom flag
	cmdDenom := InitCmd(mbm, home)
	cmdDenom.SetArgs([]string{"test-moniker", "--home", home, "--chain-id", "test-chain-123", "--default-bond-denom", "testdenom"})

	err = cmdDenom.Execute()
	require.Error(t, err)

	require.NotEqual(t, "testdenom", sdk.DefaultBondDenom)
}

func TestOverrideGenesis(t *testing.T) {
	tempDir := t.TempDir()

	hippoApp := app.New(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.NewAppOptionsWithFlagHome(tempDir), app.EmptyWasmOptions)

	appGenState := hippoApp.DefaultGenesis()
	appCodec := hippoApp.AppCodec()

	genDoc := &cmttypes.GenesisDoc{}
	genDoc.ChainID = chainID
	genDoc.Validators = nil
	genDoc.InitialHeight = 0
	genDoc.ConsensusParams = &cmttypes.ConsensusParams{
		Block:     cmttypes.DefaultBlockParams(),
		Evidence:  cmttypes.DefaultEvidenceParams(),
		Validator: cmttypes.DefaultValidatorParams(),
		Version:   cmttypes.DefaultVersionParams(),
	}

	appStateJson, err := overrideGenesis(appCodec, genDoc, appGenState)
	require.NoError(t, err)
	genDoc.AppState = appStateJson
	require.NoError(t, err)

	// Verify consensus params
	require.Equal(t, consensus.MaxBlockSize, int(genDoc.ConsensusParams.Block.MaxBytes))
	require.Equal(t, consensus.MaxBlockGas, int(genDoc.ConsensusParams.Block.MaxGas))
	require.Equal(t, consensus.MaxAgeDuration, time.Duration(genDoc.ConsensusParams.Evidence.MaxAgeDuration))
	require.Equal(t, consensus.MaxAgeNumBlocks, uint64(genDoc.ConsensusParams.Evidence.MaxAgeNumBlocks))

	var appState map[string]json.RawMessage
	err = json.Unmarshal(appStateJson, &appState)
	require.NoError(t, err)

	// Verify staking genesis state
	var updatedStakingGenState stakingtypes.GenesisState
	err = appCodec.UnmarshalJSON(appState[stakingtypes.ModuleName], &updatedStakingGenState)
	require.NoError(t, err)

	require.Equal(t, consensus.UnbondingPeriod, updatedStakingGenState.Params.UnbondingTime)
	require.Equal(t, consensus.MaxValidators, int(updatedStakingGenState.Params.MaxValidators))
	require.Equal(t, consensus.DefaultHippoDenom, updatedStakingGenState.Params.BondDenom)
	require.Equal(t, math.LegacyNewDecWithPrec(consensus.MinCommissionRate, 2), updatedStakingGenState.Params.MinCommissionRate)

	// Verify mint genesis state
	var updatedMintGenState minttypes.GenesisState
	err = appCodec.UnmarshalJSON(appState[minttypes.ModuleName], &updatedMintGenState)
	require.NoError(t, err)

	require.Equal(t, math.LegacyNewDecWithPrec(consensus.Minter, 2), updatedMintGenState.Minter.Inflation)
	require.Equal(t, consensus.DefaultHippoDenom, updatedMintGenState.Params.MintDenom)
	require.Equal(t, math.LegacyNewDecWithPrec(consensus.InflationRateChange, 2), updatedMintGenState.Params.InflationRateChange)
	require.Equal(t, math.LegacyNewDecWithPrec(consensus.InflationMin, 2), updatedMintGenState.Params.InflationMin)
	require.Equal(t, math.LegacyNewDecWithPrec(consensus.InflationMax, 2), updatedMintGenState.Params.InflationMax)
	require.Equal(t, consensus.BlocksPerYear, updatedMintGenState.Params.BlocksPerYear)

	//Verify distribution genesis state
	var updatedDistrGenState distrtypes.GenesisState
	err = appCodec.UnmarshalJSON(appState[distrtypes.ModuleName], &updatedDistrGenState)
	require.NoError(t, err)

	require.Equal(t, math.LegacyNewDecWithPrec(consensus.CommunityTax, 2), updatedDistrGenState.Params.CommunityTax)

	//verify gov genesis state
	var updatedGovGenState govv1types.GenesisState
	err = appCodec.UnmarshalJSON(appState[govtypes.ModuleName], &updatedGovGenState)
	require.NoError(t, err)

	minDepositTokens := sdk.TokensFromConsensusPower(consensus.MinDepositTokens, sdk.DefaultPowerReduction)
	require.Equal(t, sdk.NewCoin(consensus.DefaultHippoDenom, minDepositTokens), updatedGovGenState.Params.MinDeposit[0])
	require.Equal(t, 1, len(updatedGovGenState.Params.MinDeposit)) // check ahp is the only registered
	require.Equal(t, consensus.MaxDepositPeriod, *updatedGovGenState.Params.MaxDepositPeriod)
	require.Equal(t, consensus.VotingPeriod, *updatedGovGenState.Params.VotingPeriod)

	//verify slashing genesis state
	var updatedSlashingGenState slashingtypes.GenesisState
	err = appCodec.UnmarshalJSON(appState[slashingtypes.ModuleName], &updatedSlashingGenState)
	require.NoError(t, err)

	require.Equal(t, consensus.SignedBlocksWindow, int(updatedSlashingGenState.Params.SignedBlocksWindow))
	require.Equal(t, math.LegacyNewDecWithPrec(consensus.MinSignedPerWindow, 2), updatedSlashingGenState.Params.MinSignedPerWindow)
	require.Equal(t, math.LegacyNewDecWithPrec(consensus.SlashFractionDoubleSign, 2), updatedSlashingGenState.Params.SlashFractionDoubleSign)
	require.Equal(t, math.LegacyNewDecWithPrec(consensus.SlashFractionDowntime*100, 4), updatedSlashingGenState.Params.SlashFractionDowntime)

}

func TestFlags(t *testing.T) {
	require.Equal(t, "overwrite", FlagOverwrite)
	require.Equal(t, "recover", FlagRecover)
	require.Equal(t, "default-denom", FlagDefaultBondDenom)
	require.Equal(t, "staking-bond-denom", FlagStakingBondDenom)
}

func TestNewPrintInfo(t *testing.T) {
	moniker := "moniker"
	chainID := "chainID"
	nodeID := "nodeID"
	genTxsDir := "genTxsDir"
	appMessage := json.RawMessage(`{"name": "John", "age": 30}`)

	printInfo := newPrintInfo(moniker, chainID, nodeID, genTxsDir, appMessage)

	require.Equal(t, moniker, printInfo.Moniker)
	require.Equal(t, chainID, printInfo.ChainID)
	require.Equal(t, nodeID, printInfo.NodeID)
	require.Equal(t, genTxsDir, printInfo.GenTxsDir)
	require.Equal(t, appMessage, printInfo.AppMessage)
}

func TestDisplayInfo(t *testing.T) {
	moniker := "moniker"
	chainID := "chainID"
	nodeID := "nodeID"
	genTxsDir := "genTxsDir"
	appMessage := json.RawMessage(`{"name": "John", "age": 30}`)

	printInfo := newPrintInfo(moniker, chainID, nodeID, genTxsDir, appMessage)
	err := displayInfo(printInfo)

	require.NoError(t, err)
}

func TestFailingOverrideGenesis(t *testing.T) {
	getParam := func(key string) (codec.JSONCodec, *cmttypes.GenesisDoc, map[string]json.RawMessage) {
		tempDir := t.TempDir()
		hippoApp := app.New(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.NewAppOptionsWithFlagHome(tempDir), app.EmptyWasmOptions)
		appGenState := hippoApp.DefaultGenesis()
		appCodec := hippoApp.AppCodec()

		genDoc := &cmttypes.GenesisDoc{}
		genDoc.ChainID = chainID
		genDoc.Validators = nil
		genDoc.InitialHeight = 0
		genDoc.ConsensusParams = &cmttypes.ConsensusParams{
			Block:     cmttypes.DefaultBlockParams(),
			Evidence:  cmttypes.DefaultEvidenceParams(),
			Validator: cmttypes.DefaultValidatorParams(),
			Version:   cmttypes.DefaultVersionParams(),
		}

		delete(appGenState, key)

		return appCodec, genDoc, appGenState
	}

	_, err := overrideGenesis(getParam("staking"))
	require.Error(t, err)

	_, err = overrideGenesis(getParam("mint"))
	require.Error(t, err)

	_, err = overrideGenesis(getParam("distribution"))
	require.Error(t, err)

	_, err = overrideGenesis(getParam("gov"))
	require.Error(t, err)

	_, err = overrideGenesis(getParam("slashing"))
	require.Error(t, err)

}

func TestFailingDisplayInfo(t *testing.T) {
	moniker := ""
	chainID := ""
	nodeID := ""
	genTxsDir := ""
	appMessage := json.RawMessage("")

	printInfo := newPrintInfo(moniker, chainID, nodeID, genTxsDir, appMessage)
	err := displayInfo(printInfo)

	require.Error(t, err)
}
