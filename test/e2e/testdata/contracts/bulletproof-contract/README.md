# Bulletproof Contract

## Overview

A CosmWasm smart contract that stores and verifies Bulletproof range proofs on-chain. Bulletproofs are short non-interactive zero-knowledge proofs that require no trusted setup.

## Features

- **Store Bulletproof Proof**: Stores a serialized bulletproof range proof and its SHA-256 hash
- **On-chain Verification**: Verifies the stored bulletproof range proof using the `bulletproofs` library
- **Gas Measurement**: Enables measuring gas consumption for on-chain bulletproof verification

## Messages

### InstantiateMsg
```json
{
  "proof_hex": "<hex-encoded proof bytes>",
  "commitment_hex": "<hex-encoded commitment, 32 bytes>",
  "num_bits": 32
}
```

### ExecuteMsg
```json
{"verify": {}}
```

### QueryMsg
```json
{"get_proof_hash": {}}
```
```json
{"verify": {}}
```

## Dependencies

- `bulletproofs` v4.0.0 - Bulletproof range proof implementation
- `curve25519-dalek-ng` v4.1.1 - Elliptic curve operations
- `merlin` v3 - Transcript-based proof protocol
- `cosmwasm-std` v1.5 - CosmWasm standard library

## Building

```bash
cargo build --target wasm32-unknown-unknown --release
```

## Proof Generation

The bulletproof proof is generated off-chain using the standard `bulletproofs` library:

```rust
use bulletproofs::{BulletproofGens, PedersenGens, RangeProof};
use curve25519_dalek_ng::scalar::Scalar;
use merlin::Transcript;
use rand::thread_rng;

let pc_gens = PedersenGens::default();
let bp_gens = BulletproofGens::new(64, 1);
let secret_value = 1037578891u64;
let mut rng = thread_rng();
let blinding = Scalar::random(&mut rng);
let mut prover_transcript = Transcript::new(b"doctest example");

let (proof, committed_value) = RangeProof::prove_single_with_rng(
    &bp_gens, &pc_gens, &mut prover_transcript,
    secret_value, &blinding, 32, &mut rng,
).expect("proof generation failed");
```

The serialized proof and commitment are then passed to the contract during instantiation.
