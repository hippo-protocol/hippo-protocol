use cosmwasm_std::StdError;
use thiserror::Error;

#[derive(Error, Debug)]
pub enum ContractError {
    #[error("{0}")]
    Std(#[from] StdError),

    #[error("Unauthorized")]
    Unauthorized {},

    #[error("Name has been taken")]
    NameTaken { name: String },

    #[error("Name does not exist (name {name})")]
    NameNotExists { name: String },

    #[error("Insufficient funds sent")]
    InsufficientFundsSend {},

    #[error("Name too short (minimum 3 characters)")]
    NameTooShort {},

    #[error("Name too long (maximum 64 characters)")]
    NameTooLong {},

    #[error("Invalid character in name")]
    InvalidCharacter {},
}
