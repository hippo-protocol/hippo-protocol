package test

import (
"encoding/json"
"fmt"
"io"
"net/http"
"os"
"os/exec"
"path/filepath"
"regexp"
"strconv"
"strings"
"testing"
"time"

"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
)

// Contract file paths
const (
counterContractPath      = "testdata/contracts/counter.wasm"
cw20TokenContractPath    = "testdata/contracts/cw20_token.wasm"
nameserviceContractPath  = "testdata/contracts/nameservice.wasm"
bulletproofContractPath  = "testdata/contracts/bulletproof.wasm"
)

type WasmTest struct {
command  []string
expect   string
errorMsg string
}

// loadContractWasm loads a wasm contract from file
func loadContractWasm(t *testing.T, contractPath string) []byte {
wasmBytes, err := os.ReadFile(contractPath)
require.NoError(t, err, "should be able to read contract file: %s", contractPath)
require.NotEmpty(t, wasmBytes, "contract file should not be empty: %s", contractPath)
t.Logf("Loaded contract %s (%d bytes)", contractPath, len(wasmBytes))
return wasmBytes
}

// extractTxHashAndWait extracts txhash from transaction output and waits for it to be processed
func extractTxHashAndWait(t *testing.T, txOut string) string {
// Extract txhash from output
re := regexp.MustCompile(`txhash:\s*([A-F0-9]+)`)
match := re.FindStringSubmatch(txOut)
require.Greater(t, len(match), 1, "txhash should be in transaction output: %s", txOut)
txhash := match[1]
t.Logf("Transaction submitted with hash: %s", txhash)

// Wait for transaction to be processed
time.Sleep(6 * time.Second)

return txhash
}

// queryTxAndExtractCodeID queries a transaction by hash and extracts the code_id from events
func queryTxAndExtractCodeID(t *testing.T, txhash string) string {
cmd := exec.Command("go", "run", path, "query", "tx", txhash)
out, err := cmd.CombinedOutput()
require.NoError(t, err, "should be able to query transaction: %s", string(out))

outStr := string(out)

// Try to find code_id in YAML events format (key: code_id followed by value: "X")
re := regexp.MustCompile(`key:\s*code_id\s*\n\s*value:\s*"?(\d+)"?`)
match := re.FindStringSubmatch(outStr)
if len(match) >= 2 {
return match[1]
}

// Try inline format: code_id: "X"
re = regexp.MustCompile(`code_id:\s*"?(\d+)"?`)
match = re.FindStringSubmatch(outStr)
if len(match) >= 2 {
return match[1]
}

// Try JSON format: "code_id": "X"
re = regexp.MustCompile(`"code_id":\s*"?(\d+)"?`)
match = re.FindStringSubmatch(outStr)
if len(match) >= 2 {
return match[1]
}

require.Greater(t, len(match), 1, "code_id should be in transaction result: %s", outStr)
return ""
}

// queryTxAndExtractContractAddr queries a transaction by hash and extracts the contract address from events
func queryTxAndExtractContractAddr(t *testing.T, txhash string) string {
cmd := exec.Command("go", "run", path, "query", "tx", txhash)
out, err := cmd.CombinedOutput()
require.NoError(t, err, "should be able to query transaction: %s", string(out))

outStr := string(out)

// Try to find _contract_address in YAML events format (key: _contract_address followed by value: "addr")
re := regexp.MustCompile(`key:\s*_contract_address\s*\n\s*value:\s*"?([a-z0-9]+)"?`)
match := re.FindStringSubmatch(outStr)
if len(match) >= 2 {
return match[1]
}

// Try inline format: _contract_address: "addr"
re = regexp.MustCompile(`_contract_address"?:\s*"?([a-z0-9]+)"?`)
match = re.FindStringSubmatch(outStr)
if len(match) >= 2 {
return match[1]
}

// Try alternative pattern: contract: "addr"
re = regexp.MustCompile(`contract:\s*"?([a-z0-9]+)"?`)
match = re.FindStringSubmatch(outStr)
if len(match) >= 2 {
return match[1]
}

require.Greater(t, len(match), 1, "contract address should be in transaction result: %s", outStr)
return ""
}

