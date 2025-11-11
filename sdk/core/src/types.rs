use serde::{Deserialize, Serialize};
use wasm_bindgen::prelude::*;

#[derive(Serialize, Deserialize, Debug, Clone)]
#[wasm_bindgen]
pub struct KeyPair {
    pubkey: String,
    privkey: String,
}

#[wasm_bindgen]
impl KeyPair {
    #[wasm_bindgen(constructor)]
    pub fn new(pubkey: String, privkey: String) -> Self {
        KeyPair { pubkey, privkey }
    }

    #[wasm_bindgen]
    pub fn to_object(&self) -> JsValue {
        serde_wasm_bindgen::to_value(&self).unwrap()
    }

    #[wasm_bindgen]
    pub fn from_object(object: JsValue) -> KeyPair {
        serde_wasm_bindgen::from_value(object)
            .map_err(|e| JsValue::from_str(&format!("Failed to deserialize: {}", e)))
            .unwrap()
    }

    #[wasm_bindgen(getter)]
    pub fn pubkey(&self) -> String {
        self.pubkey.clone()
    }

    #[wasm_bindgen(getter)]
    pub fn privkey(&self) -> String {
        self.privkey.clone()
    }

    #[wasm_bindgen(setter)]
    pub fn set_pubkey(&mut self, pubkey: String) {
        self.pubkey = pubkey;
    }

    #[wasm_bindgen(setter)]
    pub fn set_privkey(&mut self, privkey: String) {
        self.privkey = privkey;
    }
}

#[derive(Serialize, Deserialize, Debug, Clone)]
#[wasm_bindgen]
pub struct Did {
    id: String,
}

#[wasm_bindgen]
impl Did {
    #[wasm_bindgen(constructor)]
    pub fn new(id: String) -> Self {
        Did { id }
    }
    #[wasm_bindgen]
    pub fn to_object(&self) -> JsValue {
        serde_wasm_bindgen::to_value(&self).unwrap()
    }
    #[wasm_bindgen]
    pub fn from_object(object: JsValue) -> Did {
        serde_wasm_bindgen::from_value(object)
            .map_err(|e| JsValue::from_str(&format!("Failed to deserialize: {}", e)))
            .unwrap()
    }
    #[wasm_bindgen(getter)]
    pub fn id(&self) -> String {
        self.id.clone()
    }
    #[wasm_bindgen(setter)]
    pub fn set_id(&mut self, id: String) {
        self.id = id;
    }
}

#[derive(Serialize, Deserialize, Debug, Clone)]
#[wasm_bindgen]
pub struct EncryptedData {
    pubkey_from: String,
    pubkey_to: String,
    data: String,
    nonce: String,
}

#[wasm_bindgen]
impl EncryptedData {
    #[wasm_bindgen(constructor)]
    pub fn new(pubkey_from: String, pubkey_to: String, data: String, nonce: String) -> Self {
        EncryptedData {
            pubkey_from,
            pubkey_to,
            data,
            nonce,
        }
    }
    #[wasm_bindgen]
    pub fn to_object(&self) -> JsValue {
        serde_wasm_bindgen::to_value(&self).unwrap()
    }
    #[wasm_bindgen]
    pub fn from_object(object: JsValue) -> EncryptedData {
        serde_wasm_bindgen::from_value(object)
            .map_err(|e| JsValue::from_str(&format!("Failed to deserialize: {}", e)))
            .unwrap()
    }
    #[wasm_bindgen(getter)]
    pub fn pubkey_from(&self) -> String {
        self.pubkey_from.clone()
    }
    #[wasm_bindgen(getter)]
    pub fn pubkey_to(&self) -> String {
        self.pubkey_to.clone()
    }
    #[wasm_bindgen(getter)]
    pub fn data(&self) -> String {
        self.data.clone()
    }
    #[wasm_bindgen(getter)]
    pub fn nonce(&self) -> String {
        self.nonce.clone()
    }
}

#[derive(Serialize, Deserialize, Debug, Clone)]
#[wasm_bindgen]
pub struct AesEncryptedData {
    data: String,
    nonce: String,
}

