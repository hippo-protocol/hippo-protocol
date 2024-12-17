# Hippo Protocol

`Hippo Protocol` is an application chain built using the Cosmos SDK for medical purposes.

## Running testnets with `hippod`

If you want to spin up a quick testnet with your friends, you can follow these steps.
Unless otherwise noted, every step must be done by everyone who wants to participate
in this testnet.

1. If you've run `hippod` before, you may need to reset your genesis file(including gentx) and database before starting a new
   testnet. You can do with the following command: `$ rm -rf .hippo` from your
   Home(or User) directory. If you want to reset only database, just do `$ go run hippod/main.go tendermint unsafe-reset-all` in root dir.
2. `$ go run hippod/main.go init hippo --chain-id hippo-protocol-testnet-1`. This will initialize a new working directory
   at the default location `~/.hippod`. You need to provide a "moniker" and a "chain id". These
   two names are "hippo" and "hippo-protocol-testnet-1" here. If you want to overwrite just genesis file(not including gentx), add `--overwrite` flag.
3. `$ go run hippod/main.go keys add alice`. This will create a new key, with a name of your choosing(for here, alice).
   Save the output of this command somewhere; you'll need the address generated here later.
4. `$ go run hippod/main.go genesis add-genesis-account alice 1084734273380000000000000000ahp`, where `key_name` is the same key name as
   before; and `1084734273380000000000000000ahp` is `amount`.
5. `$ go run hippod/main.go genesis gentx alice 1000000000000000000ahp --chain-id hippo-protocol-testnet-1`. This will create the genesis
   transaction for your new chain. Here `amount` should be at least `1000000000000000000ahp`. If you
   provide too much or too little, you will encounter an error when starting your node.
6. Now, one person needs to create the genesis file `genesis.json` using the genesis transactions
   from every participant, by gathering all the genesis transactions under `config/gentx` and then
   calling `$ go run hippod/main.go genesis collect-gentxs`. This will create a new `genesis.json` file that includes data
   from all the validators (we sometimes call it the "super genesis file" to distinguish it from
   single-validator genesis files).
7. Once you've received the super genesis file, overwrite your original `genesis.json` file with
   the new super `genesis.json`.
8. Modify your `config/config.toml` (in the simapp working directory) to include the other participants as
   persistent peers:

   ```text
   # Comma separated list of nodes to keep persistent connections to
   persistent_peers = "[validator_address]@[ip_address]:[port],[validator_address]@[ip_address]:[port]"
   ```

   You can find `validator_address` by running `$ go run hippod/main.go tendermint show-node-id`. The output will
   be the hex-encoded `validator_address`. The default `port` is 26656.

9. Now you can start your nodes: `$ go run hippod/main.go start`.

Now you have a hippod testnet that you can use to try out changes to the Cosmos SDK or Tendermint!

NOTE: Sometimes creating the network through the `collect-gentxs` will fail, and validators will start
in a funny state (and then panic). If this happens, you can try to create and start the network first
with a single validator and then add additional validators using a `create-validator` transaction.
