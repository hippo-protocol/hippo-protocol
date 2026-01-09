package test

import (
	"cosmossdk.io/log"
	dbm "github.com/cosmos/cosmos-db"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/hippocrat-dao/hippo-protocol/app"
	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
)

func GetApp() app.App {

	// To prevent config seal error, not using SetWalletConfig()
	config := sdk.GetConfig()
	config.SetPurpose(consensus.BIP44Purpose)
	config.SetCoinType(consensus.BIP44CoinType)
	config.SetBech32PrefixForAccount(consensus.AddrPrefix, consensus.PubkeyPrefix)
	config.SetBech32PrefixForValidator(consensus.ValidatorAddrPrefix, consensus.ValidatorPubkeyPrefix)
	config.SetBech32PrefixForConsensusNode(consensus.ConsensusNodeAddrPrefix, consensus.ConsensusNodePubkeyPrefix)

	return *app.New(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.NewAppOptionsWithFlagHome(app.DefaultNodeHome), app.EmptyWasmOptions)
}
