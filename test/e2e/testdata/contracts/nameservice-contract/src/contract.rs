#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{
    to_json_binary, Binary, Coin, Deps, DepsMut, Env, MessageInfo, Response, StdResult,
};

use crate::error::ContractError;
use crate::msg::{ConfigResponse, ExecuteMsg, InstantiateMsg, QueryMsg, ResolveRecordResponse};
use crate::state::{Config, NameRecord, CONFIG, NAME_RESOLVER};

const MIN_NAME_LENGTH: usize = 3;
const MAX_NAME_LENGTH: usize = 64;

fn validate_name(name: &str) -> Result<(), ContractError> {
    let length = name.len();
    if length < MIN_NAME_LENGTH {
        return Err(ContractError::NameTooShort {});
    }
    if length > MAX_NAME_LENGTH {
        return Err(ContractError::NameTooLong {});
    }

    // Only allow alphanumeric and hyphens
    if !name.chars().all(|c| c.is_alphanumeric() || c == '-') {
        return Err(ContractError::InvalidCharacter {});
    }

    Ok(())
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    let config = Config {
        purchase_price: msg.purchase_price,
        transfer_price: msg.transfer_price,
    };
    CONFIG.save(deps.storage, &config)?;

    Ok(Response::new().add_attribute("method", "instantiate"))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::Register { name } => execute_register(deps, info, name),
        ExecuteMsg::Transfer { name, to } => execute_transfer(deps, info, name, to),
    }
}

pub fn execute_register(
    deps: DepsMut,
    info: MessageInfo,
    name: String,
) -> Result<Response, ContractError> {
    validate_name(&name)?;

    let config = CONFIG.load(deps.storage)?;
    
    // Check if the name is already registered
    if NAME_RESOLVER.has(deps.storage, &name) {
        return Err(ContractError::NameTaken { name });
    }

    // Verify payment
    if let Some(price) = config.purchase_price {
        let sent = info
            .funds
            .iter()
            .find(|c| c.denom == price.denom)
            .map(|c| c.amount)
            .unwrap_or_default();
        
        if sent < price.amount {
            return Err(ContractError::InsufficientFundsSend {});
        }
    }

    // Register the name
    let record = NameRecord {
        owner: info.sender.to_string(),
    };
    NAME_RESOLVER.save(deps.storage, &name, &record)?;

    Ok(Response::new()
        .add_attribute("action", "register")
        .add_attribute("name", name)
        .add_attribute("owner", info.sender))
}

pub fn execute_transfer(
    deps: DepsMut,
    info: MessageInfo,
    name: String,
    to: String,
) -> Result<Response, ContractError> {
    validate_name(&name)?;
    
    let config = CONFIG.load(deps.storage)?;
    
    // Load the name record
    let mut record = NAME_RESOLVER
        .may_load(deps.storage, &name)?
        .ok_or(ContractError::NameNotExists { name: name.clone() })?;

    // Check ownership
    if record.owner != info.sender {
        return Err(ContractError::Unauthorized {});
    }

    // Verify payment
    if let Some(price) = config.transfer_price {
        let sent = info
            .funds
            .iter()
            .find(|c| c.denom == price.denom)
            .map(|c| c.amount)
            .unwrap_or_default();
        
        if sent < price.amount {
            return Err(ContractError::InsufficientFundsSend {});
        }
    }

    // Transfer the name
    let to_addr = deps.api.addr_validate(&to)?;
    record.owner = to_addr.to_string();
    NAME_RESOLVER.save(deps.storage, &name, &record)?;

    Ok(Response::new()
        .add_attribute("action", "transfer")
        .add_attribute("name", name)
        .add_attribute("from", info.sender)
        .add_attribute("to", to))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::ResolveRecord { name } => to_json_binary(&query_resolver(deps, name)?),
        QueryMsg::Config {} => to_json_binary(&query_config(deps)?),
    }
}

fn query_resolver(deps: Deps, name: String) -> StdResult<ResolveRecordResponse> {
    let record = NAME_RESOLVER.may_load(deps.storage, &name)?;
    
    let address = record.map(|r| r.owner);
    Ok(ResolveRecordResponse { address })
}

fn query_config(deps: Deps) -> StdResult<ConfigResponse> {
    let config = CONFIG.load(deps.storage)?;
    Ok(ConfigResponse {
        purchase_price: config.purchase_price,
        transfer_price: config.transfer_price,
    })
}

#[cfg(test)]
mod tests {
    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{coins, from_json, Uint128};

    #[test]
    fn proper_initialization() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            purchase_price: Some(Coin {
                denom: "token".to_string(),
                amount: Uint128::new(100),
            }),
            transfer_price: Some(Coin {
                denom: "token".to_string(),
                amount: Uint128::new(50),
            }),
        };
        let info = mock_info("creator", &coins(2, "token"));

        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());
    }

    #[test]
    fn register_name() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            purchase_price: Some(Coin {
                denom: "token".to_string(),
                amount: Uint128::new(100),
            }),
            transfer_price: Some(Coin {
                denom: "token".to_string(),
                amount: Uint128::new(50),
            }),
        };
        let info = mock_info("creator", &coins(2, "token"));
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Register a name
        let info = mock_info("alice", &coins(100, "token"));
        let msg = ExecuteMsg::Register {
            name: "alice".to_string(),
        };
        let _res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Query the name
        let res = query(
            deps.as_ref(),
            mock_env(),
            QueryMsg::ResolveRecord {
                name: "alice".to_string(),
            },
        )
        .unwrap();
        let value: ResolveRecordResponse = from_json(&res).unwrap();
        assert_eq!(Some("alice".to_string()), value.address);
    }

    #[test]
    fn transfer_name() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            purchase_price: Some(Coin {
                denom: "token".to_string(),
                amount: Uint128::new(100),
            }),
            transfer_price: Some(Coin {
                denom: "token".to_string(),
                amount: Uint128::new(50),
            }),
        };
        let info = mock_info("creator", &coins(2, "token"));
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Register a name
        let info = mock_info("alice", &coins(100, "token"));
        let msg = ExecuteMsg::Register {
            name: "alice".to_string(),
        };
        let _res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Transfer the name
        let info = mock_info("alice", &coins(50, "token"));
        let msg = ExecuteMsg::Transfer {
            name: "alice".to_string(),
            to: "bob".to_string(),
        };
        let _res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Query the name - should now be owned by bob
        let res = query(
            deps.as_ref(),
            mock_env(),
            QueryMsg::ResolveRecord {
                name: "alice".to_string(),
            },
        )
        .unwrap();
        let value: ResolveRecordResponse = from_json(&res).unwrap();
        assert_eq!(Some("bob".to_string()), value.address);
    }

    #[test]
    fn fails_on_name_too_short() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            purchase_price: None,
            transfer_price: None,
        };
        let info = mock_info("creator", &coins(2, "token"));
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        let info = mock_info("alice", &[]);
        let msg = ExecuteMsg::Register {
            name: "ab".to_string(),
        };
        let err = execute(deps.as_mut(), mock_env(), info, msg).unwrap_err();
        match err {
            ContractError::NameTooShort {} => {}
            e => panic!("unexpected error: {:?}", e),
        }
    }
}
