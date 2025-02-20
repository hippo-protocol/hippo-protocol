use std::str::FromStr;

use aes_gcm::{
    aead::{Aead, AeadCore, OsRng as AesRng},
    Aes256Gcm, Key, KeyInit, Nonce,
};
use secp256k1::hashes::{hex::FromHex, sha256, Hash};
use secp256k1::rand::rngs::OsRng as Secp256k1Rng;
use secp256k1::{ecdh, Message, PublicKey, Secp256k1, SecretKey};
use secp256k1::{ecdsa::Signature, hashes::hex::DisplayHex};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Tx {
    /// coin name(e.g. HP).
    pub coin: String,
    /// mostly signer address
    pub from: String,
    /// address to send
    pub to: String,
    /// amount to send
    pub amount: String,
    /// fee to broadcast tx
    pub fee: String,
    /// data to be included in tx
    pub data: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct KeyPair {
    pub pubkey: String,
    pub privkey: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct Did {
    pub id: String,
}

#[derive(Serialize, Deserialize, Debug, Clone)]
pub struct EncryptedData {
    pub pubkey_from: String,
    pub pubkey_to: String,
    pub data: String,
    pub nonce: String,
}

pub struct Hippo {}

pub trait Sdk {
    /*
    /// chain
    fn write(data: String, privkey: String) -> Tx;
    fn read(tx_hash: String) -> Tx;
    fn send(from: String, to: String, privkey: String, data: Option<String>) -> Tx;
     */
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

impl Sdk for Hippo {
    fn encrypt(data: String, pubkey: String) -> EncryptedData {
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

        EncryptedData {
            pubkey_from: public_key.to_string(),
            pubkey_to: pubkey,
            data: ciphertext.to_lower_hex_string(),
            nonce: nonce.to_lower_hex_string(),
        }
    }

    fn decrypt(data: EncryptedData, privkey: String) -> String {
        // Alice: one-off pubkeyr to caculate shared secret.
        let encrypt_from_pubkey = secp256k1::PublicKey::from_str(&data.pubkey_from).unwrap();
        // Bob: the one who's able to decrypt the data with privkey.
        let privkey_to_decrypt = secp256k1::SecretKey::from_str(&privkey).unwrap();
        // Only Alice and Bob can know the secret, which is used as a key to encrypt data.
        let shared_secret =
            ecdh::SharedSecret::new(&encrypt_from_pubkey, &privkey_to_decrypt).secret_bytes();

        // Secret key just fit into AES key
        let aes_key = Key::<Aes256Gcm>::from_slice(&shared_secret);
        let cipher = Aes256Gcm::new(&aes_key);
        // Nonce hex to bytes
        let nonce_bytes = Vec::from_hex(&data.nonce).unwrap();
        let nonce = Nonce::from_slice(&nonce_bytes);
        // Data hex to bytes
        let data_bytes = Vec::from_hex(&data.data).unwrap();
        let decrypted_data = cipher.decrypt(&nonce, data_bytes.as_slice()).unwrap();

        String::from_utf8(decrypted_data).unwrap()
    }

    fn create_keypair() -> KeyPair {
        let secp = Secp256k1::new();
        let (secret_key, public_key) = secp.generate_keypair(&mut Secp256k1Rng);
        KeyPair {
            pubkey: public_key.to_string(),
            privkey: secret_key.display_secret().to_string(),
        }
    }

    fn key_to_did(pubkey: String) -> Did {
        Did {
            id: "did:hp:".to_owned() + &pubkey,
        }
    }

    fn did_to_key(did: Did) -> String {
        did.id.strip_prefix("did:hp:").unwrap().to_string()
    }

    fn sign(data: String, privkey: String) -> String {
        let secp = Secp256k1::new();
        let digest = sha256::Hash::hash(data.as_bytes());
        let message = Message::from_digest(digest.to_byte_array());
        let secret_key = SecretKey::from_str(&privkey).unwrap();

        secp.sign_ecdsa(&message, &secret_key).to_string()
    }

    fn verify(data: String, sig: String, pubkey: String) -> bool {
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

    fn sha256(data: String) -> String {
        sha256::Hash::hash(data.as_bytes()).to_string()
    }

    fn ecdh(privkey: String, pubkey: String) -> String {
        // Alice: param privkey is the one who's able to decrypt the data.
        let alice = secp256k1::SecretKey::from_str(&privkey).unwrap();
        // Bob: param pubkey is the one who's able to decrypt the data.
        let bob = secp256k1::PublicKey::from_str(&pubkey).unwrap();
        // Only Alice and Bob can know the secret, which is used as a key to sign or encrypt data(or any other).
        let shared_secret = ecdh::SharedSecret::new(&bob, &alice);
        shared_secret.display_secret().to_string()
    }
}

#[cfg(test)]
mod tests {
    use super::*; // Bring the outer functions into the test module's scope

    #[test]
    fn test_enc_dec() {
        // given
        let data = String::from("datag허ㅜㅏ니ㅜ2#@_!##ㅏ!~2ㅡ₩ㅡ1    ㅁAl;;A;;:{}()[]");
        let alice = Hippo::create_keypair();
        // when
        let enc_data = Hippo::encrypt(data.clone(), alice.pubkey);
        let dec_data = Hippo::decrypt(enc_data, alice.privkey);
        // then
        assert_eq!(data, dec_data);
    }

    #[test]
    fn test_did() {
        // given
        let key_pair = Hippo::create_keypair();
        // when
        let pubkey_to_did = Hippo::key_to_did(key_pair.pubkey.clone());
        let did_to_pubkey = Hippo::did_to_key(pubkey_to_did.clone());
        // then
        assert_eq!(key_pair.pubkey, did_to_pubkey);
        assert!(pubkey_to_did.id.starts_with("did:hp"));
    }

    #[test]
    fn test_ecdsa() {
        // given
        let key_pair = Hippo::create_keypair();
        let data = String::from("data");
        // when
        let sig = Hippo::sign(data.clone(), key_pair.privkey);
        let is_verified = Hippo::verify(data, sig, key_pair.pubkey);
        // then
        assert!(is_verified);
    }
}

fn main() {}
