#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{to_json_binary, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult};

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, GetProofHashResponse, InstantiateMsg, QueryMsg, VerifyResponse};
use crate::state::{ProofState, PROOF_STATE};

use bulletproofs::{BulletproofGens, PedersenGens, RangeProof};
use curve25519_dalek_ng::ristretto::CompressedRistretto;
use merlin::Transcript;
use rand_chacha::ChaCha20Rng;
use rand_core::SeedableRng;
use sha2::{Digest, Sha256};

/// Compute SHA-256 hash of bytes and return hex string
fn compute_hash(data: &[u8]) -> String {
    let mut hasher = Sha256::new();
    hasher.update(data);
    let result = hasher.finalize();
    hex::encode(result)
}

/// Verify a bulletproof range proof
fn verify_bulletproof(proof_bytes: &[u8], commitment_bytes: &[u8], num_bits: u32) -> Result<bool, ContractError> {
    let pc_gens = PedersenGens::default();
    let bp_gens = BulletproofGens::new(64, 1);

    let proof = RangeProof::from_bytes(proof_bytes)
        .map_err(|e| ContractError::InvalidProof {
            reason: format!("failed to deserialize proof: {:?}", e),
        })?;

    if commitment_bytes.len() != 32 {
        return Err(ContractError::InvalidProof {
            reason: format!("commitment must be 32 bytes, got {}", commitment_bytes.len()),
        });
    }
    let commitment = CompressedRistretto::from_slice(commitment_bytes);

    let mut verifier_transcript = Transcript::new(b"doctest example");

    // Use a deterministic RNG for verification in the wasm environment
    let mut rng = ChaCha20Rng::from_seed([0u8; 32]);

    Ok(proof
        .verify_single_with_rng(
            &bp_gens,
            &pc_gens,
            &mut verifier_transcript,
            &commitment,
            num_bits as usize,
            &mut rng,
        )
        .is_ok())
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    // Validate hex encoding
    let proof_bytes = hex::decode(&msg.proof_hex).map_err(|e| ContractError::InvalidProof {
        reason: format!("invalid proof hex: {}", e),
    })?;

    let commitment_bytes = hex::decode(&msg.commitment_hex).map_err(|e| ContractError::InvalidProof {
        reason: format!("invalid commitment hex: {}", e),
    })?;

    if commitment_bytes.len() != 32 {
        return Err(ContractError::InvalidProof {
            reason: format!("commitment must be 32 bytes, got {}", commitment_bytes.len()),
        });
    }

    // Validate num_bits
    if msg.num_bits != 8 && msg.num_bits != 16 && msg.num_bits != 32 && msg.num_bits != 64 {
        return Err(ContractError::InvalidProof {
            reason: format!("num_bits must be 8, 16, 32, or 64, got {}", msg.num_bits),
        });
    }

    // Compute SHA-256 hash of proof bytes
    let proof_hash = compute_hash(&proof_bytes);

    let state = ProofState {
        proof_hex: msg.proof_hex,
        commitment_hex: msg.commitment_hex,
        num_bits: msg.num_bits,
        proof_hash: proof_hash.clone(),
        owner: info.sender.to_string(),
    };
    PROOF_STATE.save(deps.storage, &state)?;

    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("owner", info.sender)
        .add_attribute("proof_hash", proof_hash)
        .add_attribute("num_bits", msg.num_bits.to_string()))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    _env: Env,
    _info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::Verify {} => execute_verify(deps),
    }
}

pub fn execute_verify(deps: DepsMut) -> Result<Response, ContractError> {
    let state = PROOF_STATE.load(deps.storage)?;

    let proof_bytes = hex::decode(&state.proof_hex).map_err(|e| ContractError::InvalidProof {
        reason: format!("invalid stored proof hex: {}", e),
    })?;
    let commitment_bytes = hex::decode(&state.commitment_hex).map_err(|e| ContractError::InvalidProof {
        reason: format!("invalid stored commitment hex: {}", e),
    })?;

    let is_valid = verify_bulletproof(&proof_bytes, &commitment_bytes, state.num_bits)?;

    Ok(Response::new()
        .add_attribute("action", "verify")
        .add_attribute("is_valid", is_valid.to_string())
        .add_attribute("proof_hash", &state.proof_hash))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetProofHash {} => to_json_binary(&query_proof_hash(deps)?),
        QueryMsg::Verify {} => to_json_binary(&query_verify(deps)?),
    }
}

