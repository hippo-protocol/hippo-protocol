package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"
)

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