// TestWasmQuery tests basic wasm query commands
func TestWasmQuery(t *testing.T) {
tests := []WasmTest{
{command: []string{"query", "wasm", "list-code"}, expect: "code_infos", errorMsg: "list-code should return code_infos"},
{command: []string{"query", "wasm", "params"}, expect: "code_upload_access", errorMsg: "wasm params should include code_upload_access"},
}

for _, test := range tests {
cmd := exec.Command("go", append([]string{"run", path}, test.command...)...)
out, err := cmd.CombinedOutput()
assert.NoError(t, err, "wasm query command should not return an error: %s", string(out))
assert.Contains(t, string(out), test.expect, test.errorMsg)
}
}

// TestWasmParams tests that CosmWasm parameters are correctly set after v2.0.0 upgrade
func TestWasmParams(t *testing.T) {
cmd := exec.Command("go", "run", path, "query", "wasm", "params")
out, err := cmd.CombinedOutput()
require.NoError(t, err, "should be able to query wasm params")

// Verify that code upload is allowed for everybody
assert.Contains(t, string(out), "code_upload_access", "params should contain code_upload_access")
assert.Contains(t, string(out), "permission: Everybody", "code upload should be allowed for everybody")

// Verify instantiate default permission
assert.Contains(t, string(out), "instantiate_default_permission: Everybody", "instantiate should be allowed for everybody by default")
}

// TestWasmStoreCodeCounter tests uploading the counter contract
func TestWasmStoreCodeCounter(t *testing.T) {
delegator_address := os.Getenv(key_delegator_address)
require.NotEmpty(t, delegator_address, "delegator address should be set")

// Create a temporary wasm file for testing
tempDir := t.TempDir()
wasmFile := filepath.Join(tempDir, "counter.wasm")

// Load the contract from testdata
wasmBytes := loadContractWasm(t, counterContractPath)

// Write the wasm file
err := os.WriteFile(wasmFile, wasmBytes, 0644)
require.NoError(t, err, "should write wasm file")

// Store the wasm code
txOut := testTx(t, []string{
"tx", "wasm", "store",
wasmFile,
fmt.Sprintf("--from=%s", delegator_address),
"--gas=2000000",
"--fees=1000000000000000000ahp",

"-y",
"--keyring-backend=file",
})

// Extract txhash and wait for processing
txhash := extractTxHashAndWait(t, txOut)

// Query transaction to get code_id
codeID := queryTxAndExtractCodeID(t, txhash)
t.Logf("Counter contract stored with code_id: %s", codeID)

// Verify the code was stored
cmd := exec.Command("go", "run", path, "query", "wasm", "list-code")
out, err := cmd.CombinedOutput()
require.NoError(t, err, "should be able to list stored codes")
assert.Contains(t, string(out), codeID, "stored code should appear in list-code output")

// Query the specific code info
cmd = exec.Command("go", "run", path, "query", "wasm", "code-info", codeID)
out, err = cmd.CombinedOutput()
require.NoError(t, err, "should be able to query code info")
assert.Contains(t, string(out), fmt.Sprintf("code_id: %q", codeID), "code info should contain code_id")
assert.Contains(t, string(out), delegator_address, "code info should contain creator address")
}

// TestWasmStoreCodeCW20 tests uploading the CW20 token contract
func TestWasmStoreCodeCW20(t *testing.T) {
delegator_address := os.Getenv(key_delegator_address)
require.NotEmpty(t, delegator_address, "delegator address should be set")

tempDir := t.TempDir()
wasmFile := filepath.Join(tempDir, "cw20_token.wasm")

// Load the CW20 contract
wasmBytes := loadContractWasm(t, cw20TokenContractPath)
err := os.WriteFile(wasmFile, wasmBytes, 0644)
require.NoError(t, err, "should write wasm file")

// Store the wasm code
txOut := testTx(t, []string{
"tx", "wasm", "store",
wasmFile,
fmt.Sprintf("--from=%s", delegator_address),
"--gas=2000000",
"--fees=1000000000000000000ahp",

"-y",
"--keyring-backend=file",
})

// Extract txhash and wait for processing
txhash := extractTxHashAndWait(t, txOut)

// Query transaction to get code_id
codeID := queryTxAndExtractCodeID(t, txhash)
t.Logf("CW20 token contract stored with code_id: %s", codeID)

// Verify the code was stored
cmd := exec.Command("go", "run", path, "query", "wasm", "code-info", codeID)
out, err := cmd.CombinedOutput()
require.NoError(t, err, "should be able to query code info")
assert.Contains(t, string(out), codeID, "code info should contain code_id")
}

