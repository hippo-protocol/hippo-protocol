package test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Test struct {
	command  []string
	expect   string
	errorMsg string
}

const key_delegator_address = "HIPPO_DELEGATOR_ADDRESS"
const key_validator_address = "HIPPO_VALIDATOR_ADDRESS"
const target_address = "hippo1mj5e9kpths3x5qsxarax9c50dadumyj8rqxq95" // any bech32 hippo address that is tx target(used for sending, ...etc)
const passphrase = "password"                                         // used when sending tx

const path = "../../hippod/main.go"

func getDelegatorAddress() (string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(`echo -e "%s" | go run ../../hippod/main.go keys show alice --keyring-backend file | awk '/address:/ {print $3}'`, passphrase))
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getValidatorAddress() (string, error) {
	cmd := exec.Command("go", "run", path, "query", "staking", "validators")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`operator_address:\s*(\S+)`)
	match := re.FindStringSubmatch(string(out))

	if len(match) > 1 {
		operatorAddress := match[1]
		return operatorAddress, nil
	} else {
		return "", fmt.Errorf("Validator address not found")
	}
}

// Our token amounts exceed int64 range, so convert string to big.Int and compare
func compareAmount(amount1 string, amount2 string) int {
	bn1, _ := new(big.Int).SetString(amount1, 10)
	bn2, _ := new(big.Int).SetString(amount2, 10)
	return bn1.Cmp(bn2)
}

func testQuery(t *testing.T, tests []Test) {
	for _, test := range tests {
		cmd := exec.Command("go", append([]string{"run", path}, test.command...)...)
		out, err := cmd.CombinedOutput()
		assert.NoError(t, err, "cli should not return an error")
		assert.Contains(t, string(out), test.expect, test.errorMsg)
	}
}

func testTx(t *testing.T, command []string) string {
	cmd := exec.Command("go", append([]string{"run", path}, command...)...)

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
	assert.Contains(t, string(out), "txhash", "txhash should be in the output")

	return string(out)
}

func TestMain(m *testing.M) {
	delegator_address, err := getDelegatorAddress()
	if err != nil {
		os.Exit(1)
	}
	validator_address, err := getValidatorAddress()
	if err != nil {
		os.Exit(1)
	}

	// setup delegator & validator address using file keyring backend for other tests
	os.Setenv(key_delegator_address, delegator_address)
	os.Setenv(key_validator_address, validator_address)

	exitCode := m.Run() // Running Tests

	// cleanup environment variables
	os.Unsetenv(key_delegator_address)
	os.Unsetenv(key_validator_address)

	os.Exit(exitCode)
}

func TestAuth(t *testing.T) {
	tests := []Test{
		{command: []string{"query", "auth", "accounts"}, expect: "accounts:", errorMsg: "all accounts should be in the output"},
		{command: []string{"query", "auth", "params"}, expect: "max_memo_characters", errorMsg: "max_memo_characters should be in auth parameters"},
	}

	testQuery(t, tests)
}

func TestBank(t *testing.T) {
	delegator_address := os.Getenv(key_delegator_address)
	tests := []Test{
		{command: []string{"query", "bank", "balances", delegator_address}, expect: "balances", errorMsg: "balances should be in the output"},
		// {command: []string{"query", "bank", "denom-metadata", "ahp"}, expect: "ahp", errorMsg: "metadata should be in the output"}, // Fail, currently metadata do not exists
		{command: []string{"query", "bank", "total"}, expect: "supply", errorMsg: "supply data for ahp should be in the output"},
	}

	testQuery(t, tests)
}

