package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	cfg "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/libs/cli"
	tmrand "github.com/cometbft/cometbft/libs/rand"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/go-bip39"
	"github.com/pkg/errors"

	genttypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	tmjson "github.com/cometbft/cometbft/libs/json"
	"github.com/cosmos/cosmos-sdk/codec"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/hippocrat-dao/hippo-protocol/app"
	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
	"github.com/spf13/cobra"
)

const (
	// FlagOverwrite defines a flag to overwrite an existing genesis JSON file.
	FlagOverwrite = "overwrite"

	// FlagSeed defines a flag to initialize the private validator key from a specific seed.
	FlagRecover = "recover"

	// FlagDefaultBondDenom defines the default denom to use in the genesis file.
	FlagDefaultBondDenom = "default-denom"

	// FlagStakingBondDenom defines a flag to specify the staking token in the genesis file.
	FlagStakingBondDenom = "staking-bond-denom"
	blockTimeSec         = consensus.BlockTimeSec    // 5s of timeout_commit + 1s
	unbondingPeriod      = consensus.UnbondingPeriod // three weeks
)

type printInfo struct {
	Moniker    string          `json:"moniker" yaml:"moniker"`
	ChainID    string          `json:"chain_id" yaml:"chain_id"`
	NodeID     string          `json:"node_id" yaml:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir" yaml:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message" yaml:"app_message"`
}