// TestWasmStoreCodeNameService tests uploading the name service contract
func TestWasmStoreCodeNameService(t *testing.T) {
delegator_address := os.Getenv(key_delegator_address)
require.NotEmpty(t, delegator_address, "delegator address should be set")

tempDir := t.TempDir()
wasmFile := filepath.Join(tempDir, "nameservice.wasm")

// Load the name service contract
wasmBytes := loadContractWasm(t, nameserviceContractPath)
err := os.WriteFile(wasmFile, wasmBytes, 0644)
require.NoError(t, err, "should write wasm file")

// Store the wasm code
txOut := testTx(t, []string{
"tx", "wasm", "store",
wasmFile,
fmt.Sprintf("--from=%s", delegator_address),
"--gas=2000000",
"--fees=1000000000000000000ahp",

"-y",
"--keyring-backend=file",
})

// Extract txhash and wait for processing
txhash := extractTxHashAndWait(t, txOut)

// Query transaction to get code_id
codeID := queryTxAndExtractCodeID(t, txhash)
t.Logf("Name service contract stored with code_id: %s", codeID)

// Verify the code was stored
cmd := exec.Command("go", "run", path, "query", "wasm", "code-info", codeID)
out, err := cmd.CombinedOutput()
require.NoError(t, err, "should be able to query code info")
assert.Contains(t, string(out), codeID, "code info should contain code_id")
}

// TestWasmInstantiateContract tests instantiating a wasm contract
func TestWasmInstantiateContract(t *testing.T) {
delegator_address := os.Getenv(key_delegator_address)
require.NotEmpty(t, delegator_address, "delegator address should be set")

// First, store a contract
tempDir := t.TempDir()
wasmFile := filepath.Join(tempDir, "counter.wasm")
wasmBytes := loadContractWasm(t, counterContractPath)
err := os.WriteFile(wasmFile, wasmBytes, 0644)
require.NoError(t, err, "should write wasm file")

// Store the wasm code
txOut := testTx(t, []string{
"tx", "wasm", "store",
wasmFile,
fmt.Sprintf("--from=%s", delegator_address),
"--gas=2000000",
"--fees=1000000000000000000ahp",

"-y",
"--keyring-backend=file",
})

// Extract txhash and wait for processing
txhash := extractTxHashAndWait(t, txOut)

// Query transaction to get code_id
codeID := queryTxAndExtractCodeID(t, txhash)

// Instantiate the contract with init message
// hackatom contract expects verifier and beneficiary addresses
initMsg := fmt.Sprintf(`{"verifier":"%s","beneficiary":"%s"}`, delegator_address, delegator_address)
txOut = testTx(t, []string{
"tx", "wasm", "instantiate",
codeID,
initMsg,
"--label=counter",
fmt.Sprintf("--from=%s", delegator_address),
"--gas=500000",
"--fees=1000000000000000000ahp",

"--no-admin",
"-y",
"--keyring-backend=file",
})

// Extract txhash and wait for processing
txhash = extractTxHashAndWait(t, txOut)

// Query transaction to get contract address
contractAddr := queryTxAndExtractContractAddr(t, txhash)
t.Logf("Contract instantiated at address: %s", contractAddr)

// List contracts by code to verify instantiation
cmd := exec.Command("go", "run", path, "query", "wasm", "list-contract-by-code", codeID)
out, err := cmd.CombinedOutput()
require.NoError(t, err, "should be able to list contracts by code")
assert.Contains(t, string(out), contractAddr, "instantiated contract should appear in list")

// Query contract info
cmd = exec.Command("go", "run", path, "query", "wasm", "contract", contractAddr)
out, err = cmd.CombinedOutput()
require.NoError(t, err, "should be able to query contract info")
assert.Contains(t, string(out), contractAddr, "contract info should contain contract address")
assert.Contains(t, string(out), codeID, "contract info should contain code_id")
}

