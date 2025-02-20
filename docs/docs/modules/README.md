---
title: SDK
order: 1
---

Here you can find an overview of the modules included on the Hippo Protocol blockchain with relevant info and
links for each one.

## hippo-sdk

[hippo-sdk](https://github.com/hippocrat-dao/hippo-protocol/tree/main/sdk/core) provides the simplest way to utilize the data protocol.

```rust
pub trait Sdk {
    /// chain
    fn write(data: String, privkey: String) -> Tx;
    fn read(tx_hash: String) -> Tx;
    fn send(from: String, to: String, privkey: String, data: Option<String>) -> Tx;
    /// did
    fn create_keypair() -> KeyPair;
    fn key_to_did(pubkey: String) -> Did;
    fn did_to_key(did: Did) -> String;
    /// cryptography
    fn encrypt(data: String, pubkey: String) -> EncryptedData;
    fn decrypt(data: EncryptedData, privkey: String) -> String;
    fn sign(data: String, privkey: String) -> String;
    fn verify(data: String, sig: String, pubkey: String) -> bool;
    fn sha256(data: String) -> String;
    fn ecdh(privkey: String, pubkey: String) -> String;
}
```

## cosmos-sdk

- [x/auth](https://docs.cosmos.network/v0.47/build/modules/auth)
- [x/authz](https://docs.cosmos.network/v0.47/build/modules/authz)
- [x/bank](https://docs.cosmos.network/v0.47/build/modules/bank)
- [x/capability](https://docs.cosmos.network/v0.47/build/modules/capability)
- [x/consensus](https://docs.cosmos.network/v0.47/build/modules/consensus)
- [x/crisis](https://docs.cosmos.network/v0.47/build/modules/crisis)
- [x/distribution](https://docs.cosmos.network/v0.47/build/modules/distribution)
- [x/evidence](https://docs.cosmos.network/v0.47/build/modules/evidence)
- [x/feegrant](https://docs.cosmos.network/v0.47/build/modules/feegrant)
- [x/gov](https://docs.cosmos.network/v0.47/build/modules/gov)
- [x/mint](https://docs.cosmos.network/v0.47/build/modules/mint)
- [x/params](https://docs.cosmos.network/v0.47/build/modules/params)
- [x/slashing](https://docs.cosmos.network/v0.47/build/modules/slashing)
- [x/staking](https://docs.cosmos.network/v0.47/build/modules/staking)
- [x/upgrade](https://docs.cosmos.network/v0.47/build/modules/upgrade)
