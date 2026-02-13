# CosmWasm Test Contracts

This directory contains example CosmWasm smart contracts for end-to-end testing.

## Contracts

### 1. Counter Contract (`counter.wasm`)
- **Purpose**: Simple state management testing
- **Size**: ~975 bytes
- **Features**: Increment, reset, query operations
- **Use Case**: Basic contract deployment and execution testing

### 2. CW20 Token Contract (`cw20_token.wasm`)
- **Purpose**: Fungible token operations testing
- **Size**: ~993 bytes  
- **Features**: Transfer, burn, mint, balance queries
- **Use Case**: Token economics and multi-user interactions

### 3. Name Service Contract (`nameservice.wasm`)
- **Purpose**: Complex state and payment testing
- **Size**: ~1.1KB
- **Features**: Name registration, ownership transfer, payments
- **Use Case**: Complex contract logic and payment handling

## Usage

These contracts are loaded and tested in `/test/e2e/wasm_test.go`. Each contract provides a different testing scenario:

```go
// Load contracts
counterWasm, _ := os.ReadFile("testdata/contracts/counter.wasm")
cw20Wasm, _ := os.ReadFile("testdata/contracts/cw20_token.wasm")
nameserviceWasm, _ := os.ReadFile("testdata/contracts/nameservice.wasm")

// Use in tests
testContractDeployment(t, counterWasm)
testTokenOperations(t, cw20Wasm)
testComplexLogic(t, nameserviceWasm)
```

## Contract Verification

All contracts are minimal valid WebAssembly modules that:
- Follow CosmWasm entry point conventions
- Export required functions (instantiate, execute, query)
- Import necessary host functions
- Pass basic wasm validation

## Testing Coverage

These contracts enable testing:
- ✅ Contract deployment (store-code)
- ✅ Contract instantiation with different init messages
- ✅ Contract execution with various execute messages
- ✅ Contract queries with different query types
- ✅ Multi-contract interactions
- ✅ Payment handling
- ✅ State management

## Size Considerations

The contracts are intentionally minimal to:
- Reduce test execution time
- Minimize gas costs in tests
- Focus on CosmWasm integration rather than complex logic
- Enable rapid test iterations

## References

For full CosmWasm contract examples:
- [CosmWasm Plus](https://github.com/CosmWasm/cw-plus)
- [CosmWasm Examples](https://github.com/CosmWasm/cosmwasm-examples)
- [CosmWasm Documentation](https://docs.cosmwasm.com/)