func TestDistribution(t *testing.T) {
	delegator_address := os.Getenv(key_delegator_address)
	validator_address := os.Getenv(key_validator_address)

	tests := []Test{
		{command: []string{"query", "distribution", "community-pool"}, expect: "ahp", errorMsg: "community pool balance should be in the output"},
		{command: []string{"query", "distribution", "params"}, expect: "community_tax", errorMsg: "community_tax should be in the distribution params"},
		{command: []string{"query", "distribution", "rewards", delegator_address}, expect: "reward:", errorMsg: "delegator rewards should be in the output"},
		{command: []string{"query", "distribution", "commission", validator_address}, expect: "commission:", errorMsg: "validator commission should be in the output"},
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
	delegator_address := os.Getenv(key_delegator_address)
	tests := []Test{
		{command: []string{"query", "staking", "delegations", delegator_address}, expect: "delegator_address", errorMsg: "delegations made by address should be in the output"},
		{command: []string{"query", "staking", "validators"}, expect: "hippovaloper", errorMsg: "validator address should be in the output"},
		{command: []string{"query", "staking", "historical-info", "1"}, expect: "header", errorMsg: "staking history in specific block height should be printed correctly"},
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

func TestSending(t *testing.T) {
	delegator_address := os.Getenv(key_delegator_address)

	// send hp from delegator_address to target_address
	testTx(t, []string{"tx", "bank", "send", delegator_address, target_address, "1000000000000000000ahp", "--fees=1000000000000000000ahp", "-y", "--keyring-backend=file"})

	// sometimes the results are not updated immediately, so wait for a new block
	time.Sleep(6 * time.Second)

	// check target_address balance
	cmd := exec.Command("go", "run", path, "query", "bank", "balances", target_address)
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, "balance should be queried correctly")
	re := regexp.MustCompile(`amount:\s*"(\d+)"\s*denom: ahp`)
	match := re.FindStringSubmatch(string(out))
	assert.Condition(t, func() bool { return len(match) > 1 }, "balance should be in the output")
	assert.Greater(t, match[1], "0", "balance should be greater than 0 after receiving hp")
}

func TestStakingTx(t *testing.T) {
	delegator_address := os.Getenv(key_delegator_address)
	validator_address := os.Getenv(key_validator_address)

	// Delegate tokens to the validator
	testTx(t, []string{"tx", "staking", "delegate", validator_address, "1000000000000000000ahp", "--fees=1000000000000000000ahp", fmt.Sprintf("--from=%s", delegator_address), "-y", "--keyring-backend=file"})

	// Wait for a new block to ensure state updates
	time.Sleep(6 * time.Second)

	// Check delegation amount
	cmd := exec.Command("go", "run", path, "query", "staking", "delegation", delegator_address, validator_address)
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, "delegation should be queried correctly")
	re := regexp.MustCompile(`amount:\s*"(\d+)"\s*denom: ahp`)
	match := re.FindStringSubmatch(string(out))
	assert.Condition(t, func() bool { return len(match) > 1 }, "delegation amount should be in the output")
	assert.Condition(t, func() bool {
		return compareAmount(match[1], "1000000000000000000") > 0
	}, "delegation amount should be greater than initial deposit after delegation")
	delegationAmount := match[1]

	// Check balance
	cmd = exec.Command("go", "run", path, "query", "bank", "balances", delegator_address)
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err, "balance should be queried correctly")
	re = regexp.MustCompile(`amount:\s*"(\d+)"\s*denom: ahp`)
	match = re.FindStringSubmatch(string(out))
	assert.Condition(t, func() bool { return len(match) > 1 }, "balance should be in the output")
	balance := match[1]

	// Delegate more tokens to the validator
	testTx(t, []string{"tx", "staking", "delegate", validator_address, "500000000000000000000ahp", "--fees=1000000000000000000ahp", fmt.Sprintf("--from=%s", delegator_address), "-y", "--keyring-backend=file"})

	// Wait for a new block to ensure state updates
	time.Sleep(6 * time.Second)

	// Compare delegation amount and balance after delegating more tokens
	cmd = exec.Command("go", "run", path, "query", "bank", "balances", delegator_address)
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err, "balance should be queried correctly")
	re = regexp.MustCompile(`amount:\s*"(\d+)"\s*denom: ahp`)
	match = re.FindStringSubmatch(string(out))
	assert.Condition(t, func() bool { return len(match) > 1 }, "balance should be in the output")
	assert.Condition(t, func() bool {
		return compareAmount(match[1], balance) < 0
	}, "balance should be decreased after delegation")

	cmd = exec.Command("go", "run", path, "query", "staking", "delegation", delegator_address, validator_address)
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err, "delegation should be queried correctly")
	re = regexp.MustCompile(`amount:\s*"(\d+)"\s*denom: ahp`)
	match = re.FindStringSubmatch(string(out))
	assert.Condition(t, func() bool { return len(match) > 1 }, "delegation amount should be in the output")
	assert.Condition(t, func() bool {
		return compareAmount(match[1], delegationAmount) > 0
	}, "delegation amount should increase after delegating more tokens")
	delegationAmount = match[1]

	// Delegation reward should be accumulated
	cmd = exec.Command("go", "run", path, "query", "distribution", "rewards", delegator_address)
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err, "rewards should be queried correctly")
	assert.Contains(t, string(out), validator_address, "rewards from validator_address should be in the output")
	assert.Contains(t, string(out), "reward:", "rewards from delegating should be in the output")
}

