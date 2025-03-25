package cmd_test

import (
	"encoding/json"
	"os"
	"path/filepath"

	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hippocrat-dao/hippo-protocol/hippod/cmd"
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

func TestAddGenesisAccountCmd_FullCoverage(t *testing.T) {
	home := t.TempDir()
	createMinimalGenesisFile(t, home)

	tests := []struct {
		name        string
		args        []string
		expectErr   bool
		errContains string
	}{
		{
			name:      "Invalid coins",
			args:      []string{"cosmos1xx", "invalidcoin", "--home", home},
			expectErr: true,
		},
		{
			name:        "Vesting amount > total",
			args:        []string{"cosmos1xx", "50stake", "--home", home, "--vesting-amount", "100stake", "--vesting-end-time", "10000"},
			expectErr:   true,
			errContains: "vesting amount cannot be greater",
		},
		{
			name: "Continuous vesting",
			args: []string{"cosmos1abc", "100stake", "--home", home, "--vesting-amount", "50stake", "--vesting-start-time", "1000", "--vesting-end-time", "2000"},
		},
		{
			name: "Delayed vesting",
			args: []string{"cosmos1def", "100stake", "--home", home, "--vesting-amount", "50stake", "--vesting-end-time", "3000"},
		},
		{
			name:        "Missing vesting params",
			args:        []string{"cosmos1vest", "100stake", "--home", home, "--vesting-amount", "50stake"},
			expectErr:   true,
			errContains: "invalid vesting parameters",
		},
		{
			name: "Successful account",
			args: []string{"cosmos1success", "100stake", "--home", home},
		},
		{
			name: "Append to existing account",
			args: []string{"cosmos1success", "50stake", "--home", home, "--append"},
		},
		{
			name:        "Duplicate without append",
			args:        []string{"cosmos1success", "50stake", "--home", home},
			expectErr:   true,
			errContains: "cannot add account at existing address",
		},
		{
			name:        "Invalid vesting coins",
			args:        []string{"cosmos1failvest", "100stake", "--home", home, "--vesting-amount", "badcoin"},
			expectErr:   true,
			errContains: "failed to parse vesting amount",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd.AddGenesisAccountCmd(home)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if tt.expectErr {
				require.Error(t, err)
				if tt.errContains != "" {
					require.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
