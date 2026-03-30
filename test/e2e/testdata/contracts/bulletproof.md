# Bulletproof Range Proof Verification Contract

| Field | Value |
|-------|-------|
| **Name** | Bulletproof Verification |
| **Version** | 1.0.0 |
| **Binary** | bulletproof.wasm |
| **Source** | bulletproof-contract/ |
| **Size** | ~339 KB |
| **CosmWasm Version** | 1.5 |

## Description

This contract stores a Bulletproof range proof and its SHA-256 hash, and provides on-chain verification using the `bulletproofs` cryptographic library. Bulletproofs are short non-interactive zero-knowledge proofs that require no trusted setup, particularly useful for range proofs.

## Use Case

- Zero-knowledge range proof verification on-chain
- Proving a committed value lies within a specific range without revealing the value
- Gas consumption benchmarking for cryptographic verification operations

## Key Dependencies

- `bulletproofs` v4.0.0 (zero-knowledge range proofs)
- `curve25519-dalek-ng` v4.1.1 (Ristretto group operations)
- `merlin` v3 (Fiat-Shamir transcript)