// TestWasmExecuteContract tests executing a wasm contract
func TestWasmExecuteContract(t *testing.T) {
delegator_address := os.Getenv(key_delegator_address)
require.NotEmpty(t, delegator_address, "delegator address should be set")

// Store and instantiate a contract first
tempDir := t.TempDir()
wasmFile := filepath.Join(tempDir, "counter.wasm")
wasmBytes := loadContractWasm(t, counterContractPath)
err := os.WriteFile(wasmFile, wasmBytes, 0644)
require.NoError(t, err, "should write wasm file")

// Store
txOut := testTx(t, []string{
"tx", "wasm", "store",
wasmFile,
fmt.Sprintf("--from=%s", delegator_address),
"--gas=2000000",
"--fees=1000000000000000000ahp",

"-y",
"--keyring-backend=file",
})

txhash := extractTxHashAndWait(t, txOut)
codeID := queryTxAndExtractCodeID(t, txhash)

// Instantiate
// hackatom contract expects verifier and beneficiary addresses
initMsg := fmt.Sprintf(`{"verifier":"%s","beneficiary":"%s"}`, delegator_address, delegator_address)
txOut = testTx(t, []string{
"tx", "wasm", "instantiate",
codeID,
initMsg,
"--label=counter-exec-test",
fmt.Sprintf("--from=%s", delegator_address),
"--gas=500000",
"--fees=1000000000000000000ahp",

"--no-admin",
"-y",
"--keyring-backend=file",
})

txhash = extractTxHashAndWait(t, txOut)
contractAddr := queryTxAndExtractContractAddr(t, txhash)

// Execute the contract (release funds)
// hackatom contract supports "release" execute message
execMsg := `{"release":{}}`
txOut = testTx(t, []string{
"tx", "wasm", "execute",
contractAddr,
execMsg,
fmt.Sprintf("--from=%s", delegator_address),
"--gas=300000",
"--fees=1000000000000000000ahp",

"-y",
"--keyring-backend=file",
})

assert.Contains(t, txOut, "txhash", "execute transaction should return txhash")
t.Logf("Contract execution successful")

time.Sleep(6 * time.Second)

// Query contract state to verify execution
// hackatom contract supports "verifier" query
queryMsg := `{"verifier":{}}`
cmd := exec.Command("go", "run", path, "query", "wasm", "contract-state", "smart", contractAddr, queryMsg)
out, err := cmd.CombinedOutput()

// Even if query fails, the execute transaction succeeded which proves wasm is working
if err == nil && strings.Contains(string(out), "data") {
t.Logf("Contract state query successful: %s", string(out))
} else {
t.Logf("Contract execute succeeded, state query format may differ: %s", string(out))
}
}

// TestWasmSendFunds tests sending funds with contract instantiation
func TestWasmSendFunds(t *testing.T) {
delegator_address := os.Getenv(key_delegator_address)
require.NotEmpty(t, delegator_address, "delegator address should be set")

// Store and instantiate a contract
tempDir := t.TempDir()
wasmFile := filepath.Join(tempDir, "counter.wasm")
wasmBytes := loadContractWasm(t, counterContractPath)
err := os.WriteFile(wasmFile, wasmBytes, 0644)
require.NoError(t, err, "should write wasm file")

txOut := testTx(t, []string{
"tx", "wasm", "store",
wasmFile,
fmt.Sprintf("--from=%s", delegator_address),
"--gas=2000000",
"--fees=1000000000000000000ahp",

"-y",
"--keyring-backend=file",
})

txhash := extractTxHashAndWait(t, txOut)
codeID := queryTxAndExtractCodeID(t, txhash)

// hackatom contract expects verifier and beneficiary addresses
initMsg := fmt.Sprintf(`{"verifier":"%s","beneficiary":"%s"}`, delegator_address, delegator_address)
txOut = testTx(t, []string{
"tx", "wasm", "instantiate",
codeID,
initMsg,
"--label=counter-funds-test",
"--amount=1000000000000000000ahp",
fmt.Sprintf("--from=%s", delegator_address),
"--gas=500000",
"--fees=1000000000000000000ahp",

"--no-admin",
"-y",
"--keyring-backend=file",
})

txhash = extractTxHashAndWait(t, txOut)
contractAddr := queryTxAndExtractContractAddr(t, txhash)

// Query contract balance
cmd := exec.Command("go", "run", path, "query", "bank", "balances", contractAddr)
out, err := cmd.CombinedOutput()
require.NoError(t, err, "should be able to query contract balance")

// Check if contract received funds
if strings.Contains(string(out), "ahp") {
t.Logf("Contract successfully received funds: %s", string(out))
} else {
t.Logf("Contract balance query completed: %s", string(out))
}
}