#[wasm_bindgen]
impl AesEncryptedData {
    #[wasm_bindgen(constructor)]
    pub fn new(data: String, nonce: String) -> Self {
        AesEncryptedData { data, nonce }
    }
    #[wasm_bindgen]
    pub fn to_object(&self) -> JsValue {
        serde_wasm_bindgen::to_value(&self).unwrap()
    }
    #[wasm_bindgen]
    pub fn from_object(object: JsValue) -> AesEncryptedData {
        serde_wasm_bindgen::from_value(object)
            .map_err(|e| JsValue::from_str(&format!("Failed to deserialize: {}", e)))
            .unwrap()
    }
    #[wasm_bindgen(getter)]
    pub fn data(&self) -> String {
        self.data.clone()
    }
    #[wasm_bindgen(getter)]
    pub fn nonce(&self) -> String {
        self.nonce.clone()
    }
}

#[derive(Serialize, Deserialize, Debug, Clone)]
#[wasm_bindgen]
pub struct Tx {
    /// coin name(e.g. HP).
    coin: String,
    /// mostly signer address
    from: String,
    /// address to send
    to: String,
    /// amount to send
    amount: String,
    /// fee to broadcast tx
    fee: String,
    /// data to be included in tx
    data: String,
}

#[wasm_bindgen]
impl Tx {
    #[wasm_bindgen(constructor)]
    pub fn new(
        coin: String,
        from: String,
        to: String,
        amount: String,
        fee: String,
        data: String,
    ) -> Self {
        Tx {
            coin,
            from,
            to,
            amount,
            fee,
            data,
        }
    }
    #[wasm_bindgen]
    pub fn to_object(&self) -> JsValue {
        serde_wasm_bindgen::to_value(&self).unwrap()
    }
    #[wasm_bindgen]
    pub fn from_object(object: JsValue) -> Tx {
        serde_wasm_bindgen::from_value(object)
            .map_err(|e| JsValue::from_str(&format!("Failed to deserialize: {}", e)))
            .unwrap()
    }
    #[wasm_bindgen(getter)]
    pub fn coin(&self) -> String {
        self.coin.clone()
    }
    #[wasm_bindgen(getter)]
    pub fn from(&self) -> String {
        self.from.clone()
    }
    #[wasm_bindgen(getter)]
    pub fn to(&self) -> String {
        self.to.clone()
    }
    #[wasm_bindgen(getter)]
    pub fn amount(&self) -> String {
        self.amount.clone()
    }
    #[wasm_bindgen(getter)]
    pub fn fee(&self) -> String {
        self.fee.clone()
    }
    #[wasm_bindgen(getter)]
    pub fn data(&self) -> String {
        self.data.clone()
    }
    #[wasm_bindgen(setter)]
    pub fn set_coin(&mut self, coin: String) {
        self.coin = coin;
    }
    #[wasm_bindgen(setter)]
    pub fn set_from(&mut self, from: String) {
        self.from = from;
    }
    #[wasm_bindgen(setter)]
    pub fn set_to(&mut self, to: String) {
        self.to = to;
    }
    #[wasm_bindgen(setter)]
    pub fn set_amount(&mut self, amount: String) {
        self.amount = amount;
    }
    #[wasm_bindgen(setter)]
    pub fn set_fee(&mut self, fee: String) {
        self.fee = fee;
    }
    #[wasm_bindgen(setter)]
    pub fn set_data(&mut self, data: String) {
        self.data = data;
    }
}

#[derive(Serialize, Deserialize, Debug, Clone)]
#[wasm_bindgen]
pub struct Commitment {
    commitment: String,
    secret_blinding_factor: String,
}

#[wasm_bindgen]
impl Commitment {
    #[wasm_bindgen(constructor)]
    pub fn new(commitment: String, secret_blinding_factor: String) -> Self {
        Commitment {
            commitment,
            secret_blinding_factor,
        }
    }
    #[wasm_bindgen]
    pub fn to_object(&self) -> JsValue {
        serde_wasm_bindgen::to_value(&self).unwrap()
    }
    #[wasm_bindgen]
    pub fn from_object(object: JsValue) -> Commitment {
        serde_wasm_bindgen::from_value(object)
            .map_err(|e| JsValue::from_str(&format!("Failed to deserialize: {}", e)))
            .unwrap()
    }
    #[wasm_bindgen(getter)]
    pub fn commitment(&self) -> String {
        self.commitment.clone()
    }
    #[wasm_bindgen(setter)]
    pub fn set_commitment(&mut self, commitment: String) {
        self.commitment = commitment;
    }
    #[wasm_bindgen(getter)]
    pub fn secret_blinding_factor(&self) -> String {
        self.secret_blinding_factor.clone()
    }
    #[wasm_bindgen(setter)]
    pub fn set_secret_blinding_factor(&mut self, secret_blinding_factor: String) {
        self.secret_blinding_factor = secret_blinding_factor;
    }
}
