package consensus_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
	"github.com/stretchr/testify/require"
)

func TestSetWalletConfig(t *testing.T) {
	// Set the wallet configuration
	consensus.SetWalletConfig()

	// Get the config from the SDK
	config := sdk.GetConfig()

	// Check prefixes for different address types
	require.Equal(t, "hippo", config.GetBech32AccountAddrPrefix(), "Prefix for general account should be 'hippo'")
	require.Equal(t, "hippopub", config.GetBech32AccountPubPrefix(), "Prefix for general pubkey should be 'hippopub'")

	require.Equal(t, "hippovaloper", config.GetBech32ValidatorAddrPrefix(), "Prefix for validator account should be 'hippovaloper'")
	require.Equal(t, "hippovaloperpub", config.GetBech32ValidatorPubPrefix(), "Prefix for validator pubkey should be 'hippovaloperpub'")

	require.Equal(t, "hippovalcons", config.GetBech32ConsensusAddrPrefix(), "Prefix for consensus node account should be 'hippovalcons'")
	require.Equal(t, "hippovalconspub", config.GetBech32ConsensusPubPrefix(), "Prefix for consensus node pubkey should be 'hippovalconspub'")

	// Check BIP44 settings
	require.Equal(t, uint32(44), config.GetPurpose(), "BIP44 purpose should be 44")
	require.Equal(t, uint32(0), config.GetCoinType(), "BIP44 coin type should be 0")
}
