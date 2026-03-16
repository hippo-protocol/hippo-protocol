# Counter Contract - Rust Source

A simple counter smart contract demonstrating basic CosmWasm functionality.

## Features

- Initialize with a starting count value
- Increment the counter
- Reset the counter to a new value
- Query the current count
- Owner tracking

## Project Structure

```
counter-contract/
├── Cargo.toml          # Rust package configuration
└── src/
    ├── lib.rs          # Library entry point
    ├── contract.rs     # Main contract logic
    ├── msg.rs          # Message definitions
    ├── state.rs        # State management
    └── error.rs        # Error types
```

## Messages

### InstantiateMsg
```rust
pub struct InstantiateMsg {
    pub count: i32,
}
```

### ExecuteMsg
```rust
pub enum ExecuteMsg {
    Increment {},
    Reset { count: i32 },
}
```

### QueryMsg
```rust
pub enum QueryMsg {
    GetCount {},
}
```

## Building

### Prerequisites

1. Install Rust: https://rustup.rs/
2. Install wasm32 target:
```bash
rustup target add wasm32-unknown-unknown
```

3. Install cargo-generate (optional, for new projects):
```bash
cargo install cargo-generate --features vendored-openssl
```

### Build Commands

#### Development Build
```bash
cd counter-contract
cargo build
```

#### Run Tests
```bash
cargo test
```

#### Production Build (Optimized WASM)
```bash
cargo build --release --target wasm32-unknown-unknown
```

The compiled contract will be at:
`target/wasm32-unknown-unknown/release/counter_contract.wasm`

#### Optimize WASM (Recommended for Production)

Install optimizer:
```bash
docker pull cosmwasm/optimizer:0.15.0
```

Run optimizer:
```bash
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/optimizer:0.15.0
```

Optimized contract will be in `artifacts/counter_contract.wasm`

## Deployment

### Using CLI

Store code:
```bash
hippod tx wasm store counter_contract.wasm \
  --from <your-key> \
  --gas 2000000 \
  --fees 1000000000000000000ahp \
  -y
```

Instantiate:
```bash
hippod tx wasm instantiate <code-id> '{"count":0}' \
  --from <your-key> \
  --label "counter" \
  --gas 500000 \
  --fees 1000000000000000000ahp \
  --no-admin \
  -y
```

Execute - Increment:
```bash
hippod tx wasm execute <contract-addr> '{"increment":{}}' \
  --from <your-key> \
  --gas 300000 \
  --fees 1000000000000000000ahp \
  -y
```

Execute - Reset:
```bash
hippod tx wasm execute <contract-addr> '{"reset":{"count":42}}' \
  --from <your-key> \
  --gas 300000 \
  --fees 1000000000000000000ahp \
  -y
```

Query:
```bash
hippod query wasm contract-state smart <contract-addr> '{"get_count":{}}'
```

## Testing

The contract includes comprehensive unit tests:

```bash
cargo test
```

Test output will show:
- Proper initialization
- Increment functionality
- Reset functionality

## Dependencies

- `cosmwasm-std`: ^1.5 - CosmWasm standard library
- `cosmwasm-schema`: ^1.5 - Schema generation
- `cw-storage-plus`: ^1.2 - Enhanced storage helpers
- `serde`: ^1.0 - Serialization framework
- `thiserror`: ^1.0 - Error handling

## Security Considerations

- No access control on increment/reset (anyone can call)
- Owner is tracked but not enforced
- For production, add proper authorization checks

## License

Apache 2.0
