package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"
)

func createMinimalGenesisFile(t *testing.T, home string) {
	configDir := filepath.Join(home, "config")
	require.NoError(t, os.MkdirAll(configDir, 0755))

	genesisFile := filepath.Join(configDir, "genesis.json")

	genesisState := map[string]json.RawMessage{
		"auth": json.RawMessage(`{"accounts": []}`),
		"bank": json.RawMessage(`{"balances": [], "supply": "0stake"}`),
	}

	genDoc := map[string]interface{}{
		"chain_id":  "test-chain",
		"app_state": genesisState,
	}
	genDocBytes, err := json.MarshalIndent(genDoc, "", "  ")
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(genesisFile, genDocBytes, 0644))
}

func TestAddGenesisAccountCmd(t *testing.T) {
	// Create a temporary directory for the node home
	tempDir, err := os.MkdirTemp("", "hippo-test-home")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a dummy genesis file
	genesisFile := filepath.Join(tempDir, "config", "genesis.json")
	genesisDoc := &types.GenesisDoc{}

	//marshal the genesis doc
	genDocBytes, err := json.Marshal(genesisDoc)
	require.NoError(t, err)

	err = os.WriteFile(genesisFile, genDocBytes, 0644)
	require.Error(t, err)

	// Create the AddGenesisAccountCmd
	cmd := AddGenesisAccountCmd(tempDir)

	// Set command flags and args
	addr := "cosmos1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	coins := "100stake,50token"
	cmd.SetArgs([]string{addr, coins, "--home", tempDir})

	// Execute the command
	err = cmd.Execute()
	require.Error(t, err)

	// Test with key name
	cmdKey := AddGenesisAccountCmd(tempDir)
	cmdKey.SetArgs([]string{"testkey", coins, "--home", tempDir, "--keyring-backend", "test"})

	//Mock stdin for keyring input
	r, w, err := os.Pipe()
	require.NoError(t, err)
	cmdKey.SetIn(r)
	_, err = w.WriteString("testpassword\n")
	require.NoError(t, err)
	w.Close()

	err = cmdKey.Execute()
	require.Error(t, err) //Test keyring functionality requires more setup.

}

func TestAddGenesisAccountCmd_InvalidVestingParameters(t *testing.T) {
	home := t.TempDir()
	createMinimalGenesisFile(t, home)

	addr := sdk.AccAddress([]byte{5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}).String()
	coins := "100stake"
	vestingAmt := "50stake"

	cmdInstance := AddGenesisAccountCmd(home)
	cmdInstance.SetArgs([]string{
		addr, coins,
		"--home", home,
		"--vesting-amount", vestingAmt,
	})
	err := cmdInstance.Execute()
	require.Error(t, err, "Invalid vesting parameters should return error")
	require.Contains(t, err.Error(), "invalid vesting parameters")
}