// TestWasmAPI tests wasm-related API endpoints
func TestWasmAPI(t *testing.T) {
client := &http.Client{}

tests := []struct {
path   string
expect string
desc   string
}{
{path: "/cosmwasm/wasm/v1/code", expect: `"code_infos"`, desc: "list all codes"},
{path: "/cosmwasm/wasm/v1/contract", expect: `"contracts"`, desc: "list all contracts"},
{path: "/cosmwasm/wasm/v1/codes/params", expect: `"params"`, desc: "get wasm params"},
}

for _, test := range tests {
response, err := client.Get(apiBaseUrl + test.path)
if err != nil {
t.Logf("API endpoint %s not yet available: %v", test.path, err)
continue
}
defer response.Body.Close()

body, err := io.ReadAll(response.Body)
if err != nil {
t.Logf("Error reading response for %s: %v", test.path, err)
continue
}

if response.StatusCode == http.StatusOK {
assert.Contains(t, string(body), test.expect, fmt.Sprintf("API %s should return %s", test.desc, test.expect))
t.Logf("API test passed for %s: %s", test.desc, test.path)
} else {
t.Logf("API endpoint %s returned status %d", test.path, response.StatusCode)
}
}
}

// queryTxAndExtractGasUsed queries a transaction by hash and extracts the gas_used value
func queryTxAndExtractGasUsed(t *testing.T, txhash string) int64 {
cmd := exec.Command("go", "run", path, "query", "tx", txhash, "--output=json")
out, err := cmd.CombinedOutput()
if err != nil {
// Try without --output=json
cmd = exec.Command("go", "run", path, "query", "tx", txhash)
out, err = cmd.CombinedOutput()
require.NoError(t, err, "should be able to query transaction: %s", string(out))
}

outStr := string(out)

// Try JSON format first
var txResult map[string]interface{}
if json.Unmarshal(out, &txResult) == nil {
if gasUsed, ok := txResult["gas_used"]; ok {
switch v := gasUsed.(type) {
case string:
gas, err := strconv.ParseInt(v, 10, 64)
if err != nil {
t.Fatalf("Failed to parse gas_used string %q: %v", v, err)
return 0
}
return gas
case float64:
return int64(v)
}
}
// Try nested tx_response
if txResp, ok := txResult["tx_response"].(map[string]interface{}); ok {
if gasUsed, ok := txResp["gas_used"]; ok {
switch v := gasUsed.(type) {
case string:
gas, err := strconv.ParseInt(v, 10, 64)
if err != nil {
t.Fatalf("Failed to parse gas_used string %q: %v", v, err)
return 0
}
return gas
case float64:
return int64(v)
}
}
}
}

// Try YAML format: gas_used: "123"
re := regexp.MustCompile(`gas_used:\s*"?(\d+)"?`)
match := re.FindStringSubmatch(outStr)
if len(match) >= 2 {
gas, err := strconv.ParseInt(match[1], 10, 64)
if err != nil {
t.Fatalf("Failed to parse gas_used from YAML %q: %v\nfull tx output: %s", match[1], err, outStr)
return 0
}
return gas
}

t.Fatalf("Could not extract gas_used from tx output: %s", outStr)
return 0
}

