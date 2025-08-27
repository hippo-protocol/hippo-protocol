import { useState } from "react";
import "./App.css";
import { create_keypair,  decrypt,  Did,  encrypt,  EncryptedData,  key_to_did,  KeyPair, sign, verify } from "hippo-sdk";

function App() {
  const [keypair, setKeypair] = useState<KeyPair>();

  const createKeypair = async () => {
    const generated = create_keypair();
    setKeypair(generated.to_object());
  };

  const [did, setDid] = useState<Did>();

  const createDid = async () => {
    if (keypair) {
      const generated = key_to_did(keypair.pubkey);
      setDid(generated.to_object());
    }
  };

  const [rawData, setRawData] = useState('');
  const [encData, setEncData] = useState<EncryptedData>();
  const [dataToDec, setEncDataToDec] = useState('');
  const [decData, setDecData] = useState('');

  const handleEncrypt = () => {
    try {
      if (rawData && keypair) {
        const generated = encrypt(rawData, keypair.pubkey);
        setEncData(generated.to_object());
      }
    } catch (e) {
      console.error(e);
      setEncData(undefined);
    }
  };

  const handleDecrypt = () => {
    try {
      if (dataToDec && keypair) {
        const fromObj = EncryptedData.from_object(JSON.parse(dataToDec));
        const generated = decrypt(fromObj, keypair.privkey);
        setDecData(generated);
      }
    } catch (e) {
      console.error(e);
      setDecData('');
    }
  };

  const [msgData, setMsgData] = useState('');
  const [sigData, setSigData] = useState<{data: string, sig: string}>();
  const [sigDataToVer, setSigToVer] = useState('');
  const [isVerified, setIsVerified] = useState<boolean>();

  const handleSign = () => {
    try {
      if (msgData && keypair) {
        const generated = sign(msgData, keypair.privkey);
        const sigWithData = {
          data: msgData,
          sig: generated,
        };
        setSigData(sigWithData);
      }
    } catch (e) {
      console.error(e);
      setSigData(undefined);
    }
  };

  const handleVerify = () => {
    try {
      if (sigDataToVer && keypair) {
        const sigWithData : {data: string, sig: string} = JSON.parse(sigDataToVer)
        const generated = verify(sigWithData.data, sigWithData.sig, keypair.pubkey);
        setIsVerified(generated);
      }
    } catch (e) {
      console.error(e);
      setIsVerified(false);
    }
  };


  return (
    <>
      <h1>Hippo SDK Example</h1>
      <div className="card">
        <button onClick={createKeypair}>
          Create Keypair
        </button>
        <div>
          {keypair? (
            <div style={{ whiteSpace: "pre-line" }}>
              {
                    JSON.stringify(keypair, null, 2)
               }
            </div>
          ) : (
            <span>No keypair generated yet</span>
          )}
        </div>
        <button onClick={createDid}>
          Create DID
        </button>
        <div>
          {did? (
            <div style={{ whiteSpace: "pre-line" }}>
              {
                   JSON.stringify(did, null, 2)
               }
            </div>
          ) : (
            <span>No DID generated yet</span>
          )}
        </div>
        <button onClick={handleEncrypt}>
          Encrypt Data with Public Key
        </button>
        <br></br>
        <textarea
          placeholder="Enter data to encrypt here..."
          value={rawData}
          onChange={(e) => setRawData(e.target.value)}
        />
        <div>
          {encData? (
            <div style={{ whiteSpace: "pre-line" }}>
              {
                    JSON.stringify(encData, null, 2)
               }
            </div>
          ) : (
            <span>Not encrypted yet</span>
          )}
        </div>
        <button onClick={handleDecrypt}>
          Decrypt Data with Private Key
        </button>
        <br></br>
        <textarea
          placeholder="Enter data to decrypt here..."
          value={dataToDec}
          onChange={(e) => setEncDataToDec(e.target.value)}
        />
        <div>
          {decData? (
            <div style={{ whiteSpace: "pre-line" }}>
              {
                    JSON.stringify(decData, null, 2)
               }
            </div>
          ) : (
            <span>Not decrypted yet</span>
          )}
        </div>
        <button onClick={handleSign}>
          Sign Data with Private Key
        </button>
        <br></br>
        <textarea
          placeholder="Enter data to sign here..."
          value={msgData}
          onChange={(e) => setMsgData(e.target.value)}
        />
        <div>
          {sigData? (
            <div style={{ whiteSpace: "pre-line" }}>
              {
                    JSON.stringify(sigData, null, 2)
               }
            </div>
          ) : (
            <span>Not signed yet</span>
          )}
        </div>
        <button onClick={handleVerify}>
          Verify Signature with Public Key
        </button>
        <br></br>
        <textarea
          placeholder="Enter signature to verify here..."
          value={sigDataToVer}
          onChange={(e) => setSigToVer(e.target.value)}
        />
        <div>
          { isVerified !== undefined ? (
            <div style={{ whiteSpace: "pre-line" }}>
              {
                    JSON.stringify(isVerified, null, 2)
               }
            </div>
          ) : (
            <span>Not verifed yet</span>
          )}
        </div>

      </div>
    </>
  );
}

export default App;
