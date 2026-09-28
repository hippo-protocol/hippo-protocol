#[cfg(test)]
mod tests {
    use crate::{
        create_keypair, decrypt, decrypt_bytes, did_to_key, encrypt, encrypt_bytes,
        generate_bulletproof, key_to_did, pedersen_commit, pedersen_reveal, sign,
        types::{Bulletproof, Commitment, EncodingType},
        verify, verify_bulletproof,
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

    #[test]
    fn test_enc_dec_bytes() {
        // given
        let data = vec![0, 1, 2, 3, 4, 5, 6, 7, 8];
        let alice = create_keypair();
        // when
        let enc_data = encrypt_bytes(data.clone(), alice.pubkey()).unwrap();
        let dec_data = decrypt_bytes(enc_data, alice.privkey()).unwrap();
        // then
        assert_eq!(data, dec_data);
    }

    #[test]
    fn test_bulletproof() {
        // given
        let bits = 32;
        let proof = generate_bulletproof(1234, bits).unwrap();
        let other_proof = generate_bulletproof(1234, bits).unwrap();
        let mut tampered_proof = hex::decode(proof.proof()).unwrap();
        tampered_proof[0] ^= 1;
        // when
        let is_verified = verify_bulletproof(&proof, bits).unwrap();
        let wrong_bits = verify_bulletproof(&proof, 64).unwrap();
        let unsupported_bits = verify_bulletproof(&proof, 10).unwrap();
        let wrong_commitment = verify_bulletproof(
            // Same value, but commitment is from another proof.
            &Bulletproof::new(proof.proof(), other_proof.commitment()),
            bits,
        )
        .unwrap();
        let wrong_proof = verify_bulletproof(
            &Bulletproof::new(hex::encode(tampered_proof), proof.commitment()),
            bits,
        )
        .unwrap();
        // then
        assert!(is_verified);
        assert!(!wrong_bits);
        assert!(!unsupported_bits);
        assert!(!wrong_commitment);
        assert!(!wrong_proof);
    }

    #[test]
    fn test_bulletproof_range_bounds() {
        for bits in [8, 16, 32, 64] {
            // both ends of the range [0, 2^bits)
            for value in [0, u64::MAX >> (64 - bits)] {
                let proof = generate_bulletproof(value, bits).unwrap();
                assert!(verify_bulletproof(&proof, bits).unwrap());
            }
        }
    }

    #[test]
    fn test_bulletproof_onchain_proof() {
        // given: pre-generated proof of the on-chain bulletproof contract e2e test(value 1037578891, 32 bits).
        let proof = Bulletproof::new(
            String::from("8026afd76427529f11bcc07e29a182e3122bab7595b61237dda31548ba96cc3e4a84148c615bb889cd99bab5519e2e7d815a2469b76b5e6bf56c1051264f9b5b0a75a84e179a21b7701de8b744612ecd96b5e73f2ad4ffda4dde5a0bf0fa5b4490bd0c7ec41b331068f3db152278f5147c876201e741a6817616ece7c58a6507a7736c1fc341bb3ab65cc6e7196855a42eed503f04b56b190fced87eab134400c9fdcb1eb43fc7fed2882b2f56b9eea62ce8a024bca4f23aa4d70afb323d4c0ad3a38d409012207bb35e174a112794008d2c3f8a0d7f4282ab718493096da30d5e432f7917f017e4ee80191990aed9a51d404700c1e441ef3c46e83129aa2f5b4a1047757dc4ce4c11d1ea429c7a95dd95bc13f7c9fd5b4c64c5aa97040948142a72e57ef4658bf2894029fc69dcd893fe5bf72d90aced60e2b4608b0bfa6a06f26414c843a86df58d95f92c1904565898262d1170ad70252445bd883ec208415ef350cb0515a602d37cbb668d78e6f6211fa4caf338513c5e551f3a36b33214c89e9681301b830da28be02204d062ca19b2edacd56fa5ce4c7e1d0a9f1fd85f7049fe27dffc601b41f35dce8b0f61b3c92a8f51ab40299e6bf452c81d95ee1880dede9a6da64b3237451715c8da5970296d3b34b0c9b585b355f31e2b71c46cd49fe004d2ec5371b3029ee2d6d0881d90d73ac81b1d16a82a74f46e36b14e33a6abaf35b81fdccbe00031d8c5918974f53d35973cd7077b839c2dbfcead70236581065006dbb5f1e1541fa6226e172e0e9471a7a0a1ed5aa627d26e9aac140f0b2ddee38a4502fe9f6327e81fdb849cd7c7698e9add48aecab22f512b56fd0b"),
            String::from("5e50cca6bdd5d8c04e1a2848d74d885647d93b883cb4f182fbb5e3bdbf00506c"),
        );
        // when
        let is_verified = verify_bulletproof(&proof, 32).unwrap();
        // then
        assert!(is_verified);
    }
}
