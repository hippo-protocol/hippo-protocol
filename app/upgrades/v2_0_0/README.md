# CosmWasm v2.0.0 Upgrade Tests

This directory contains comprehensive tests for the v2.0.0 upgrade that implements CosmWasm module support in the Hippo Protocol.

## Test Files

### upgrade_test.go
Unit tests for the upgrade configuration and basic functionality:
- `TestUpgradeNameConstant` - Verifies upgrade name is correct
- `TestUpgradeStoreConfiguration` - Validates store upgrade configuration
- `TestUpgradeConfiguration` - Checks upgrade object structure
- `TestUpgradeHandlerCreation` - Tests handler creation
- `TestWasmModuleConfiguration` - Verifies wasm module is properly configured
- `TestCosmWasmDefaultParamsStructure` - Validates CosmWasm parameter structure
- `TestUpgradePlan` - Tests upgrade plan creation
- `TestCosmWasmAccessTypes` - Verifies access type constants
- `TestStoreUpgradesNotEmpty` - Ensures store upgrades are configured
- `TestUpgradeHandlerSignature` - Validates handler signature

### integration_test.go
Integration tests for the upgrade handler with realistic setup:
- `TestUpgradeHandlerIntegration` - Full integration test of upgrade handler execution
  - Sets up complete keeper infrastructure
  - Executes upgrade handler
  - Verifies CosmWasm parameters are set correctly
  - Confirms WasmKeeper is functional after upgrade
- `TestWasmKeeperAccessibility` - Verifies WasmKeeper initialization and accessibility
- `TestWasmModuleInStoreUpgrades` - Validates store upgrade configuration details
- `TestUpgradeReadiness` - Production readiness checklist

## E2E Tests

The repository also includes end-to-end tests in `/test/e2e/wasm_test.go`:

### CosmWasm Functionality Tests
- `TestWasmQuery` - Basic wasm query commands
- `TestWasmParams` - Verifies parameters after v2.0.0 upgrade
- `TestWasmStoreCode` - Tests contract deployment (store-code)
- `TestWasmInstantiateContract` - Tests contract instantiation
- `TestWasmExecuteContract` - Tests contract execution (sendTx)
- `TestWasmSendFunds` - Tests sending funds to contracts
- `TestWasmAPI` - Tests wasm API endpoints

## Running Tests

### Unit and Integration Tests
```bash
# Run all upgrade tests
go test ./app/upgrades/v2_0_0/... -v

# Run specific test
go test ./app/upgrades/v2_0_0 -run TestUpgradeHandlerIntegration -v

# Run with coverage
go test ./app/upgrades/v2_0_0/... -cover
```

### E2E Tests
E2E tests require a running node. They are executed as part of the CI pipeline:

```bash
# Start a test node (see .github/workflows/go.yml for setup)
go run hippod/main.go init hippo --chain-id hippo-protocol-testnet-1
# ... additional setup ...
go run hippod/main.go start &

# Run e2e tests
go test ./test/e2e -v
```

## Test Coverage

The test suite covers:

✅ **Upgrade Configuration**
- Upgrade name and version
- Store upgrade configuration (wasm module addition)
- Handler creation and signature

✅ **CosmWasm Parameters**
- Code upload access set to `AllowEverybody`
- Instantiate default permission set to `AccessTypeEverybody`
- Parameter structure validation

✅ **WasmKeeper Functionality**
- Keeper initialization
- Keeper accessibility
- Parameter getting/setting

✅ **End-to-End Workflows**
- Contract deployment (store-code)
- Contract instantiation
- Contract execution (call/sendTx)
- Contract queries
- Fund transfers to contracts
- API endpoint accessibility

## Production Deployment

Before deploying the v2.0.0 upgrade to production:

1. ✅ All unit tests pass
2. ✅ All integration tests pass
3. ✅ All e2e tests pass in CI
4. ✅ CosmWasm parameters verified
5. ✅ Store upgrades configured correctly
6. ✅ Upgrade handler tested with realistic setup

## CI/CD Integration

Tests are automatically run in the GitHub Actions workflow (`.github/workflows/go.yml`):

- **Unit Tests**: Run on every PR
- **Integration Tests**: Run on every PR
- **E2E Tests**: Run in genesis test job with full node setup

## Upgrade Process

When the network upgrades to v2.0.0:

1. The upgrade handler will be triggered at the specified block height
2. The wasm module will be added to the store
3. CosmWasm parameters will be initialized:
   - Code upload: Open to everyone
   - Contract instantiation: Open to everyone by default
4. The network will continue operating with CosmWasm support enabled

## Security Considerations

- CosmWasm parameters are set to allow everyone to upload and instantiate contracts
- This is intentional for this upgrade to enable permissionless smart contract deployment
- Individual contracts can still implement their own access control
- Gas limits protect against resource exhaustion

## Support

For issues or questions about the upgrade:
- Check test logs for specific failure details
- Review CosmWasm documentation: https://docs.cosmwasm.com/
- Consult the upgrade handler code in `upgrade.go`
