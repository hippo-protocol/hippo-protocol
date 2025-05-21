import { useEffect, useState } from "react";
import "./App.css";
import init, { create_keypair, KeyPair } from "./assets/pkg/core.js";

function App() {
  const [keypair, setKeypair] = useState<KeyPair[]>([]);
  const [initialized, setInitialized] = useState(false);

  const createKeypair = async () => {
    const generated = create_keypair();
    setKeypair([...keypair, generated]);
  };

  useEffect(() => {
    init().then(() => {
      setInitialized(true);
    });
  });

  return (
    <>
      <script src="pkg/core.js" />
      <h1>Hippo SDK Example</h1>
      <div className="card">
        <button onClick={createKeypair} disabled={!initialized}>
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