fn query_proof_hash(deps: Deps) -> StdResult<GetProofHashResponse> {
    let state = PROOF_STATE.load(deps.storage)?;
    Ok(GetProofHashResponse {
        proof_hash: state.proof_hash,
    })
}

fn query_verify(deps: Deps) -> StdResult<VerifyResponse> {
    let state = PROOF_STATE.load(deps.storage)?;

    let proof_bytes = hex::decode(&state.proof_hex)
        .unwrap_or_default();
    let commitment_bytes = hex::decode(&state.commitment_hex)
        .unwrap_or_default();

    let is_valid = verify_bulletproof(&proof_bytes, &commitment_bytes, state.num_bits)
        .unwrap_or(false);

    Ok(VerifyResponse {
        is_valid,
        proof_hash: state.proof_hash,
        num_bits: state.num_bits,
    })
}

#[cfg(test)]
mod tests {
    use super::*;
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{coins, from_json};

    // Pre-generated bulletproof test data
    // Secret value: 1037578891, num_bits: 32, transcript label: "doctest example"
    const TEST_PROOF_HEX: &str = "8026afd76427529f11bcc07e29a182e3122bab7595b61237dda31548ba96cc3e4a84148c615bb889cd99bab5519e2e7d815a2469b76b5e6bf56c1051264f9b5b0a75a84e179a21b7701de8b744612ecd96b5e73f2ad4ffda4dde5a0bf0fa5b4490bd0c7ec41b331068f3db152278f5147c876201e741a6817616ece7c58a6507a7736c1fc341bb3ab65cc6e7196855a42eed503f04b56b190fced87eab134400c9fdcb1eb43fc7fed2882b2f56b9eea62ce8a024bca4f23aa4d70afb323d4c0ad3a38d409012207bb35e174a112794008d2c3f8a0d7f4282ab718493096da30d5e432f7917f017e4ee80191990aed9a51d404700c1e441ef3c46e83129aa2f5b4a1047757dc4ce4c11d1ea429c7a95dd95bc13f7c9fd5b4c64c5aa97040948142a72e57ef4658bf2894029fc69dcd893fe5bf72d90aced60e2b4608b0bfa6a06f26414c843a86df58d95f92c1904565898262d1170ad70252445bd883ec208415ef350cb0515a602d37cbb668d78e6f6211fa4caf338513c5e551f3a36b33214c89e9681301b830da28be02204d062ca19b2edacd56fa5ce4c7e1d0a9f1fd85f7049fe27dffc601b41f35dce8b0f61b3c92a8f51ab40299e6bf452c81d95ee1880dede9a6da64b3237451715c8da5970296d3b34b0c9b585b355f31e2b71c46cd49fe004d2ec5371b3029ee2d6d0881d90d73ac81b1d16a82a74f46e36b14e33a6abaf35b81fdccbe00031d8c5918974f53d35973cd7077b839c2dbfcead70236581065006dbb5f1e1541fa6226e172e0e9471a7a0a1ed5aa627d26e9aac140f0b2ddee38a4502fe9f6327e81fdb849cd7c7698e9add48aecab22f512b56fd0b";
    const TEST_COMMITMENT_HEX: &str = "5e50cca6bdd5d8c04e1a2848d74d885647d93b883cb4f182fbb5e3bdbf00506c";

