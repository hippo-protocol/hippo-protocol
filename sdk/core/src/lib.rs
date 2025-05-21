#[cfg(test)]
mod tests;
mod types;
use std::str::FromStr;

use aes_gcm::{
    aead::{Aead, AeadCore, OsRng as AesRng},
    Aes256Gcm, Key, KeyInit, Nonce,
};
use secp256k1::hashes::{hex::FromHex, sha256, Hash};
use secp256k1::rand::rngs::OsRng as Secp256k1Rng;
use secp256k1::{ecdh, Message, PublicKey, Secp256k1, SecretKey};
use secp256k1::{ecdsa::Signature, hashes::hex::DisplayHex};
use wasm_bindgen::prelude::*;

use types::{Did, EncryptedData, KeyPair};

#[wasm_bindgen]
pub fn create_keypair() -> KeyPair {
    let secp = Secp256k1::new();
    let (secret_key, public_key) = secp.generate_keypair(&mut Secp256k1Rng);
    KeyPair::new(
        public_key.to_string(),
        secret_key.display_secret().to_string(),
    )
}

#[wasm_bindgen]
pub fn key_to_did(pubkey: String) -> Did {
    Did::new("did:hp:".to_owned() + &pubkey)
}

#[wasm_bindgen]
pub fn did_to_key(did: Did) -> String {
    did.id().strip_prefix("did:hp:").unwrap().to_string()
}

#[wasm_bindgen]
pub fn encrypt(data: String, pubkey: String) -> EncryptedData {
    let secp = Secp256k1::new();
    // Alice: one-off key pair to caculate shared secret.
    let (secret_key, public_key) = secp.generate_keypair(&mut Secp256k1Rng);
    // Bob: param pubkey is the one who's able to decrypt the data.
    let encrypt_to_pubkey = secp256k1::PublicKey::from_str(&pubkey).unwrap();
    // Only Alice and Bob can know the secret, which is used as a key to encrypt data.
    let shared_secret = ecdh::SharedSecret::new(&encrypt_to_pubkey, &secret_key).secret_bytes();

    // Secret key just fit into AES key
    let aes_key = Key::<Aes256Gcm>::from_slice(&shared_secret);
    let cipher = Aes256Gcm::new(&aes_key);
    // Nonce must be random and non-reusable value
    let nonce = Aes256Gcm::generate_nonce(&mut AesRng); // 96-bits; unique per message
    let ciphertext = cipher.encrypt(&nonce, data.as_bytes()).unwrap();
    EncryptedData::new(
        public_key.to_string(),
        pubkey,
        ciphertext.to_lower_hex_string(),
        nonce.to_lower_hex_string(),
    )
}
#[wasm_bindgen]
pub fn decrypt(data: EncryptedData, privkey: String) -> String {
    // Alice: one-off pubkey to calculate shared secret.
    let encrypt_from_pubkey = secp256k1::PublicKey::from_str(&data.pubkey_from()).unwrap();
    // Bob: the one who's able to decrypt the data with privkey.
    let privkey_to_decrypt = secp256k1::SecretKey::from_str(&privkey).unwrap();
    // Only Alice and Bob can know the secret, which is used as a key to encrypt data.
    let shared_secret =
        ecdh::SharedSecret::new(&encrypt_from_pubkey, &privkey_to_decrypt).secret_bytes();

    // Secret key just fit into AES key
    let aes_key = Key::<Aes256Gcm>::from_slice(&shared_secret);
    let cipher = Aes256Gcm::new(&aes_key);
    // Nonce hex to bytes
    let nonce_bytes = Vec::from_hex(&data.nonce()).unwrap();
    let nonce = Nonce::from_slice(&nonce_bytes);
    // Data hex to bytes
    let data_bytes = Vec::from_hex(&data.data()).unwrap();
    let decrypted_data = cipher.decrypt(&nonce, data_bytes.as_slice()).unwrap();

    String::from_utf8(decrypted_data).unwrap()
}

#[wasm_bindgen]
pub fn sign(data: String, privkey: String) -> String {
    let secp = Secp256k1::new();
    let digest = sha256::Hash::hash(data.as_bytes());
    let message = Message::from_digest(digest.to_byte_array());
    let secret_key = SecretKey::from_str(&privkey).unwrap();

    secp.sign_ecdsa(&message, &secret_key).to_string()
}
#[wasm_bindgen]
pub fn verify(data: String, sig: String, pubkey: String) -> bool {
    let secp = Secp256k1::new();
    let digest = sha256::Hash::hash(data.as_bytes());
    let message = Message::from_digest(digest.to_byte_array());
    let signature = Signature::from_str(&sig).unwrap();
    let public_key = PublicKey::from_str(&pubkey).unwrap();

    match secp.verify_ecdsa(&message, &signature, &public_key) {
        Ok(_) => true,
        Err(_) => false,
    }
}
#[wasm_bindgen]
pub fn sha256(data: String) -> String {
    sha256::Hash::hash(data.as_bytes()).to_string()
}
#[wasm_bindgen]
pub fn ecdh(privkey: String, pubkey: String) -> String {
    // Alice: param privkey is the one who's able to decrypt the data.
    let alice = secp256k1::SecretKey::from_str(&privkey).unwrap();
    // Bob: param pubkey is the one who's able to decrypt the data.
    let bob = secp256k1::PublicKey::from_str(&pubkey).unwrap();
    // Only Alice and Bob can know the secret, which is used as a key to sign or encrypt data(or any other).
    let shared_secret = ecdh::SharedSecret::new(&bob, &alice);
    shared_secret.display_secret().to_string()
}
