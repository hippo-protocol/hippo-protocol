#[cfg(test)]
mod tests {
    use crate::{
        create_keypair, decrypt, did_to_key, encrypt, key_to_did, pedersen_commit, pedersen_reveal,
        sign,
        types::{Commitment, EncodingType},
        verify,
    };
    use base64::{engine::general_purpose::STANDARD, Engine as _};

    #[test]
    fn test_enc_dec() {
        // given
        let utf8_data = String::from("datag허ㅜㅏ니ㅜ2#@_!##ㅏ!~2ㅡ₩ㅡ1    ㅁAl;;A;;:{}()[]");
        let hex_data = hex::encode(utf8_data.clone());
        let base64_data = STANDARD.encode(utf8_data.clone());
        let alice = create_keypair();
        // when
        let utf8_enc_data = encrypt(utf8_data.clone(), alice.pubkey(), EncodingType::UTF8).unwrap();
        let utf8_dec_data = decrypt(utf8_enc_data, alice.privkey(), EncodingType::UTF8).unwrap();
        let hex_enc_data = encrypt(hex_data.clone(), alice.pubkey(), EncodingType::HEX).unwrap();
        let hex_dec_data = decrypt(hex_enc_data, alice.privkey(), EncodingType::HEX).unwrap();
        let base64_enc_data =
            encrypt(base64_data.clone(), alice.pubkey(), EncodingType::BASE64).unwrap();
        let base64_dec_data =
            decrypt(base64_enc_data, alice.privkey(), EncodingType::BASE64).unwrap();
        // then
        assert_eq!(utf8_data, utf8_dec_data);
        assert_eq!(hex_data, hex_dec_data);
        assert_eq!(base64_data, base64_dec_data);
    }

    #[test]
    fn test_did() {
        // given
        let key_pair = create_keypair();
        // when
        let pubkey_to_did = key_to_did(key_pair.pubkey());
        let did_to_pubkey = did_to_key(pubkey_to_did.clone()).unwrap();
        // then
        assert_eq!(key_pair.pubkey(), did_to_pubkey);
        assert!(pubkey_to_did.id().starts_with("did:hp"));
    }

    #[test]
    fn test_ecdsa() {
        // given
        let key_pair = create_keypair();
        let data = String::from("data");
        // when
        let sig = sign(data.clone(), key_pair.privkey()).unwrap();
        let is_verified = verify(data, sig, key_pair.pubkey()).unwrap();
        // then
        assert!(is_verified);
    }

    #[test]
    fn test_pedersen_commit() {
        // given
        let tag = String::from("hippo");
        let value = 100_u64;
        let wrong_tag = String::from("wrong hippo");
        let wrong_value = 0_u64;
        // when
        let commitment = pedersen_commit(value, tag.clone());
        let commitment_same_value = pedersen_commit(value, tag.clone());
        let is_verified = pedersen_reveal(commitment.clone(), value, tag.clone()).unwrap();
        // These should return Ok(false) or Err?
        // Current implementation: verify_commitments_sum_to_equal returns bool.
        // And we wrapped it in Ok(). So it returns Ok(bool).
        // If the *structure* is wrong (blinding factor format etc), it returns Err.
        // If structure is right but math is wrong, it returns Ok(false).
        // We need to check what the test expects.
        // Original code: pedersen_reveal returned bool.
        // new code: pedersen_reveal returns Result<bool, JsError>.
        // So we expect Ok(false) for logic mismatch.
        let wrong_value_and_tag =
            pedersen_reveal(commitment.clone(), wrong_value, wrong_tag).unwrap();
        let wrong_blinding_factor_with_same_value = pedersen_reveal(
            Commitment::new(
                commitment_same_value.commitment(),
                // Value is same but blinding factor is different.
                commitment.secret_blinding_factor(),
            ),
            value,
            tag.clone(),
        )
        .unwrap();
        let wrong_commitment_with_correct_blinding_factor = pedersen_reveal(
            Commitment::new(
                // Even if the value is same, blinding factor makes the commitment different.
                commitment.commitment(),
                commitment_same_value.secret_blinding_factor(), // Blinding factor is correct.
            ),
            value,
            tag,
        )
        .unwrap();
        // then
        assert!(is_verified);
        assert!(!wrong_value_and_tag);
        assert!(!wrong_blinding_factor_with_same_value);
        assert!(!wrong_commitment_with_correct_blinding_factor)
    }
}