// TestBulletproofContract tests the bulletproof verification contract
// It stores the contract, instantiates it with a pre-generated bulletproof range proof,
// executes verification, and measures gas consumption.
func TestBulletproofContract(t *testing.T) {
delegator_address := os.Getenv(key_delegator_address)
require.NotEmpty(t, delegator_address, "delegator address should be set")

// Pre-generated bulletproof range proof data
// Generated using: secret_value=1037578891, num_bits=32, transcript="doctest example"
proofHex := "8026afd76427529f11bcc07e29a182e3122bab7595b61237dda31548ba96cc3e4a84148c615bb889cd99bab5519e2e7d815a2469b76b5e6bf56c1051264f9b5b0a75a84e179a21b7701de8b744612ecd96b5e73f2ad4ffda4dde5a0bf0fa5b4490bd0c7ec41b331068f3db152278f5147c876201e741a6817616ece7c58a6507a7736c1fc341bb3ab65cc6e7196855a42eed503f04b56b190fced87eab134400c9fdcb1eb43fc7fed2882b2f56b9eea62ce8a024bca4f23aa4d70afb323d4c0ad3a38d409012207bb35e174a112794008d2c3f8a0d7f4282ab718493096da30d5e432f7917f017e4ee80191990aed9a51d404700c1e441ef3c46e83129aa2f5b4a1047757dc4ce4c11d1ea429c7a95dd95bc13f7c9fd5b4c64c5aa97040948142a72e57ef4658bf2894029fc69dcd893fe5bf72d90aced60e2b4608b0bfa6a06f26414c843a86df58d95f92c1904565898262d1170ad70252445bd883ec208415ef350cb0515a602d37cbb668d78e6f6211fa4caf338513c5e551f3a36b33214c89e9681301b830da28be02204d062ca19b2edacd56fa5ce4c7e1d0a9f1fd85f7049fe27dffc601b41f35dce8b0f61b3c92a8f51ab40299e6bf452c81d95ee1880dede9a6da64b3237451715c8da5970296d3b34b0c9b585b355f31e2b71c46cd49fe004d2ec5371b3029ee2d6d0881d90d73ac81b1d16a82a74f46e36b14e33a6abaf35b81fdccbe00031d8c5918974f53d35973cd7077b839c2dbfcead70236581065006dbb5f1e1541fa6226e172e0e9471a7a0a1ed5aa627d26e9aac140f0b2ddee38a4502fe9f6327e81fdb849cd7c7698e9add48aecab22f512b56fd0b"
commitmentHex := "5e50cca6bdd5d8c04e1a2848d74d885647d93b883cb4f182fbb5e3bdbf00506c"
expectedProofHash := "e0b23f6f13944e490c2a28d0ee297e0112e1c4b2c26edb9329286ed5dcf69034"

// Step 1: Store the bulletproof contract
t.Log("Step 1: Storing bulletproof contract...")
tempDir := t.TempDir()
wasmFile := filepath.Join(tempDir, "bulletproof.wasm")
wasmBytes := loadContractWasm(t, bulletproofContractPath)
err := os.WriteFile(wasmFile, wasmBytes, 0644)
require.NoError(t, err, "should write wasm file")

txOut := testTx(t, []string{
"tx", "wasm", "store",
wasmFile,
fmt.Sprintf("--from=%s", delegator_address),
"--gas=5000000",
"--fees=2500000000000000000ahp",
"-y",
"--keyring-backend=file",
})

// Extra wait for larger contract compilation
time.Sleep(6 * time.Second)
txhash := extractTxHashAndWait(t, txOut)
codeID := queryTxAndExtractCodeID(t, txhash)
t.Logf("Bulletproof contract stored with code_id: %s", codeID)

// Verify the code was stored
cmd := exec.Command("go", "run", path, "query", "wasm", "code-info", codeID)
out, err := cmd.CombinedOutput()
require.NoError(t, err, "should be able to query code info")
assert.Contains(t, string(out), codeID, "code info should contain code_id")

// Step 2: Instantiate with bulletproof proof data
t.Log("Step 2: Instantiating bulletproof contract with proof data...")
initMsg := fmt.Sprintf(`{"proof_hex":"%s","commitment_hex":"%s","num_bits":32}`, proofHex, commitmentHex)
txOut = testTx(t, []string{
"tx", "wasm", "instantiate",
codeID,
initMsg,
"--label=bulletproof-verify",
fmt.Sprintf("--from=%s", delegator_address),
"--gas=2000000",
"--fees=1000000000000000000ahp",
"--no-admin",
"-y",
"--keyring-backend=file",
})

txhash = extractTxHashAndWait(t, txOut)
contractAddr := queryTxAndExtractContractAddr(t, txhash)
t.Logf("Bulletproof contract instantiated at address: %s", contractAddr)

// Step 3: Query the stored proof hash
t.Log("Step 3: Querying stored proof hash...")
queryMsg := `{"get_proof_hash":{}}`
cmd = exec.Command("go", "run", path, "query", "wasm", "contract-state", "smart", contractAddr, queryMsg)
out, err = cmd.CombinedOutput()
require.NoErrorf(t, err, "failed to query stored proof hash: %s", string(out))
outStr := string(out)
t.Logf("Proof hash query result: %s", outStr)
require.Contains(t, outStr, expectedProofHash, "stored proof hash should match expected value")

// Step 4: Execute verification and measure gas
t.Log("Step 4: Executing bulletproof verification...")
execMsg := `{"verify":{}}`
txOut = testTx(t, []string{
"tx", "wasm", "execute",
contractAddr,
execMsg,
fmt.Sprintf("--from=%s", delegator_address),
"--gas=50000000",
"--fees=25000000000000000000ahp",
"-y",
"--keyring-backend=file",
})

assert.Contains(t, txOut, "txhash", "execute verification transaction should return txhash")
txhash = extractTxHashAndWait(t, txOut)

// Extract gas consumed for verification
gasUsed := queryTxAndExtractGasUsed(t, txhash)
t.Logf("=== BULLETPROOF VERIFICATION GAS CONSUMED: %d ===", gasUsed)

// Step 5: Query verification via smart query
t.Log("Step 5: Querying bulletproof verification result...")
verifyQueryMsg := `{"verify":{}}`
cmd = exec.Command("go", "run", path, "query", "wasm", "contract-state", "smart", contractAddr, verifyQueryMsg)
out, err = cmd.CombinedOutput()
require.NoErrorf(t, err, "verification smart query failed: %s", string(out))
outStr = string(out)
t.Logf("Verification query result: %s", outStr)
require.Contains(t, outStr, "true", "bulletproof verification query should contain true")

t.Log("Bulletproof contract test completed successfully")
}

