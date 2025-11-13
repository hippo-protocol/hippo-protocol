#[cfg(test)]
mod tests;
mod types;
use std::str::FromStr;

use aes_gcm::{
    aead::{Aead, AeadCore, OsRng as AesRng},
    Aes256Gcm, Key, KeyInit, Nonce,
};
use secp256k1_zkp::{
    ecdh,
    rand::{self, rngs::OsRng as Secp256k1Rng},
    verify_commitments_sum_to_equal, Generator, Message, PedersenCommitment, PublicKey, Secp256k1,
    SecretKey,
};
use secp256k1_zkp::{ecdsa::Signature, Tweak};
use secp256k1_zkp::{
    hashes::{sha256, Hash},
    Tag,
};
use wasm_bindgen::prelude::*;

use types::{AesEncryptedData, Commitment, Did, EncodingType, EncryptedData, KeyPair};

use base64::{engine::general_purpose::STANDARD, Engine as _};

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
pub fn encrypt(data: String, pubkey: String, encoding_type: EncodingType) -> EncryptedData {
    let secp = Secp256k1::new();
    // Alice: one-off key pair to caculate shared secret.
    let (secret_key, public_key) = secp.generate_keypair(&mut Secp256k1Rng);
    // Bob: param pubkey is the one who's able to decrypt the data.
    let encrypt_to_pubkey = secp256k1_zkp::PublicKey::from_str(&pubkey).unwrap();
    // Only Alice and Bob can know the secret, which is used as a key to encrypt data.
    let shared_secret = ecdh::SharedSecret::new(&encrypt_to_pubkey, &secret_key).secret_bytes();

    let aes_encrypted_data = encrypt_aes(data, hex::encode(&shared_secret), encoding_type);

    EncryptedData::new(
        public_key.to_string(),
        pubkey,
        aes_encrypted_data.data(),
        aes_encrypted_data.nonce(),
    )
}

#[wasm_bindgen]
pub fn decrypt(data: EncryptedData, privkey: String, encoding_type: EncodingType) -> String {
    // Alice: one-off pubkey to calculate shared secret.
    let encrypt_from_pubkey = secp256k1_zkp::PublicKey::from_str(&data.pubkey_from()).unwrap();
    // Bob: the one who's able to decrypt the data with privkey.
    let privkey_to_decrypt = secp256k1_zkp::SecretKey::from_str(&privkey).unwrap();
    // Only Alice and Bob can know the secret, which is used as a key to encrypt data.
    let shared_secret =
        ecdh::SharedSecret::new(&encrypt_from_pubkey, &privkey_to_decrypt).secret_bytes();

    decrypt_aes(
        AesEncryptedData::new(data.data(), data.nonce()),
        hex::encode(&shared_secret),
        encoding_type,
    )
}

