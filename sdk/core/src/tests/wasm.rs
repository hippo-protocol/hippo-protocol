use wasm_bindgen_test::*;

use crate::types::EncodingType;

#[wasm_bindgen_test]
fn encrypt_decrypt() {
    use crate::{decrypt, encrypt};
    let message = String::from("Hello, world!");
    let alice = crate::create_keypair();
    let encrypted_message = encrypt(message.clone(), alice.pubkey(), EncodingType::UTF8).unwrap();
    assert_ne!(message, encrypted_message.data());
    let decrypted_message =
        decrypt(encrypted_message, alice.privkey(), EncodingType::UTF8).unwrap();
    assert_eq!(message, decrypted_message);
}

#[wasm_bindgen_test]
fn did_conversion() {
    use crate::{create_keypair, did_to_key, key_to_did};
    let alice = create_keypair();
    let did = key_to_did(alice.pubkey());
    assert_eq!(alice.pubkey(), did_to_key(did).unwrap());
}

#[wasm_bindgen_test]
fn sign_verify() {
    use crate::{sign, verify};
    let message = String::from("Hello, world!");
    let alice = crate::create_keypair();
    let signature = sign(message.clone(), alice.privkey()).unwrap();
    let verified = verify(message, signature, alice.pubkey()).unwrap();
    assert_eq!(true, verified);
}

#[wasm_bindgen_test]
fn sha256_hash() {
    use crate::sha256;
    let message = String::from("Hello, world!");
    let hash = sha256(message.clone());
    assert_eq!(
        hash,
        "315f5bdb76d078c43b8ac0064e4a0164612b1fce77c869345bfc94c75894edd3"
    );
}

#[wasm_bindgen_test]
fn ecdh_shared_secret() {
    use crate::{create_keypair, ecdh};
    let alice = create_keypair();
    let bob = create_keypair();
    let shared_secret_alice = ecdh(alice.privkey(), bob.pubkey()).unwrap();
    let shared_secret_bob = ecdh(bob.privkey(), alice.pubkey()).unwrap();
    assert_eq!(shared_secret_alice, shared_secret_bob);
}
