#[cfg(test)]
mod tests {
    use crate::{
        create_keypair, decrypt, did_to_key, encrypt, key_to_did, pedersen_commit, pedersen_reveal,
        sign, types::Commitment, verify,
    };

    #[test]
    fn test_enc_dec() {
        // given
        let data = String::from("datag허ㅜㅏ니ㅜ2#@_!##ㅏ!~2ㅡ₩ㅡ1    ㅁAl;;A;;:{}()[]");
        let alice = create_keypair();
        // when
        let enc_data = encrypt(data.clone(), alice.pubkey());
        let dec_data = decrypt(enc_data, alice.privkey());
        // then
        assert_eq!(data, dec_data);
    }

    #[test]
    fn test_did() {
        // given
        let key_pair = create_keypair();
        // when
        let pubkey_to_did = key_to_did(key_pair.pubkey());
        let did_to_pubkey = did_to_key(pubkey_to_did.clone());
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
        let sig = sign(data.clone(), key_pair.privkey());
        let is_verified = verify(data, sig, key_pair.pubkey());
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
        let is_verified = pedersen_reveal(commitment.clone(), value, tag.clone());
        let wrong_value_and_tag = pedersen_reveal(commitment.clone(), wrong_value, wrong_tag);
        let wrong_blinding_factor_with_same_value = pedersen_reveal(
            Commitment::new(
                commitment_same_value.commitment(),
                // Value is same but blinding factor is different.
                commitment.secret_blinding_factor(),
            ),
            value,
            tag.clone(),
        );
        let wrong_commitment_with_correct_blinding_factor = pedersen_reveal(
            Commitment::new(
                // Even if the value is same, blinding factor makes the commitment different.
                commitment.commitment(),
                commitment_same_value.secret_blinding_factor(), // Blinding factor is correct.
            ),
            value,
            tag,
        );
        // then
        assert!(is_verified);
        assert!(!wrong_value_and_tag);
        assert!(!wrong_blinding_factor_with_same_value);
        assert!(!wrong_commitment_with_correct_blinding_factor)
    }
}
