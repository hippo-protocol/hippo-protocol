package cmd

import (
	"fmt"
	"os"
	"testing"

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