func TestValidator(t *testing.T) {
	delegator_address := os.Getenv(key_delegator_address)

	// default value of moniker and min_self_delegation
	cmd := exec.Command("go", "run", path, "query", "staking", "validators")
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, "validator should be queried correctly")
	assert.Condition(t, func() bool { return strings.Contains(string(out), "moniker: hippo") }, "moniker should be hippo")
	assert.Condition(t, func() bool { return strings.Contains(string(out), `min_self_delegation: "1"`) }, "min_self_delegation should be 1")

	// edit validator
	testTx(t, []string{"tx", "staking", "edit-validator", "--new-moniker=newnewHippo", "--min-self-delegation=5", "--fees=1000000000000000000ahp", fmt.Sprintf("--from=%s", delegator_address), "-y", "--keyring-backend=file"})

	// Wait for a new block to ensure state updates
	time.Sleep(6 * time.Second)

	// changed value of moniker and min_self_delegation
	cmd = exec.Command("go", "run", path, "query", "staking", "validators")
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err, "validator should be queried correctly")
	assert.Condition(t, func() bool { return strings.Contains(string(out), "moniker: newnewHippo") }, "moniker should be changed")
	assert.Condition(t, func() bool { return strings.Contains(string(out), `min_self_delegation: "5"`) }, "min_self_delegation should be changed")

}

func TestProposal(t *testing.T) {
	delegator_address := os.Getenv(key_delegator_address)

	content := map[string]interface{}{
		"metadata":  "ipfs://CID",
		"deposit":   "1000000000000000000000ahp",
		"title":     "test proposal title",
		"summary":   "summary",
		"expedited": false,
	}
	jsonData, err := json.MarshalIndent(content, "", "  ")
	assert.NoError(t, err, "Failed to marshal JSON")

	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, "draft_proposal.json")

	err = os.WriteFile(filePath, jsonData, 0644)
	assert.NoError(t, err, "Failed to write JSON to file")

	testTx(t, []string{"tx", "gov", "submit-proposal", filePath, fmt.Sprintf("--from=%s", delegator_address), "--fees=1000000000000000000ahp", "-y", "--keyring-backend=file"})

	// Wait for a new block to ensure state updates
	time.Sleep(6 * time.Second)

	cmd := exec.Command("go", "run", path, "query", "gov", "proposals")
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Proposals should be queried correctly")
	assert.Condition(t, func() bool { return strings.Contains(string(out), "title: test proposal title") }, "proposal title should be in the output")

}

func TestCommission(t *testing.T) {
	delegator_address := os.Getenv(key_delegator_address)
	validator_address := os.Getenv(key_validator_address)

	cmd := exec.Command("go", "run", path, "query", "distribution", "commission", validator_address)
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, "validator commission should be queried correctly")

	re := regexp.MustCompile(`-\s*(\d+).*\d*ahp`)
	match := re.FindStringSubmatch(string(out))

	assert.Condition(t, func() bool { return len(match) > 1 }, "commission should be in the output")
	commission := match[1]

	testTx(t, []string{"tx", "distribution", "withdraw-rewards", "--commission", validator_address, fmt.Sprintf("--from=%s", delegator_address), "--fees=1000000000000000000ahp", "-y", "--keyring-backend=file"})

	time.Sleep(6 * time.Second)

	cmd = exec.Command("go", "run", path, "query", "distribution", "commission", validator_address)
	out, err = cmd.CombinedOutput()
	assert.NoError(t, err, "validator commission should be queried correctly")
	match = re.FindStringSubmatch(string(out))

	assert.Condition(t, func() bool {
		return compareAmount(match[1], commission) < 0
	}, "commimssion should be decreased after withdraw commission")
}
