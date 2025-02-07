use std::str::FromStr;

use secp256k1::ecdsa::Signature;
use secp256k1::hashes::{sha256, Hash};
use secp256k1::rand::rngs::OsRng;
use secp256k1::{Message, PublicKey, Secp256k1, SecretKey};
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
    //fn encrypt(data: String, pubkey: String) -> String;
    //fn decrypt(data: String, privkey: String) -> String;
    fn sign(data: String, privkey: String) -> String;
    fn verify(data: String, sig: String, pubkey: String) -> bool;
    fn sha256(data: String) -> String;
}

impl Sdk for Hippo {
    fn create_keypair() -> KeyPair {
        let secp = Secp256k1::new();
        let (secret_key, public_key) = secp.generate_keypair(&mut OsRng);
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
}

fn main() {
    let key_pair = Hippo::create_keypair();
    let pubkey_to_did = Hippo::key_to_did(key_pair.pubkey.clone());
    let did_to_pubkey = Hippo::did_to_key(pubkey_to_did.clone());

    let data = String::from("data");
    let sig = Hippo::sign(data.clone(), key_pair.privkey.clone());
    let is_verified = Hippo::verify(data, sig.clone(), key_pair.pubkey.clone());

    println!("key_pair: {:?}", key_pair);
    println!("did: {:?}", pubkey_to_did);
    println!("pubkey: {:?}", did_to_pubkey);
    println!("sig: {:?}", sig);
    println!("is_verified: {:?}", is_verified)
}
