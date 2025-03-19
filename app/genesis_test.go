package app

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenesisState(t *testing.T) {
	// Create sample raw genesis messages for testing
	module1Genesis := json.RawMessage(`{"module1_field": "value1"}`)
	module2Genesis := json.RawMessage(`{"module2_field": "wrong_value"}`)

	// Initialize a new GenesisState map
	genesisState := GenesisState{
		"module1": module1Genesis,
		"module2": module2Genesis,
	}

	// Test: Check if module1's genesis state is correctly retrieved
	t.Run("Test retrieving module1 genesis", func(t *testing.T) {
		result, exists := genesisState["module1"]
		assert.True(t, exists, "Expected module1 genesis state to exist")
		assert.JSONEq(t, string(module1Genesis), string(result), "Module1 genesis state mismatch")
	})

	// Test: Check if module2's genesis state is correctly retrieved
	t.Run("Test retrieving module2 genesis", func(t *testing.T) {
		result, exists := genesisState["module2"]
		assert.True(t, exists, "Expected module2 genesis state to exist")
		assert.JSONEq(t, string(module2Genesis), string(result), "Module2 genesis state mismatch")
	})

	// Test: Check if a non-existent module returns false
	t.Run("Test retrieving non-existent module", func(t *testing.T) {
		_, exists := genesisState["module3"]
		assert.False(t, exists, "Expected non-existent module to not exist")
	})
}
