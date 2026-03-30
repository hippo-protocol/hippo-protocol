use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized")]
    Unauthorized {},

    #[error("Invalid proof: {reason}")]
    InvalidProof { reason: String },

    #[error("Verification failed")]
    VerificationFailed {},
}
