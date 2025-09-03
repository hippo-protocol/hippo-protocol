# 🦛 hippo-sdk

`hippo-sdk` is a WebAssembly-based npm library for cryptographic operations, including **key pair generation**, **ECDH encryption**, **ECDSA signing**, and **DID management**. It's built with Rust for high performance and compiled to WebAssembly for seamless use in web applications.

---

## 📦 Installation

To add `hippo-sdk` to your project, use either `npm` or `yarn`:

```bash
npm install hippo-sdk
```

## 🚀 Usage

This library provides a set of cryptographic utilities. You can import the functions directly and use them in your JavaScript or TypeScript projects.

### Example

Here's a simple example demonstrating how to generate a key pair and then use it to encrypt and decrypt data.

```javascript
import {
  create_keypair,
  key_to_did,
  encrypt,
  decrypt,
  sign,
  verify,
} from "hippo-sdk";

// 1. Generate a new key pair
const keypair = create_keypair();
const { pubkey, privkey } = keypair.to_object();

console.log("Public Key:", pubkey);
console.log("Private Key:", privkey);

// 2. Create a DID from the public key
const myDid = key_to_did(pubkey);

console.log("DID:", myDid.to_object());
// Expected output: { "id": "did:hp:<pubkey>" }

// 3. Encrypt some data using the public key
const rawData = "Hello, wasm!";
const encryptedData = encrypt(rawData, pubkey);

console.log("Encrypted Data:", encryptedData.to_object());

// 4. Decrypt the data using the private key
const decryptedData = decrypt(encryptedData, privkey);

console.log("Decrypted Data:", decryptedData);
// Expected output: "Hello, wasm!"

// 5. Sign a message using the private key
const messageToSign = "This is a message to be signed.";
const signature = sign(messageToSign, privkey);

console.log("Signature:", signature);

// 6. Verify the signature using the public key
const isVerified = verify(messageToSign, signature, pubkey);

console.log("Signature Verified:", isVerified);
// Expected output: true
```

## 📖 API

The library exposes several functions for various cryptographic tasks.

### `create_keypair(): KeyPair`

Generates a new **secp256k1** key pair. Returns a `KeyPair` object containing the public and private keys.

- **Returns**: `KeyPair`

### `key_to_did(pubkey: string): Did`

Converts a public key string into a DID (Decentralized Identifier) with the format `did:hp:<pubkey>`.

- **`pubkey`**: The public key string.
- **Returns**: `Did`

### `did_to_key(did: Did): string`

Extracts the public key string from a DID.

- **`did`**: The `Did` object.
- **Returns**: `string`

### `encrypt(data: string, pubkey: string): EncryptedData`

Encrypts a string of data using **ECDH** key exchange and **AES-256-GCM** encryption.

- **`data`**: The string to be encrypted.
- **`pubkey`**: The public key of the recipient.
- **Returns**: `EncryptedData`

### `decrypt(data: EncryptedData, privkey: string): string`

Decrypts data that was encrypted using the `encrypt` function.

- **`data`**: The `EncryptedData` object.
- **`privkey`**: The private key of the recipient.
- **Returns**: `string`

### `sign(data: string, privkey: string): string`

Signs a string of data using **ECDSA**.

- **`data`**: The string to be signed.
- **`privkey`**: The private key of the signer.
- **Returns**: The signature string.

### `verify(data: string, sig: string, pubkey: string): boolean`

Verifies an **ECDSA** signature against the original data and a public key.

- **`data`**: The original string of data.
- **`sig`**: The signature string.
- **`pubkey`**: The public key of the signer.
- **Returns**: `boolean` indicating if the signature is valid.

### `sha256(data: string): string`

Computes the **SHA-256** hash of a string.

- **`data`**: The string to hash.
- **Returns**: The hash as a string.

### `ecdh(privkey: string, pubkey: string): string`

Performs an **Elliptic Curve Diffie-Hellman (ECDH)** key exchange to compute a shared secret.

- **`privkey`**: Your private key.
- **`pubkey`**: The other party's public key.
- **Returns**: The shared secret as a hex string.

---

## 🤝 Contributing

We welcome contributions! Feel free to open an issue or submit a pull request on the [GitHub repository](https://github.com/hippo-protocol/hippo-protocol).

---

## 📄 License

This project is licensed under the [MIT License](https://github.com/hippo-protocol/hippo-protocol/blob/main/LICENSE).
