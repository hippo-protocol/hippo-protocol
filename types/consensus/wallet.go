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

// Set prefix of address and pubkey for general, validator and consensus node
func SetWalletPrefix() {
	config := sdk.GetConfig()
	config.SetPurpose(44)
	config.SetCoinType(0)
	config.SetBech32PrefixForAccount(AddrPrefix, PubkeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddrPrefix, ValidatorPubkeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsensusNodeAddrPrefix, ConsensusNodePubkeyPrefix)
	config.Seal()
}