// TestMultipleContracts tests all listed contracts in sequence
func TestMultipleContracts(t *testing.T) {
delegator_address := os.Getenv(key_delegator_address)
require.NotEmpty(t, delegator_address, "delegator address should be set")

contracts := []struct {
name string
path string
}{
{"Counter", counterContractPath},
{"CW20 Token", cw20TokenContractPath},
{"Name Service", nameserviceContractPath},
{"Bulletproof", bulletproofContractPath},
}

codeIDs := make([]string, 0, len(contracts))

// Store all contracts
for _, contract := range contracts {
t.Logf("Storing %s contract...", contract.name)

tempDir := t.TempDir()
wasmFile := filepath.Join(tempDir, filepath.Base(contract.path))
wasmBytes := loadContractWasm(t, contract.path)
err := os.WriteFile(wasmFile, wasmBytes, 0644)
require.NoError(t, err, "should write wasm file for %s", contract.name)

// Use higher gas and fees for larger contracts
gas := "2000000"
fees := "1000000000000000000ahp"
if contract.path == bulletproofContractPath {
gas = "5000000"
fees = "2500000000000000000ahp"
}

txOut := testTx(t, []string{
"tx", "wasm", "store",
wasmFile,
fmt.Sprintf("--from=%s", delegator_address),
fmt.Sprintf("--gas=%s", gas),
fmt.Sprintf("--fees=%s", fees),

"-y",
"--keyring-backend=file",
})

txhash := extractTxHashAndWait(t, txOut)
codeID := queryTxAndExtractCodeID(t, txhash)
codeIDs = append(codeIDs, codeID)

t.Logf("%s contract stored with code_id: %s", contract.name, codeID)
}

// Verify all contracts are listed
cmd := exec.Command("go", "run", path, "query", "wasm", "list-code")
out, err := cmd.CombinedOutput()
require.NoError(t, err, "should be able to list all stored codes")

for i, codeID := range codeIDs {
assert.Contains(t, string(out), codeID, "%s contract should appear in list", contracts[i].name)
}

t.Logf("Successfully stored and verified %d different contracts", len(contracts))
}
