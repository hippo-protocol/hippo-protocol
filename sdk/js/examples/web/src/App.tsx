import { useState } from "react";
import "./App.css";
import { create_keypair, KeyPair } from "hippo-sdk";

function App() {
  const [keypair, setKeypair] = useState<KeyPair[]>([]);

  const createKeypair = async () => {
    const generated = create_keypair();
    setKeypair([...keypair, generated]);
  };


  return (
    <>
      <h1>Hippo SDK Example</h1>
      <div className="card">
        <button onClick={createKeypair}>
          Create Keypair
        </button>
        <div>
          {keypair.length > 0 ? (
            <div style={{ whiteSpace: "pre-line" }}>
              {keypair
                .map(
                  (obj) =>
                    `privateKey: ${obj.privkey}\n publicKey: ${obj.pubkey}`
                )
                .join("\n")}
            </div>
          ) : (
            <span>No keypair generated yet</span>
          )}
        </div>
      </div>
    </>
  );
}

export default App;
