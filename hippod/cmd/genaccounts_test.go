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

	// Marshal the genesis doc
	genDocBytes, err := json.Marshal(genesisDoc)
	require.NoError(t, err)

	err = os.WriteFile(genesisFile, genDocBytes, 0644)
	require.NoError(t, err)

	// Create the AddGenesisAccountCmd
	cmd := AddGenesisAccountCmd(tempDir)

	// Set command flags and args
	addr := "cosmos1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	coins := "100stake,50token"
	cmd.SetArgs([]string{addr, coins, "--home", tempDir})

	// Execute the command
	err = cmd.Execute()
	require.Error(t, err, "Expected error due to missing keyring configuration")

	// Test with invalid coins format
	invalidCoins := "100stake,50"
	cmd.SetArgs([]string{addr, invalidCoins, "--home", tempDir})
	err = cmd.Execute()
	require.Error(t, err, "Expected error due to invalid coin format")

	// Test with key name (keyring)
	cmdKey := AddGenesisAccountCmd(tempDir)
	cmdKey.SetArgs([]string{"testkey", coins, "--home", tempDir, "--keyring-backend", "test"})

	// Mock stdin for keyring input
	r, w, err := os.Pipe()
	require.NoError(t, err)
	cmdKey.SetIn(r)
	_, err = w.WriteString("testpassword\n")
	require.NoError(t, err)
	w.Close()

	// Simulate keyring error by adding the keyring creation step
	err = cmdKey.Execute()
	require.Error(t, err, "Expected error due to failed keyring setup")

	// Test with valid key name and valid arguments
	// Mock keyring (or use a mock in the test framework)
	cmdKey.SetArgs([]string{"testkey", coins, "--home", tempDir, "--keyring-backend", "os"})

	// Mock stdin for keyring input
	r, w, err = os.Pipe()
	require.NoError(t, err)
	cmdKey.SetIn(r)
	_, err = w.WriteString("testpassword\n")
	require.NoError(t, err)
	w.Close()

	err = cmdKey.Execute()
	require.NoError(t, err, "Expected no error for valid keyring input")

	// Test with vesting params (successful vesting creation)
	cmdVesting := AddGenesisAccountCmd(tempDir)
	vestingCoins := "100stake,50token"
	cmdVesting.SetArgs([]string{addr, vestingCoins, "--home", tempDir, "--vesting-start-time", "1000", "--vesting-end-time", "2000", "--vesting-amount", "200stake"})

	err = cmdVesting.Execute()
	require.NoError(t, err, "Expected no error for valid vesting parameters")

	// Test invalid vesting parameters
	cmdInvalidVesting := AddGenesisAccountCmd(tempDir)
	cmdInvalidVesting.SetArgs([]string{addr, vestingCoins, "--home", tempDir, "--vesting-start-time", "1000", "--vesting-amount", "200stake"})

	err = cmdInvalidVesting.Execute()
	require.Error(t, err, "Expected error for missing vesting end time")

	// Test append mode
	cmdAppend := AddGenesisAccountCmd(tempDir)
	cmdAppend.SetArgs([]string{addr, vestingCoins, "--home", tempDir, "--append"})

	err = cmdAppend.Execute()
	require.NoError(t, err, "Expected no error when appending to existing account")

	// Test for valid successful account creation
	cmdCreate := AddGenesisAccountCmd(tempDir)
	cmdCreate.SetArgs([]string{addr, "100stake", "--home", tempDir})

	err = cmdCreate.Execute()
	require.NoError(t, err, "Expected no error for valid account creation")

	// Test for non-existent address and invalid keyring arguments
	nonExistentAddr := "cosmos1nonexistentaddressxxxxxxxxxxxxx"
	cmd.SetArgs([]string{nonExistentAddr, coins, "--home", tempDir})
	err = cmd.Execute()
	require.Error(t, err, "Expected error for non-existent address")

	// Simulate failed JSON marshal for genesis state
	// Create an invalid genesis file to test error handling
	invalidGenesisFile := filepath.Join(tempDir, "config", "genesis_invalid.json")
	err = os.WriteFile(invalidGenesisFile, []byte("invalid json"), 0644)
	require.NoError(t, err)

	// Test with invalid genesis file
	cmd.SetArgs([]string{addr, coins, "--home", tempDir, "--vesting-start-time", "1000"})
	err = cmd.Execute()
	require.Error(t, err, "Expected error for invalid genesis file")
}
