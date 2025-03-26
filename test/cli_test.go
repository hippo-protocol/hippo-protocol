package test

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Test struct {
	command  []string
	expect   string
	errorMsg string
}

const delegator_address = "hippo18cdh0u0tg6yc623mna3glf3fmx4qq5y0627tw2"        // Address of the delegator
const validator_address = "hippovaloper18cdh0u0tg6yc623mna3glf3fmx4qq5y0tss90e" // Address of the validator
const target_address = "hippo1mj5e9kpths3x5qsxarax9c50dadumyj8rqxq95"           // any bech32 hippo address that is tx target(used for sending, ...etc)
const passphrase = "12345678"                                                   // used when sending tx

const path = "../hippod/main.go"

func testQuery(t *testing.T, tests []Test) {
	for _, test := range tests {
		cmd := exec.Command("go", append([]string{"run", path}, test.command...)...)
		out, err := cmd.CombinedOutput()
		assert.NoError(t, err, "cli should not return an error")
		assert.Contains(t, string(out), test.expect, test.errorMsg)
	}
}

func TestAuth(t *testing.T) {
	tests := []Test{
		{command: []string{"query", "auth", "accounts"}, expect: "accounts:", errorMsg: "all accounts should be in the output"},
		{command: []string{"query", "auth", "params"}, expect: "max_memo_characters", errorMsg: "max_memo_characters should be in auth parameters"},
	}

	testQuery(t, tests)
}

func TestBank(t *testing.T) {
	tests := []Test{
		{command: []string{"query", "bank", "balances", delegator_address}, expect: "balances", errorMsg: "balances should be in the output"},
		// {command: []string{"query", "bank", "denom-metadata", "ahp"}, expect: "ahp", errorMsg: "metadata should be in the output"}, // Fail, currently metadata do not exists
		{command: []string{"query", "bank", "total"}, expect: "supply", errorMsg: "supply data for ahp should be in the output"},
	}

	testQuery(t, tests)
}

func TestDistribution(t *testing.T) {
	tests := []Test{
		{command: []string{"query", "distribution", "community-pool"}, expect: "ahp", errorMsg: "community pool balance should be in the output"},
		{command: []string{"query", "distribution", "params"}, expect: "community_tax", errorMsg: "community_tax should be in the distribution params"},
		{command: []string{"query", "distribution", "rewards", delegator_address}, expect: "reward:", errorMsg: "delegator rewards should be in the output"},
	}

	testQuery(t, tests)
}

func TestGov(t *testing.T) {
	tests := []Test{
		{command: []string{"query", "gov", "proposals"}, expect: "pagination", errorMsg: "all proposals should be in the output"},
		{command: []string{"query", "gov", "params"}, expect: "min_deposit", errorMsg: "min_deposit should be in gov parameters"},
	}

	testQuery(t, tests)
}

func TestMint(t *testing.T) {
	tests := []Test{
		{command: []string{"query", "mint", "inflation"}, expect: "inflation:", errorMsg: "inflation should be calculated"},
		{command: []string{"query", "mint", "params"}, expect: "blocks_per_year", errorMsg: "blocks_per_year should be in mint params"},
	}

	testQuery(t, tests)
}

func TestStaking(t *testing.T) {
	tests := []Test{
		{command: []string{"query", "staking", "delegations", delegator_address}, expect: "delegator_address", errorMsg: "delegations made by address should be in the output"},
		{command: []string{"query", "staking", "validators"}, expect: "hippovaloper", errorMsg: "validator address should be in the output"},
		{command: []string{"query", "staking", "historical-info", "10"}, expect: "header", errorMsg: "staking history in specific block height should be printed correctly"},
		{command: []string{"query", "staking", "params"}, expect: "min_commission_rate:", errorMsg: "min_commission_rate should be inside staking params"},
	}

	testQuery(t, tests)
}

func TestUpgrade(t *testing.T) {
	tests := []Test{
		{command: []string{"query", "upgrade", "module_versions"}, expect: "name: staking", errorMsg: "staking module should be in the output"},
		{command: []string{"query", "upgrade", "plan"}, expect: "", errorMsg: "upgrade plan should be queried correctly"}, // just checking it works
	}

	testQuery(t, tests)
}

func TestTx(t *testing.T) {
	tests := []Test{
		{command: []string{"tx", "bank", "send", delegator_address, target_address, "1000000000000000000ahp", "--fees=1000000000000000000ahp", "-y"}, expect: "txhash", errorMsg: "txhash should be in the output"},
		{command: []string{"tx", "staking", "delegate", validator_address, "1000000000000000000ahp", "--fees=1000000000000000000ahp", fmt.Sprintf("--from=%s", delegator_address), "-y"}, expect: "txhash", errorMsg: "txhash should be in the output"},
	}

	for _, test := range tests {
		cmd := exec.Command("go", append([]string{"run", path}, test.command...)...)

		// Create a pipe for stdin.
		stdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatalf("Failed to create stdin pipe: %v", err)
		}

		// Write the input to stdin.
		_, err = stdin.Write([]byte(passphrase))
		if err != nil {
			t.Fatalf("Failed to write to stdin: %v", err)
		}

		// Close stdin to signal EOF.
		err = stdin.Close()
		if err != nil {
			t.Fatalf("Failed to close stdin: %v", err)
		}

		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Command failed: %v, output: %s", err, out)
		}
		assert.NoError(t, err, "Sending Tx should work")
		assert.Contains(t, string(out), test.expect, test.errorMsg)
	}

}
