package consensus

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	AddrPrefix                = "hippo"
	PubkeyPrefix              = AddrPrefix + "pub"
	ValidatorAddrPrefix       = AddrPrefix + "valoper"
	ValidatorPubkeyPrefix     = AddrPrefix + "valoperpub"
	ConsensusNodeAddrPrefix   = AddrPrefix + "valcons"
	ConsensusNodePubkeyPrefix = AddrPrefix + "valconspub"
)

const (
	BIP44Purpose  = 44
	BIP44CoinType = 0
)

// Set prefix of address and pubkey for general, validator and consensus node.
//
// Set BIP44 path purpose and coin type.
func SetWalletConfig() {
	config := sdk.GetConfig()
	config.SetPurpose(BIP44Purpose)
	config.SetCoinType(BIP44CoinType)
	config.SetBech32PrefixForAccount(AddrPrefix, PubkeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddrPrefix, ValidatorPubkeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsensusNodeAddrPrefix, ConsensusNodePubkeyPrefix)
	config.Seal()
}
