# Counter Contract

A simple counter smart contract for testing CosmWasm functionality.

## Features

- Initialize with a starting count value
- Increment the counter
- Reset the counter to a new value
- Query the current count

## Messages

### InstantiateMsg
```json
{
  "count": 0
}
```

### ExecuteMsg

Increment the counter:
```json
{
  "increment": {}
}
```

Reset the counter:
```json
{
  "reset": {
    "count": 42
  }
}
```

### QueryMsg

Get current count:
```json
{
  "get_count": {}
}
```

Response:
```json
{
  "count": 42
}
```

## Usage in Tests

This contract is used to test:
- Contract deployment (store-code)
- Contract instantiation
- Contract execution (increment/reset)
- Contract queries (get_count)

## Contract Size

Approximately 975 bytes compiled
