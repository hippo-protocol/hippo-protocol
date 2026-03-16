# CW20 Token Contract - Rust Source

A CW20-compatible fungible token smart contract for creating and managing tokens on CosmWasm.

## Features

- Create fungible tokens with custom name, symbol, and decimals
- Transfer tokens between addresses
- Burn tokens
- Mint new tokens (with minter permission)
- Query balances and token info
- Optional minter with supply cap

## Project Structure

```
cw20-token-contract/
├── Cargo.toml          # Rust package configuration
└── src/
    ├── lib.rs          # Library entry point
    ├── contract.rs     # Main contract logic (transfer, mint, burn)
    ├── msg.rs          # Message definitions
    ├── state.rs        # State management (balances, token info)
    └── error.rs        # Error types
```

## Messages

### InstantiateMsg
```rust
pub struct InstantiateMsg {
    pub name: String,
    pub symbol: String,
    pub decimals: u8,
    pub initial_balances: Vec<Cw20Coin>,
    pub mint: Option<MinterResponse>,
    pub marketing: Option<InstantiateMarketingInfo>,
}
```

### ExecuteMsg
```rust
pub enum ExecuteMsg {
    Transfer { recipient: String, amount: Uint128 },
    Burn { amount: Uint128 },
    Send { contract: String, amount: Uint128, msg: String },
    Mint { recipient: String, amount: Uint128 },
}
```

### QueryMsg
```rust
pub enum QueryMsg {
    Balance { address: String },
    TokenInfo {},
    Minter {},
}
```

## Building

### Prerequisites

1. Install Rust: https://rustup.rs/
2. Add wasm32 target:
```bash
rustup target add wasm32-unknown-unknown
```

### Build Commands

#### Development Build
```bash
cd cw20-token-contract
cargo build
```

#### Run Tests
```bash
cargo test
```

#### Production Build
```bash
cargo build --release --target wasm32-unknown-unknown
```

Output: `target/wasm32-unknown-unknown/release/cw20_token_contract.wasm`

#### Optimize for Production
```bash
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/optimizer:0.15.0
```

Optimized: `artifacts/cw20_token_contract.wasm`

## Deployment

### Store Code
```bash
hippod tx wasm store cw20_token_contract.wasm \
  --from <your-key> \
  --gas 2000000 \
  --fees 1000000000000000000ahp \
  -y
```

### Instantiate Token
```bash
hippod tx wasm instantiate <code-id> '{
  "name": "My Token",
  "symbol": "MTK",
  "decimals": 6,
  "initial_balances": [
    {
      "address": "hippo1...",
      "amount": "1000000000"
    }
  ],
  "mint": {
    "minter": "hippo1...",
    "cap": null
  }
}' \
  --from <your-key> \
  --label "my-token" \
  --gas 500000 \
  --fees 1000000000000000000ahp \
  --no-admin \
  -y
```

### Transfer Tokens
```bash
hippod tx wasm execute <contract-addr> '{
  "transfer": {
    "recipient": "hippo1...",
    "amount": "1000000"
  }
}' \
  --from <your-key> \
  --gas 300000 \
  --fees 1000000000000000000ahp \
  -y
```

### Burn Tokens
```bash
hippod tx wasm execute <contract-addr> '{
  "burn": {
    "amount": "500000"
  }
}' \
  --from <your-key> \
  --gas 300000 \
  --fees 1000000000000000000ahp \
  -y
```

### Mint Tokens (if authorized)
```bash
hippod tx wasm execute <contract-addr> '{
  "mint": {
    "recipient": "hippo1...",
    "amount": "1000000"
  }
}' \
  --from <minter-key> \
  --gas 300000 \
  --fees 1000000000000000000ahp \
  -y
```

### Query Balance
```bash
hippod query wasm contract-state smart <contract-addr> '{
  "balance": {
    "address": "hippo1..."
  }
}'
```

### Query Token Info
```bash
hippod query wasm contract-state smart <contract-addr> '{"token_info":{}}'
```

## Testing

Run comprehensive unit tests:

```bash
cargo test
```

Tests include:
- Token initialization
- Transfer functionality
- Balance tracking
- Burn operations

## Dependencies

- `cosmwasm-std`: ^1.5 - CosmWasm standard library
- `cw-storage-plus`: ^1.2 - Storage helpers
- `cw2`: ^1.1 - Contract versioning
- `cw20`: ^1.1 - CW20 standard
- `serde`: ^1.0 - Serialization

## CW20 Standard Compliance

This contract implements the CW20 fungible token standard:
- Standard transfer/burn/mint operations
- Balance and supply queries
- Optional minter with cap
- Compatible with CW20 ecosystem

## Security Considerations

- Minter authorization checked on mint
- Cap enforcement prevents unlimited minting
- Zero amount transfers/burns rejected
- Overflow protection on all arithmetic
- Address validation on all operations

## Production Checklist

- [ ] Set appropriate minter (or None for fixed supply)
- [ ] Configure supply cap if needed
- [ ] Test with small amounts first
- [ ] Verify all addresses before instantiation
- [ ] Use optimized WASM build
- [ ] Audit code if handling significant value

## License

Apache 2.0