    #[test]
    fn proper_initialization_with_invalid_hex() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            proof_hex: "invalid_hex".to_string(),
            commitment_hex: "0000000000000000000000000000000000000000000000000000000000000000".to_string(),
            num_bits: 32,
        };
        let info = mock_info("creator", &coins(1000, "earth"));

        let res = instantiate(deps.as_mut(), mock_env(), info, msg);
        assert!(res.is_err());
    }

    #[test]
    fn proper_initialization_with_invalid_bits() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            proof_hex: "00".to_string(),
            commitment_hex: "0000000000000000000000000000000000000000000000000000000000000000".to_string(),
            num_bits: 15,
        };
        let info = mock_info("creator", &coins(1000, "earth"));

        let res = instantiate(deps.as_mut(), mock_env(), info, msg);
        assert!(res.is_err());
    }

    #[test]
    fn query_proof_hash_works() {
        let mut deps = mock_dependencies();

        // Use a valid hex string for proof (doesn't need to be a valid proof for hash test)
        let proof_hex = "abcdef0123456789".to_string();
        let commitment_hex = "0000000000000000000000000000000000000000000000000000000000000000".to_string();

        let msg = InstantiateMsg {
            proof_hex: proof_hex.clone(),
            commitment_hex,
            num_bits: 32,
        };
        let info = mock_info("creator", &coins(1000, "earth"));
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Query the proof hash
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetProofHash {}).unwrap();
        let value: GetProofHashResponse = from_json(&res).unwrap();

        // Compute expected hash
        let proof_bytes = hex::decode(&proof_hex).unwrap();
        let expected_hash = compute_hash(&proof_bytes);
        assert_eq!(value.proof_hash, expected_hash);
    }

    #[test]
    fn verify_valid_proof() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            proof_hex: TEST_PROOF_HEX.to_string(),
            commitment_hex: TEST_COMMITMENT_HEX.to_string(),
            num_bits: 32,
        };
        let info = mock_info("creator", &coins(1000, "earth"));
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Query verify
        let res = query(deps.as_ref(), mock_env(), QueryMsg::Verify {}).unwrap();
        let value: VerifyResponse = from_json(&res).unwrap();
        assert!(value.is_valid, "bulletproof verification should succeed");
        assert_eq!(value.num_bits, 32);
    }

    #[test]
    fn verify_invalid_proof_returns_false() {
        let mut deps = mock_dependencies();

        // Use valid hex but invalid proof bytes (all zeros)
        let msg = InstantiateMsg {
            proof_hex: hex::encode([0u8; 608]),
            commitment_hex: TEST_COMMITMENT_HEX.to_string(),
            num_bits: 32,
        };
        let info = mock_info("creator", &coins(1000, "earth"));
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Query verify - should fail because proof is invalid
        let res = query(deps.as_ref(), mock_env(), QueryMsg::Verify {}).unwrap();
        let value: VerifyResponse = from_json(&res).unwrap();
        assert!(!value.is_valid, "invalid proof should fail verification");
    }

    #[test]
    fn execute_verify_valid_proof() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            proof_hex: TEST_PROOF_HEX.to_string(),
            commitment_hex: TEST_COMMITMENT_HEX.to_string(),
            num_bits: 32,
        };
        let info = mock_info("creator", &coins(1000, "earth"));
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Execute verify
        let info = mock_info("anyone", &coins(2, "token"));
        let msg = ExecuteMsg::Verify {};
        let res = execute(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Check attributes
        let is_valid_attr = res.attributes.iter().find(|a| a.key == "is_valid").unwrap();
        assert_eq!(is_valid_attr.value, "true");

        let action_attr = res.attributes.iter().find(|a| a.key == "action").unwrap();
        assert_eq!(action_attr.value, "verify");
    }

    #[test]
    fn generate_and_verify_fresh_proof() {
        use curve25519_dalek_ng::scalar::Scalar;

        let mut deps = mock_dependencies();

        // Generate a fresh proof
        let pc_gens = PedersenGens::default();
        let bp_gens = BulletproofGens::new(64, 1);
        let secret_value = 42u64;
        let mut rng = rand::thread_rng();
        let blinding = Scalar::random(&mut rng);
        let mut prover_transcript = Transcript::new(b"doctest example");

        let (proof, committed_value) = RangeProof::prove_single_with_rng(
            &bp_gens,
            &pc_gens,
            &mut prover_transcript,
            secret_value,
            &blinding,
            32,
            &mut rng,
        ).unwrap();

        let proof_bytes: Vec<u8> = proof.to_bytes();
        let commitment_bytes: [u8; 32] = committed_value.to_bytes();
        let proof_hex = hex::encode(&proof_bytes);
        let commitment_hex = hex::encode(&commitment_bytes);

        // Instantiate with fresh proof
        let msg = InstantiateMsg {
            proof_hex,
            commitment_hex,
            num_bits: 32,
        };
        let info = mock_info("creator", &coins(1000, "earth"));
        let _res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();

        // Verify via query
        let res = query(deps.as_ref(), mock_env(), QueryMsg::Verify {}).unwrap();
        let value: VerifyResponse = from_json(&res).unwrap();
        assert!(value.is_valid, "freshly generated proof should verify");
    }
}
