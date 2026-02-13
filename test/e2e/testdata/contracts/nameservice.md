# Name Service Contract

A decentralized name service smart contract for registering and managing names.

## Features

- Register names with payment
- Transfer name ownership
- Resolve names to addresses
- Query name configuration
- Set purchase and transfer prices

## Messages

### InstantiateMsg
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

### ExecuteMsg

Register a name:
```json
{
  "register": {
    "name": "alice"
  }
}
```

Transfer a name:
```json
{
  "transfer": {
    "name": "alice",
    "to": "hippo1..."
  }
}
```

### QueryMsg

Resolve a name:
```json
{
  "resolve_record": {
    "name": "alice"
  }
}
```

Response:
```json
{
  "address": "hippo1..."
}
```

Get configuration:
```json
{
  "config": {}
}
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

## Usage in Tests

This contract is used to test:
- Complex contract initialization with multiple parameters
- Payment handling with contract execution
- Name-to-address mapping functionality
- Ownership transfer logic
- Query operations with parameters

## Contract Size

Approximately 1.1KB compiled
