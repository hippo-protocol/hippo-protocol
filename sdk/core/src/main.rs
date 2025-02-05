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

pub struct KeyPair {
    pub pubkey: String,
    pub privkey: String,
}

pub struct Did {
    pub id: String,
}

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
    fn encrypt(data: String, pubkey: String) -> String;
    fn decrypt(data: String, privkey: String) -> String;
    fn sign(data: String, privkey: String) -> String;
    fn verify(data: String, sig: String, pubkey: String) -> String;
    fn sha256(data: String) -> String;
}

fn main() {
    println!("Hello, world!");
}
