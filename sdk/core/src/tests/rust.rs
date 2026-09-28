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
        let tag = String::from("hippo");
        let proof = generate_bulletproof(1234, bits, tag.clone()).unwrap();
        let other_proof = generate_bulletproof(1234, bits, tag.clone()).unwrap();
        let mut tampered_proof = hex::decode(proof.proof()).unwrap();
        tampered_proof[0] ^= 1;
        // when
        let is_verified = verify_bulletproof(&proof, bits, tag.clone()).unwrap();
        let wrong_tag = verify_bulletproof(&proof, bits, String::from("wrong hippo")).unwrap();
        let wrong_bits = verify_bulletproof(&proof, 64, tag.clone()).unwrap();
        let unsupported_bits = verify_bulletproof(&proof, 10, tag.clone()).unwrap();
        let wrong_commitment = verify_bulletproof(
            // Same value, but commitment is from another proof.
            &Bulletproof::new(proof.proof(), other_proof.commitment()),
            bits,
            tag.clone(),
        )
        .unwrap();
        let wrong_proof = verify_bulletproof(
            &Bulletproof::new(hex::encode(tampered_proof), proof.commitment()),
            bits,
            tag,
        )
        .unwrap();
        // then
        assert!(is_verified);
        assert!(!wrong_tag);
        assert!(!wrong_bits);
        assert!(!unsupported_bits);
        assert!(!wrong_commitment);
        assert!(!wrong_proof);
    }

    #[test]
    fn test_bulletproof_range_bounds() {
        let tag = String::from("hippo");
        for bits in [8, 16, 32, 64] {
            // both ends of the range [0, 2^bits)
            for value in [0, u64::MAX >> (64 - bits)] {
                let proof = generate_bulletproof(value, bits, tag.clone()).unwrap();
                assert!(verify_bulletproof(&proof, bits, tag.clone()).unwrap());
            }
        }
    }

    #[test]
    fn test_bulletproof_known_proof() {
        // given: pre-generated proof(value 1037578891, 32 bits, tag "hippo"), so a change of proof format is caught.
        let proof = Bulletproof::new(
            String::from("e212dffd521fb6d3d8b50ce340282dae99550da537e974c44a14aed10fb6181056a6712ede1569a2e5322f9f7597541b9e4d4db86ec4e8d85ad7f81907e85322e80253330072ae2cbfc952499f13179823445cf670c970cd77a6c085f9bbb57d8c1521c815565be19e48562a6b8272e9d5ab632057e3501b22de200c28318a33a1de7c84bdc052aa856526c15b9aca47c28e3bcb1e46ae55a6327078eab45c022371548ae03c43b7eb6a34bcecd84036368c544f76848ba90289e8df08b6e10512034a06222a0126c72c2e891aa19274b80a66b8564dfb41e3a27a58aeb52207325875ed52ae98057782a316cea9bf0a5da35e8ecda7183a3d5e6ba8d0af66219802e8d59a16bf9e61972939d0088c694f4bc572822008329577a77435752307ac9c4e77c4f23e7183816148fcd9e26063150348bc30d083c0bac539183cf1350cb6ab3b0e16c869506b42254cc92b755a56ee6b95f553458a183f75377e1c1de260546d1dee299a03d662f28f26279b4bf1a091755a0aa7d3f204608a37f26860c55f9a5f8b4af37ebc573e0817df2a1b7e2fec4219cb13628c54b87717ee5066814d23d78bc6cf669bbb3d656d48781c53f806106c6de758c5a6450dd7a62f7649bb544e79db21bdd6b4043b89bc23fdb3dad7516e76964986e2b7ce1c6101dcd6ca57accc97d4565b3202d6c37e7f1c87119cdd632c8fd4c1e7b7695bbb2454a3baffbe506379246629949bbf9fc72a77a3b5d53caefdafd5aa7e8d98be5963d6778b3046aae2ff666f0aacc1488381d6917eaad65d489e8ae029ea52fc0a357ccd683140694e2934e8b4cb3bc8ae8ffdcf2cd2b8810bd37f25aac610b502"),
            String::from("30129e12b0a685521465275031da1d57b7386c7f6078644bc8ed88fb36d1ef51"),
        );
        // when
        let is_verified = verify_bulletproof(&proof, 32, String::from("hippo")).unwrap();
        // then
        assert!(is_verified);
    }
}
