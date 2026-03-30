use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use cw_storage_plus::Item;

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, JsonSchema)]
pub struct ProofState {
    /// Hex-encoded serialized bulletproof range proof
    pub proof_hex: String,
    /// Hex-encoded serialized committed value (CompressedRistretto point, 32 bytes)
    pub commitment_hex: String,
    /// Bit size used for the range proof (e.g. 32 or 64)
    pub num_bits: u32,
    /// SHA-256 hash of the proof bytes, hex-encoded
    pub proof_hash: String,
    /// Address of the account that stored the proof
    pub owner: String,
}

pub const PROOF_STATE: Item<ProofState> = Item::new("proof_state");
