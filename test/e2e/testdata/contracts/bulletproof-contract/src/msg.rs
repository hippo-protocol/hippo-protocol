use cosmwasm_schema::{cw_serde, QueryResponses};

#[cw_serde]
pub struct InstantiateMsg {
    /// Hex-encoded serialized bulletproof range proof bytes
    pub proof_hex: String,
    /// Hex-encoded serialized committed value (CompressedRistretto, 32 bytes)
    pub commitment_hex: String,
    /// Bit size used for the range proof (8, 16, 32, or 64)
    pub num_bits: u32,
}

#[cw_serde]
pub enum ExecuteMsg {
    /// Verify the stored bulletproof range proof
    Verify {},
}

#[cw_serde]
#[derive(QueryResponses)]
pub enum QueryMsg {
    /// Returns the stored proof hash
    #[returns(GetProofHashResponse)]
    GetProofHash {},
    /// Verifies the stored bulletproof proof and returns the result
    #[returns(VerifyResponse)]
    Verify {},
}

#[cw_serde]
pub struct GetProofHashResponse {
    pub proof_hash: String,
}

#[cw_serde]
pub struct VerifyResponse {
    pub is_valid: bool,
    pub proof_hash: String,
    pub num_bits: u32,
}