#[wasm_bindgen]
pub fn encrypt_aes(data: String, key: String, encoding_type: EncodingType) -> AesEncryptedData {
    // Secret key just fit into AES key
    let hex_key = hex::decode(key).expect("Key is malformed.");
    let aes_key = Key::<Aes256Gcm>::from_slice(&hex_key);
    let cipher = Aes256Gcm::new(&aes_key);
    // Nonce must be random and non-reusable value
    let nonce = Aes256Gcm::generate_nonce(&mut AesRng); // 96-bits; unique per message
    let data_bytes = match encoding_type {
        EncodingType::UTF8 => data.as_bytes(),
        EncodingType::HEX => &hex::decode(data).expect("Wrong hex format data"),
        EncodingType::BASE64 => &STANDARD.decode(data).expect("Wrong base64 format data"),
    };
    let ciphertext = cipher.encrypt(&nonce, data_bytes).unwrap();
    AesEncryptedData::new(hex::encode(&ciphertext), hex::encode(&nonce))
}
#[wasm_bindgen]
pub fn decrypt_aes(data: AesEncryptedData, key: String, encoding_type: EncodingType) -> String {
    // Secret key just fit into AES key
    let hex_key = hex::decode(key).expect("Key is malformed.");
    let aes_key = Key::<Aes256Gcm>::from_slice(&hex_key);
    let decipher = Aes256Gcm::new(&aes_key);
    // Nonce hex to bytes
    let nonce_bytes = hex::decode(&data.nonce()).unwrap();
    let nonce = Nonce::from_slice(&nonce_bytes);
    // Data hex to bytes
    let data_bytes = hex::decode(&data.data()).unwrap();
    let decrypted_data = decipher.decrypt(&nonce, data_bytes.as_slice()).unwrap();

    return match encoding_type {
        EncodingType::UTF8 => String::from_utf8(decrypted_data).expect("Wrong utf8 format data"),
        EncodingType::HEX => hex::encode(decrypted_data),
        EncodingType::BASE64 => STANDARD.encode(decrypted_data),
    };
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
    let alice = secp256k1_zkp::SecretKey::from_str(&privkey).unwrap();
    // Bob: param pubkey is the one who's able to decrypt the data.
    let bob = secp256k1_zkp::PublicKey::from_str(&pubkey).unwrap();
    // Only Alice and Bob can know the secret, which is used as a key to sign or encrypt data(or any other).
    let shared_secret = ecdh::SharedSecret::new(&bob, &alice);
    shared_secret.display_secret().to_string()
}
// Pederson commitment is building block for zkp.
// commitment is additively homomorphic.
// blinding_factor must be kept as secret, to prevent the confidentiality of the committed value.
// Tag is for domain(or purpose) separation.
#[wasm_bindgen]
pub fn pedersen_commit(value: u64, tag: String) -> Commitment {
    let secp = Secp256k1::new();
    let blinding_factor = Tweak::new(&mut Secp256k1Rng);
    let tag = Tag::from(sha256::Hash::hash(tag.as_bytes()).to_byte_array());
    Commitment::new(
        PedersenCommitment::new(
            &secp,
            value,
            blinding_factor,
            // Blinded Generator is used in MimbleWimble for
            // 1. Unlinkability: G' unique to a specific transaction, any commitments using it not linked to commitments from other transactions.
            // 2. Proof of Ownership: the sum of all blinding factors in a transaction must equal a public key.
            // Original pedersen: C=v⋅G+r⋅H
            // Blinded Generator pedersen: C=v⋅G'+r⋅H (where G' = G+g⋅H, g is distinct blinding factor).
            // Here, not using Blinded Generator as we only care the confidentiality of value.
            Generator::new_unblinded(&secp, tag),
        )
        .to_string(),
        blinding_factor.to_string(),
    )
}
// Perderson verify by revealing value.
#[wasm_bindgen]
pub fn pedersen_reveal(commitment: Commitment, value: u64, tag: String) -> bool {
    let secp = Secp256k1::new();
    let blinding_factor =
        Tweak::from_str(&commitment.secret_blinding_factor()).expect("Wrong blinding factor");
    let tag = Tag::from(sha256::Hash::hash(tag.as_bytes()).to_byte_array());

    verify_commitments_sum_to_equal(
        &secp,
        &vec![PedersenCommitment::new(
            &secp,
            value,
            blinding_factor,
            // Blinded Generator is used in MimbleWimble for
            // 1. Unlinkability: G' unique to a specific transaction, any commitments using it not linked to commitments from other transactions.
            // 2. Proof of Ownership: the sum of all blinding factors in a transaction must equal a public key
            // Original pedersen: C=v⋅G+r⋅H
            // Blinded Generator pedersen: C=v⋅G'+r⋅H (where G' = G+g⋅H, g is distinct blinding factor)
            // Here, not using Blinded Generator as we only care the confidentiality of value
            Generator::new_unblinded(&secp, tag),
        )],
        &vec![PedersenCommitment::from_str(&commitment.commitment()).expect("Wrong commitment")],
    )
}
