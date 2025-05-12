try {
  const {
    create_keypair,
    encrypt,
    decrypt,
    key_to_did,
  } = require("../pkg/core");

  const keyPair = create_keypair();
  console.log(keyPair.privkey, keyPair.pubkey);

  const message = "Hello, world!";
  console.log(
    "encrypt and decrypt a message: ",
    message,
    decrypt(encrypt(message, keyPair.pubkey), keyPair.privkey)
  );

  const did = key_to_did(keyPair.pubkey);
  console.log("did:", did.id);
} catch (e) {
  if (e.code === "MODULE_NOT_FOUND") {
    console.log("You shoulde build the sdk using wasm-pack");
  } else {
    console.error(e);
  }
}
