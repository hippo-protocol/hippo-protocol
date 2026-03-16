# CW20 Token Contract

A CW20-compatible fungible token smart contract for testing token operations.

## Features

- Create a new fungible token with name, symbol, and decimals
- Transfer tokens between addresses
- Query token balance
- Query token metadata
- Burn tokens
- Mint tokens (if authorized)

## Messages

### InstantiateMsg
```json
{
  "name": "My Token",
  "symbol": "MTK",
  "decimals": 6,
  "initial_balances": [
    {
      "address": "hippo1...",
      "amount": "1000000"
    }
  ]
}
```

### ExecuteMsg

Transfer tokens:
```json
{
  "transfer": {
    "recipient": "hippo1...",
    "amount": "100000"
  }
}
```

Burn tokens:
```json
{
  "burn": {
    "amount": "50000"
  }
}
```

Mint tokens (if authorized):
```json
{
  "mint": {
    "recipient": "hippo1...",
    "amount": "100000"
  }
}
```

### QueryMsg

Get balance:
```json
{
  "balance": {
    "address": "hippo1..."
  }
}
```

Get token info:
```json
{
  "token_info": {}
}
```

Response:
```json
{
  "name": "My Token",
  "symbol": "MTK",
  "decimals": 6,
  "total_supply": "1000000"
}
```

## Usage in Tests

This contract is used to test:
- Token creation and initialization
- Token transfers
- Balance queries
- Multi-user token interactions
- Token burn/mint operations

## Contract Size

Approximately 993 bytes compiled
