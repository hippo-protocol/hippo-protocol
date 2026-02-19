# CosmWasm Test Contracts

This directory contains example CosmWasm smart contracts for end-to-end testing, with **full Rust source code** included.

## 📦 Contracts Overview

### 1. Counter Contract (`counter.wasm`)
- **Purpose**: Simple state management testing
- **Size**: ~975 bytes (compiled)
- **Features**: Increment, reset, query operations
- **Use Case**: Basic contract deployment and execution testing
- **Source**: [`counter-contract/`](./counter-contract/)

### 2. CW20 Token Contract (`cw20_token.wasm`)
- **Purpose**: Fungible token operations testing
- **Size**: ~993 bytes (compiled)
- **Features**: Transfer, burn, mint, balance queries
- **Use Case**: Token economics and multi-user interactions
- **Source**: [`cw20-token-contract/`](./cw20-token-contract/)

### 3. Name Service Contract (`nameservice.wasm`)
- **Purpose**: Complex state and payment testing
- **Size**: ~1.1KB (compiled)
- **Features**: Name registration, ownership transfer, payments
- **Use Case**: Complex contract logic and payment handling
- **Source**: [`nameservice-contract/`](./nameservice-contract/)

## 🦀 Rust Source Code

Each contract includes **production-ready Rust source code**:

```
contracts/
├── counter-contract/
│   ├── Cargo.toml          # Dependencies and build config
│   ├── README.md           # Contract-specific documentation
│   └── src/
│       ├── lib.rs          # Library entry point
│       ├── contract.rs     # Contract logic
│       ├── msg.rs          # Message definitions
│       ├── state.rs        # State management
│       └── error.rs        # Error types
├── cw20-token-contract/
│   ├── Cargo.toml
│   ├── README.md
│   └── src/...
├── nameservice-contract/
│   ├── Cargo.toml
│   ├── README.md
│   └── src/...
├── counter.wasm            # Compiled binary
├── cw20_token.wasm         # Compiled binary
└── nameservice.wasm        # Compiled binary
```

## 🚀 Building from Source

### Prerequisites

1. **Install Rust**:
```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

2. **Add WebAssembly target**:
```bash
rustup target add wasm32-unknown-unknown
```

3. **Install optimizer (for production)**:
```bash
docker pull cosmwasm/optimizer:0.15.0
```

### Build Any Contract

```bash
# Navigate to contract directory
cd counter-contract

# Run tests
cargo test

# Development build
cargo build

# Production build
cargo build --release --target wasm32-unknown-unknown

# Optimized build (recommended)
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/optimizer:0.15.0
```

### Build All Contracts

```bash
# Build script for all contracts
for contract in counter-contract cw20-token-contract nameservice-contract; do
  echo "Building $contract..."
  cd $contract
  cargo build --release --target wasm32-unknown-unknown
  cd ..
done
```

## 🧪 Testing in Go Tests

These contracts are loaded and tested in `/test/e2e/wasm_test.go`:

```go
// Load contracts from compiled binaries
counterWasm := loadContractWasm(t, "testdata/contracts/counter.wasm")
cw20Wasm := loadContractWasm(t, "testdata/contracts/cw20_token.wasm")
nameserviceWasm := loadContractWasm(t, "testdata/contracts/nameservice.wasm")

// Test deployment
testContractDeployment(t, counterWasm)
testTokenOperations(t, cw20Wasm)
testComplexLogic(t, nameserviceWasm)
```

## ✅ Testing Coverage

These contracts enable comprehensive testing:

- ✅ **Contract deployment** (store-code)
- ✅ **Contract instantiation** with different init messages
- ✅ **Contract execution** with various execute messages
- ✅ **Contract queries** with different query types
- ✅ **Multi-contract interactions**
- ✅ **Payment handling** (nameservice)
- ✅ **State management** (counter)
- ✅ **Token operations** (CW20)

## 📚 Contract Documentation

Each contract has detailed documentation:

- **[Counter Contract README](./counter-contract/README.md)** - Basic state management
- **[CW20 Token README](./cw20-token-contract/README.md)** - Fungible tokens
- **[Name Service README](./nameservice-contract/README.md)** - Name registration

## 🔧 Contract Features Comparison

| Feature | Counter | CW20 Token | Name Service |
|---------|---------|------------|--------------|
| State Management | ✅ Simple | ✅ Complex | ✅ Complex |
| Payment Handling | ❌ | ❌ | ✅ |
| Access Control | ❌ | ✅ Minter | ✅ Owner |
| Multiple Users | ✅ | ✅ | ✅ |
| Tests Included | ✅ | ✅ | ✅ |
| Production Ready | ✅ | ✅ | ✅ |

## 🎯 Use Cases

### Counter Contract
- Learning CosmWasm basics
- Testing contract lifecycle
- Simple state operations

### CW20 Token
- Creating fungible tokens
- Token transfers and burns
- DeFi applications
- Governance tokens

### Name Service
- Domain name systems
- Identity management
- Payment-based services
- Ownership transfers

## 📦 Dependencies

All contracts use standard CosmWasm dependencies:

- **cosmwasm-std**: ^1.5 - Core CosmWasm library
- **cosmwasm-schema**: ^1.5 - Schema generation
- **cw-storage-plus**: ^1.2 - Enhanced storage
- **serde**: ^1.0 - Serialization
- **thiserror**: ^1.0 - Error handling

CW20 Token additionally uses:
- **cw2**: ^1.1 - Contract versioning
- **cw20**: ^1.1 - CW20 standard

## 🔐 Security

All contracts include:
- Comprehensive unit tests
- Error handling
- Input validation
- Safe arithmetic (overflow protection)

For production use:
- Audit the code
- Test on testnet first
- Use optimized builds
- Follow CosmWasm best practices

## 📖 Additional Resources

- **CosmWasm Documentation**: https://docs.cosmwasm.com/
- **CosmWasm Book**: https://book.cosmwasm.com/
- **CosmWasm Plus**: https://github.com/CosmWasm/cw-plus
- **CosmWasm Examples**: https://github.com/CosmWasm/cosmwasm-examples
- **Rust Book**: https://doc.rust-lang.org/book/

## 🚢 Production Deployment

1. **Build optimized WASM**:
```bash
docker run --rm -v "$(pwd)":/code \
  cosmwasm/optimizer:0.15.0
```

2. **Store on chain**:
```bash
hippod tx wasm store contract.wasm --from <key> --gas 2000000
```

3. **Instantiate**:
```bash
hippod tx wasm instantiate <code-id> '<init-msg>' --from <key> --label "my-contract"
```

4. **Interact**:
```bash
hippod tx wasm execute <contract-addr> '<execute-msg>' --from <key>
hippod query wasm contract-state smart <contract-addr> '<query-msg>'
```

## 📝 License

All contracts are licensed under Apache 2.0
