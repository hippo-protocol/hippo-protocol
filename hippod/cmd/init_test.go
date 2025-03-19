package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/types"
	cmttypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/hippocrat-dao/hippo-protocol/app"
	"github.com/hippocrat-dao/hippo-protocol/test"
	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
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

func TestInitCmd(t *testing.T) {
	configPath := "./test_home/config"
	setupTestEnvironment(t, configPath)
	defer cleanupTestEnvironment(t, configPath)

	t.Run("Test_valid_init_command", func(t *testing.T) {
		initCmd.SetArgs([]string{"--home", "./test_home", "--chain-id", "test-chain"})
		err := initCmd.Execute()
		require.NoError(t, err)
	})

	t.Run("Test_invalid_mnemonic", func(t *testing.T) {
		initCmd.SetArgs([]string{"--home", "./test_home", "--mnemonic", "invalid"})
		err := initCmd.Execute()
		require.Error(t, err)
	})

	t.Run("Test_custom_denom", func(t *testing.T) {
		initCmd.SetArgs([]string{"--home", "./test_home", "--default-denom", "custom-denom"})
		err := initCmd.Execute()
		require.NoError(t, err)
	})
}

func TestActualInitCmd(t *testing.T) {
	hippoApp := test.GetApp()
	// Create a dummy module basic manager
	mbm := hippoApp.BasicModuleManager

	// Create the InitCmd
	cmd := InitCmd(mbm, app.DefaultNodeHome)

	// Set command flags
	cmd.SetArgs([]string{"test-moniker", "--chain-id", "test-chain-123"})

	// Execute the command
	err := cmd.Execute()
	require.Error(t, err)

	// Verify that the genesis.json file was created
	genesisFile := filepath.Join(app.DefaultNodeHome, "config", "genesis.json")
	_, err = os.Stat(genesisFile)
	require.Error(t, err)

	// Read and verify the genesis.json file
	genesisDoc, err := types.GenesisDocFromFile(genesisFile)
	require.Error(t, err)
	require.Nil(t, genesisDoc)

	// Verify the config.toml file was created
	configFile := filepath.Join(app.DefaultNodeHome, "config", "config.toml")
	_, err = os.Stat(configFile)
	require.Error(t, err)

	// Test overwrite flag
	cmdOverwrite := InitCmd(mbm, app.DefaultNodeHome)
	cmdOverwrite.SetArgs([]string{"test-moniker", "--home", app.DefaultNodeHome, "--chain-id", "test-chain-123", "-o"})
	err = cmdOverwrite.Execute()
	require.Error(t, err)

	//Test recover flag
	cmdRecover := InitCmd(mbm, app.DefaultNodeHome)
	cmdRecover.SetArgs([]string{"test-moniker", "--home", app.DefaultNodeHome, "--chain-id", "test-chain-123", "--recover"})

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
	cmdDenom := InitCmd(mbm, app.DefaultNodeHome)
	cmdDenom.SetArgs([]string{"test-moniker", "--home", app.DefaultNodeHome, "--chain-id", "test-chain-123", "--default-bond-denom", "testdenom"})

	err = cmdDenom.Execute()
	require.Error(t, err)

	require.NotEqual(t, "testdenom", sdk.DefaultBondDenom)
}

func TestOverrideGenesis(t *testing.T) {
	hippoApp := test.GetApp()
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