// InitCmd wraps the genutilcli.InitCmd to inject specific parameters for Hippo.
// It reads the default genesis.json, modifies and exports it.
// Reference: https://github.com/public-awesome/stargaze/blob/b92bf9847559b9c7f4ac08576d056d3d00efe12c/cmd/starsd/cmd/init.go
func InitCmd(mbm module.BasicManager, defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [moniker]",
		Short: "Initialize private validator, p2p, genesis, and application configuration files",
		Long:  `Initialize validators's and node's configuration files.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config
			config.SetRoot(clientCtx.HomeDir)

			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			switch {
			case chainID != "":
			case clientCtx.ChainID != "":
				chainID = clientCtx.ChainID
			default:
				chainID = fmt.Sprintf("test-chain-%v", tmrand.Str(6))
			}

			// Get bip39 mnemonic
			var mnemonic string
			recover, _ := cmd.Flags().GetBool(FlagRecover)
			if recover {
				inBuf := bufio.NewReader(cmd.InOrStdin())
				value, err := input.GetString("Enter your bip39 mnemonic", inBuf)
				if err != nil {
					return err
				}

				mnemonic = value
				if !bip39.IsMnemonicValid(mnemonic) {
					return errors.New("invalid mnemonic")
				}
			}

			// Get initial height
			initHeight, _ := cmd.Flags().GetInt64(flags.FlagInitHeight)
			if initHeight < 1 {
				initHeight = 1
			}

			nodeID, _, err := genutil.InitializeNodeValidatorFilesFromMnemonic(config, mnemonic)
			if err != nil {
				return err
			}

			config.Moniker = args[0]

			genFile := config.GenesisFile()
			overwrite, _ := cmd.Flags().GetBool(FlagOverwrite)
			defaultDenom, _ := cmd.Flags().GetString(FlagDefaultBondDenom)

			// use os.Stat to check if the file exists
			_, err = os.Stat(genFile)
			if !overwrite && !os.IsNotExist(err) {
				return fmt.Errorf("genesis.json file already exists: %v", genFile)
			}

			// Overwrites the SDK default denom for side-effects
			if defaultDenom != "" {
				sdk.DefaultBondDenom = defaultDenom
			}
			appGenState := mbm.DefaultGenesis(cdc)
			genDoc := &types.GenesisDoc{}
			if _, err := os.Stat(genFile); err != nil {
				if !os.IsNotExist(err) {
					return err
				}
			} else {
				genDoc, err = types.GenesisDocFromFile(genFile)
				if err != nil {
					return errors.Wrap(err, "Failed to read genesis doc from file")
				}
			}

			genDoc.ChainID = chainID
			genDoc.Validators = nil
			genDoc.InitialHeight = initHeight
			genDoc.ConsensusParams = &types.ConsensusParams{
				Block:     types.DefaultBlockParams(),
				Evidence:  types.DefaultEvidenceParams(),
				Validator: types.DefaultValidatorParams(),
				Version:   types.DefaultVersionParams(),
			}

			appState, err := overrideGenesis(cdc, genDoc, appGenState)
			if err != nil {
				return errors.Wrap(err, "Failed to marshal default genesis state")
			}

			genDoc.AppState = appState

			// v0.50.12 requires AppGenesis to pass to generate Genesis File
			// reference: https://github.com/cosmos/cosmos-sdk/blob/v0.50.12/x/genutil/client/cli/init.go#L144-L163
			appGenesis := &genttypes.AppGenesis{}
			appGenesis.AppName = app.Name
			appGenesis.AppVersion = app.Version
			appGenesis.ChainID = chainID
			appGenesis.AppState = appState
			appGenesis.InitialHeight = initHeight
			appGenesis.Consensus = &genttypes.ConsensusGenesis{
				Validators: nil,
				Params: &types.ConsensusParams{
					Block: types.BlockParams{
						MaxBytes: genDoc.ConsensusParams.Block.MaxBytes,
						MaxGas:   genDoc.ConsensusParams.Block.MaxGas,
					},
					Evidence: types.EvidenceParams{
						MaxAgeNumBlocks: genDoc.ConsensusParams.Evidence.MaxAgeNumBlocks,
						MaxAgeDuration:  genDoc.ConsensusParams.Evidence.MaxAgeDuration,
						MaxBytes:        genDoc.ConsensusParams.Evidence.MaxBytes,
					},
					Validator: types.ValidatorParams{
						PubKeyTypes: genDoc.ConsensusParams.Validator.PubKeyTypes,
					},
					Version: types.VersionParams{
						App: genDoc.ConsensusParams.Version.App,
					},
				},
			}

			if err = genutil.ExportGenesisFile(appGenesis, genFile); err != nil {
				return errors.Wrap(err, "Failed to export genesis file")
			}

			toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", appState)
			cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)
			return displayInfo(toPrint)
		},
	}

	cmd.Flags().String(cli.HomeFlag, defaultNodeHome, "node's home directory")
	cmd.Flags().BoolP(FlagOverwrite, "o", false, "overwrite the genesis.json file")
	cmd.Flags().Bool(FlagRecover, false, "provide seed phrase to recover existing key instead of creating")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(FlagDefaultBondDenom, "", "genesis file default denomination, if left blank default value is 'stake'")
	cmd.Flags().Int64(flags.FlagInitHeight, 1, "specify the initial block height at genesis")

	return cmd
}

// overrideGenesis overrides some parameters in the genesis doc to the hippo-specific values.
func overrideGenesis(cdc codec.JSONCodec, genDoc *types.GenesisDoc, appState map[string]json.RawMessage) (json.RawMessage, error) {
	genDoc.ConsensusParams.Block.MaxBytes = consensus.MaxBlockSize // 4MB
	genDoc.ConsensusParams.Block.MaxGas = consensus.MaxBlockGas    // 100 milion

	var bankGenState banktypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[banktypes.ModuleName], &bankGenState); err != nil {
		return nil, err
	}

	hpMetadata := banktypes.Metadata{
		Name:        "Hippo", // Often the full name
		Symbol:      "HP",    // The commonly used ticker symbol
		Description: "The native staking token of the Hippo Protocol.",
		DenomUnits: []*banktypes.DenomUnit{ // Note: DenomUnits is a slice of *pointers* to DenomUnit
			{
				Denom:    "ahp", // base unit
				Exponent: 0,
			},
		},
		Base:    "ahp", // The base denomination (exponent 0)
		Display: "ahp", // The denomination users typically see
	}

	bankGenState.DenomMetadata = append(bankGenState.DenomMetadata, hpMetadata)

	appState[banktypes.ModuleName] = cdc.MustMarshalJSON(&bankGenState)

	var stakingGenState stakingtypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[stakingtypes.ModuleName], &stakingGenState); err != nil {
		return nil, err
	}

	stakingGenState.Params.UnbondingTime = unbondingPeriod
	stakingGenState.Params.MaxValidators = consensus.MaxValidators
	stakingGenState.Params.BondDenom = consensus.DefaultHippoDenom
	stakingGenState.Params.MinCommissionRate = math.LegacyNewDecWithPrec(consensus.MinCommissionRate, 2)
	appState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingGenState)

	var mintGenState minttypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[minttypes.ModuleName], &mintGenState); err != nil {
		return nil, err
	}
	mintGenState.Minter = minttypes.InitialMinter(math.LegacyNewDecWithPrec(consensus.Minter, 2)) // 25% inflation
	mintGenState.Params.MintDenom = consensus.DefaultHippoDenom
	mintGenState.Params.InflationRateChange = math.LegacyNewDecWithPrec(consensus.InflationRateChange, 2) // 25%
	mintGenState.Params.InflationMin = math.LegacyNewDecWithPrec(consensus.InflationMin, 2)               // 0%
	mintGenState.Params.InflationMax = math.LegacyNewDecWithPrec(consensus.InflationMax, 2)               // 25%
	mintGenState.Params.BlocksPerYear = consensus.BlocksPerYear
	appState[minttypes.ModuleName] = cdc.MustMarshalJSON(&mintGenState)

	var distrGenState distrtypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[distrtypes.ModuleName], &distrGenState); err != nil {
		return nil, err
	}
	distrGenState.Params.CommunityTax = math.LegacyNewDecWithPrec(consensus.CommunityTax, 2)
	appState[distrtypes.ModuleName] = cdc.MustMarshalJSON(&distrGenState)

	var govGenState govv1types.GenesisState
	if err := cdc.UnmarshalJSON(appState[govtypes.ModuleName], &govGenState); err != nil {
		return nil, err
	}
	minDepositTokens := sdk.TokensFromConsensusPower(consensus.MinDepositTokens, sdk.DefaultPowerReduction) // 50,000 HP
	govGenState.Params.MinDeposit = sdk.Coins{sdk.NewCoin(consensus.DefaultHippoDenom, minDepositTokens)}
	expeditedMinDepositTokens := sdk.TokensFromConsensusPower(consensus.ExpeditedMinDeposit, sdk.DefaultPowerReduction) // 100,000 HP
	govGenState.Params.ExpeditedMinDeposit = sdk.Coins{sdk.NewCoin(consensus.DefaultHippoDenom, expeditedMinDepositTokens)}
	maxDepositPeriod := consensus.MaxDepositPeriod // 14 days
	govGenState.Params.MaxDepositPeriod = &maxDepositPeriod
	votingPeriod := consensus.VotingPeriod
	govGenState.Params.VotingPeriod = &votingPeriod
	expeditedVotingPeriod := consensus.ExpeditedVotingPeriod
	govGenState.Params.ExpeditedVotingPeriod = &expeditedVotingPeriod
	appState[govtypes.ModuleName] = cdc.MustMarshalJSON(&govGenState)

	var slashingGenState slashingtypes.GenesisState
	if err := cdc.UnmarshalJSON(appState[slashingtypes.ModuleName], &slashingGenState); err != nil {
		return nil, err
	}
	slashingGenState.Params.SignedBlocksWindow = consensus.SignedBlocksWindow
	slashingGenState.Params.MinSignedPerWindow = math.LegacyNewDecWithPrec(consensus.MinSignedPerWindow, 2)
	slashingGenState.Params.SlashFractionDoubleSign = math.LegacyNewDecWithPrec(consensus.SlashFractionDoubleSign, 2) // 5%
	slashingGenState.Params.SlashFractionDowntime = math.LegacyNewDecWithPrec(consensus.SlashFractionDowntime*100, 4) // 0.01%
	appState[slashingtypes.ModuleName] = cdc.MustMarshalJSON(&slashingGenState)

	// MaxAgeDuration and MaxAgeNumBlocks values should be longer than unbonding period.
	// Otherwise it may allow malicious validators to escape penalties.
	// https://github.com/advisories/GHSA-555p-m4v6-cqxv
	genDoc.ConsensusParams.Evidence.MaxAgeDuration = consensus.MaxAgeDuration          // 30 days
	genDoc.ConsensusParams.Evidence.MaxAgeNumBlocks = int64(consensus.MaxAgeNumBlocks) // 30 days

	return tmjson.Marshal(appState)
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string, appMessage json.RawMessage) printInfo {
	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

func displayInfo(info printInfo) error {
	out, err := json.MarshalIndent(info, "", " ")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stderr, "%s\n", sdk.MustSortJSON(out))

	return err
}
