# Name Service Contract - Rust Source

A decentralized name service smart contract for registering and managing human-readable names on CosmWasm.

## Features

- Register names with payment
- Transfer name ownership
- Resolve names to addresses
- Configurable purchase and transfer prices
- Name validation (length, characters)
- Payment verification

## Project Structure

```
nameservice-contract/
├── Cargo.toml          # Rust package configuration
└── src/
    ├── lib.rs          # Library entry point
    ├── contract.rs     # Main contract logic (register, transfer, resolve)
    ├── msg.rs          # Message definitions
    ├── state.rs        # State management (names, config)
    └── error.rs        # Error types
```

## Messages

### InstantiateMsg
```rust
pub struct InstantiateMsg {
    pub purchase_price: Option<Coin>,
    pub transfer_price: Option<Coin>,
}
```

### ExecuteMsg
```rust
pub enum ExecuteMsg {
    Register { name: String },
    Transfer { name: String, to: String },
}
```

### QueryMsg
```rust
pub enum QueryMsg {
    ResolveRecord { name: String },
    Config {},
}
```

## Name Validation Rules

- **Minimum length**: 3 characters
- **Maximum length**: 64 characters
- **Allowed characters**: alphanumeric (a-z, A-Z, 0-9) and hyphens (-)
- **Case sensitive**: "Alice" and "alice" are different names

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
cd nameservice-contract
cargo build
```

#### Run Tests
```bash
cargo test
```

Tests include:
- Name registration
- Name transfers
- Payment verification
- Name validation
- Ownership checks

#### Production Build
```bash
cargo build --release --target wasm32-unknown-unknown
```

Output: `target/wasm32-unknown-unknown/release/nameservice_contract.wasm`

#### Optimize for Production
```bash
docker run --rm -v "$(pwd)":/code \
  --mount type=volume,source="$(basename "$(pwd)")_cache",target=/target \
  --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
  cosmwasm/optimizer:0.15.0
```

Optimized: `artifacts/nameservice_contract.wasm`

## Deployment

### Store Code
```bash
hippod tx wasm store nameservice_contract.wasm \
  --from <your-key> \
  --gas 2000000 \
  --fees 1000000000000000000ahp \
  -y
```

### Instantiate Service
```bash
hippod tx wasm instantiate <code-id> '{
  "purchase_price": {
    "denom": "ahp",
    "amount": "1000000000000000000"
  },
  "transfer_price": {
    "denom": "ahp",
    "amount": "500000000000000000"
  }
}' \
  --from <your-key> \
  --label "nameservice" \
  --gas 500000 \
  --fees 1000000000000000000ahp \
  --no-admin \
  -y
```

### Register a Name
```bash
hippod tx wasm execute <contract-addr> '{
  "register": {
    "name": "alice"
  }
}' \
  --from <your-key> \
  --amount 1000000000000000000ahp \
  --gas 300000 \
  --fees 1000000000000000000ahp \
  -y
```

### Transfer Name Ownership
```bash
hippod tx wasm execute <contract-addr> '{
  "transfer": {
    "name": "alice",
    "to": "hippo1..."
  }
}' \
  --from <owner-key> \
  --amount 500000000000000000ahp \
  --gas 300000 \
  --fees 1000000000000000000ahp \
  -y
```

### Resolve Name to Address
```bash
hippod query wasm contract-state smart <contract-addr> '{
  "resolve_record": {
    "name": "alice"
  }
}'
```

Response:
```json
{
  "address": "hippo1..."
}
```

### Query Configuration
```bash
hippod query wasm contract-state smart <contract-addr> '{"config":{}}'
```

Response:
```json
{
  "purchase_price": {
    "denom": "ahp",
    "amount": "1000000000000000000"
  },
  "transfer_price": {
    "denom": "ahp",
    "amount": "500000000000000000"
  }
}
```

## Use Cases

### Personal Name Service
Users can register their names for easy identification:
```bash
# Alice registers her name
hippod tx wasm execute <contract> '{"register":{"name":"alice"}}' --from alice --amount 1000000000000000000ahp

# Others can resolve alice's address
hippod query wasm contract-state smart <contract> '{"resolve_record":{"name":"alice"}}'
```

### Domain Names
Register domains for dApps or services:
```bash
# Register a domain
hippod tx wasm execute <contract> '{"register":{"name":"my-dapp"}}' --amount 1000000000000000000ahp

# Transfer to organization
hippod tx wasm execute <contract> '{"transfer":{"name":"my-dapp","to":"hippo1org..."}}' --amount 500000000000000000ahp
```

### Business Names
Companies can register their brand names:
```bash
hippod tx wasm execute <contract> '{"register":{"name":"acme-corp"}}' --amount 1000000000000000000ahp
```

## Testing

Run all tests:
```bash
cargo test
```

Run specific test:
```bash
cargo test test_name -- --nocapture
```

Tests cover:
- ✅ Name registration with payment
- ✅ Name transfer with payment
- ✅ Name resolution
- ✅ Ownership verification
- ✅ Payment verification
- ✅ Name validation (too short, too long, invalid chars)
- ✅ Duplicate name prevention

## Dependencies

- `cosmwasm-std`: ^1.5 - CosmWasm standard library
- `cw-storage-plus`: ^1.2 - Storage helpers
- `serde`: ^1.0 - Serialization
- `thiserror`: ^1.0 - Error handling

## Error Handling

The contract includes comprehensive error types:

- `NameTaken`: Name already registered
- `NameNotExists`: Name not found
- `Unauthorized`: Not the name owner
- `InsufficientFundsSend`: Payment too low
- `NameTooShort`: Less than 3 characters
- `NameTooLong`: More than 64 characters
- `InvalidCharacter`: Non-alphanumeric (except hyphens)

## Security Considerations

- **Payment verification**: All registrations and transfers require correct payment
- **Ownership checks**: Only owners can transfer names
- **Name validation**: Prevents invalid names
- **No admin override**: Once registered, only owner controls name
- **Immutable history**: Name ownership changes are transparent

## Production Deployment Checklist

- [ ] Set appropriate prices for your use case
- [ ] Consider price in smallest denomination (e.g., 1 ahp = 10^-18 HP)
- [ ] Test registration flow on testnet
- [ ] Test transfer flow on testnet
- [ ] Verify name validation rules
- [ ] Use optimized WASM build
- [ ] Document price structure for users

## Future Enhancements

Potential additions for advanced versions:
- Name expiration/renewal
- Sub-names (e.g., alice.hippo)
- Name metadata (avatar, description)
- Auction system for premium names
- Bulk name operations

## License

Apache 2.0
