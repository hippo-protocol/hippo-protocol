---
title: Service Providers
order: 4
---

'Service Providers' are defined as entities that provide services for end-users that involve some form of interaction with the Hippo Protocol. More specifically, this document is focused on interactions with tokens.

Service Providers are expected to act as trusted points of contact to the blockchain for their end-users. This Service Providers section does not apply to wallet builders that want to provide Light Client functionalities.

This document describes:

- [Connection Options](#connection-options)
- [Running a Full Node](#running-a-full-node)
  - [What is a Full Node?](#what-is-a-full-node)
  - [Installation and Configuration](#installation-and-configuration)
- [Command-Line Interface](#command-line-interface)
  - [Available Commands](#available-commands)
  - [Remote Access to hippod](#remote-access-to-hippod)
  - [Create a Key Pair](#create-a-key-pair)
    - [Check your Account](#check-your-account)
  - [Check your Balance](#check-your-balance)
    - [Send Coins Using the CLI](#send-coins-using-the-cli)
- [REST API](#rest-api)
  - [Listen for Incoming Transactions](#listen-for-incoming-transactions)

## Connection Options

There are four main technologies to consider to connect to the Hippo Protocol:

- Full Nodes: Interact with the blockchain.
- REST Server: Serves for HTTP calls.
- REST API: Use available endpoints for the REST Server.
- GRPC: Connect to the Hippo Protocol using gRPC.

## Running a Full Node

### What is a Full Node?

A Full Node is a network node that syncs up with the state of the blockchain. It provides blockchain data to others by using RESTful APIs, a replica of the database by exposing data with interfaces. A Full Node keeps in syncs with the rest of the blockchain nodes and stores the state on disk. If the full node does not have the queried block on disk the full node can go find the blockchain where the queried data lives.

### Installation and Configuration

This section describes the steps to run and interact with a full node for the Hippo Protocol.

First, you need to [install the software](../getting-started/installation).

Consider running your own [Hippo Protocol Full Node](../hub-tutorials/join-mainnet).

## Command-Line Interface

The command-line interface (CLI) is the most powerful tool to access the Hippo Protocol and use hippo.
To use the CLI, you must install the latest version of `hippo` on your machine.

Compare your version with the [latest release version](https://github.com/hippocrat-dao/hippo-protocol.git)

```bash
hippod version --long
```

### Available Commands

All available CLI commands are shown when you run the `hippod` command:

```bash
hippod
```

````bash
Hippo App

Usage:
hippod [command]

Available Commands:
config Create or query an application CLI configuration file
debug Tool for helping with debugging your application
export Export state to JSON
genesis Application's genesis-related subcommands
help Help about any command
init Initialize private validator, p2p, genesis, and application configuration files
keys Manage your application's keys
prune Prune app history states by keeping the recent heights and deleting old heights
query Querying subcommands
rollback rollback cosmos-sdk and tendermint state by one height
snapshots Manage local snapshots
start Run the full node
status Query remote node for status
tendermint Tendermint subcommands
tx Transactions subcommands
version Print the application binary version information

Flags:
-h, --help help for hippod
--home string directory for config and data (default "/home/kek0114/.hippo")
--log_format string The logging format (json|plain) (default "plain")
--log_level string The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
--log_no_color Disable colored logs
--trace print out full stack trace on errors

Use "hippod [command] --help" for more information about a command.```

For each displayed command, you can use the `--help` flag to get further information.

```bash
hippod query --help
Usage:
  hippod query [flags]
  hippod query [command]

Aliases:
  query, q

Available Commands:
  account                  Query for account by address
  auth                     Querying commands for the auth module
  bank                     Querying commands for the bank module
  block                    Get verified data for a the block at given height
  distribution             Querying commands for the distribution module
  evidence                 Query for evidence by hash or for all (paginated) submitted evidence
  gov                      Querying commands for the governance module
  ibc                      Querying commands for the IBC module
  ibc-transfer             IBC fungible token transfer query subcommands
  mint                     Querying commands for the minting module
  params                   Querying commands for the params module
  slashing                 Querying commands for the slashing module
  staking                  Querying commands for the staking module
  tendermint-validator-set Get the full tendermint validator set at given height
  tx                       Query for a transaction by hash in a committed block
  txs                      Query for paginated transactions that match a set of events
  upgrade                  Querying commands for the upgrade module

Flags:
      --chain-id string   The network chain ID
  -h, --help              help for query

Global Flags:
      --home string         directory for config and data (default "/Users/tobias/.hippo")
      --log_format string   The logging format (json|plain) (default "plain")
      --log_level string    The logging level (trace|debug|info|warn|error|fatal|panic) (default "info")
      --trace               print out full stack trace on errors

Use "hippod query [command] --help" for more information about a command.
````

### Remote Access to hippod

When choosing to remote access a Full Node and hippod, you need a Full Node running and hippo installed on your local machine.

`hippod` is the tool that enables you to interact with the node that runs on the Hippo Protocol network, whether you run it yourself or not.

To set up `hippod` on a local machine and connect to an existing full node, use the following command:

```bash
hippod config <flag> <value>
```

First, set up the address of the full node you want to connect to:

```bash
hippod config node <host>:<port

// example: hippod config node https://77.87.106.33:26657 (note: this is a placeholder)
```

If you run your own full node locally, use `tcp://localhost:26657` as the address.

Finally, set the `chain-id` of the blockchain you want to interact with:

```bash
hippod config chain-id hippo-protocol-1
```

Next, learn to use CLI commands to interact with the full node.
You can run these commands as remote control or when you are running it on your local machine.

### Create a Key Pair

The default key is `secp256k1 elliptic curve`. Use the `hippod keys` command to list the keys and generate a new key.

```bash
hippod keys add <your_key_name>
```

You will be asked to create a password (at least 8 characters) for this key-pair. This will return the information listed below:

- `NAME`: Name of your key
- `TYPE`: Type of your key, always `local`.
- `ADDRESS`: Your address. Used to receive funds.
- `PUBKEY`: Your public key. Useful for validators.
- `MNEMONIC`: 24-word phrase. **Save this mnemonic somewhere safe**. This phrase is required to recover your private key in case you forget the password. The mnemonic is displayed at the end of the output.

You can see all available keys by typing:

```bash
hippod keys list
```

Use the `--recover` flag to add a key that imports a mnemonic to your keyring.

```bash
hippod keys add <your_key_name> --recover
```

#### Check your Account

You can view your account by using the `query account` command.

```bash
hippod query account <YOUR_ADDRESS>
```

It will display your account type, account number, public key and current account sequence.

```bash
'@type': /cosmos.auth.v1beta1.BaseAccount
account_number: "xxxx"
address: hippoxxxx
pub_key:
  '@type': /cosmos.crypto.secp256k1.PubKey
  key: xxx
sequence: "x"
```

### Check your Balance

Query the account balance with the command:

```bash
hippod query bank balances <YOUR_ADDRESS>
```

The response contains keys `balances` and `pagination`.
Each `balances` entry contains an `amount` held, connected to a `denom` identifier.
The typical $HP token is identified by the denom `ahp`. Where 1 `ahp` is 0.000001 HP.

```bash
balances:
- amount: "12345678"
  denom: ahp
pagination:
  next_key: null
  total: "0"
```

When you query an account that has not received any token yet, the `balances` entry is shown as an empty array.

```bash
balances: []
pagination:
  next_key: null
  total: "0"
```

#### Send Coins Using the CLI

To send coins using the CLI:

```bash
hippod tx bank send [from_key_or_address] [to_address] [amount] [flags]
```

Parameters:

- `<from_key_or_address>`: Key name or address of sending account.
- `<to_address>`: Address of the recipient.
- `<amount>`: This parameter accepts the format `<value|coinName>`, such as `1000000ahp`.

Flags:

- `--chain-id`: This flag allows you to specify the id of the chain. There are different ids for different testnet chains and mainnet chains.
- `--gas-prices`: This flag allows you to specify the gas prices you pay for the transaction. The format is used as `0.0025ahp`

## REST API

The REST API documents list all the available endpoints that you can use to interact
with your full node. Learn [how to enable the REST API](../hub-tutorials/join-mainnet.md#enable-the-rest-api) on your full node.

### Listen for Incoming Transactions

The recommended way to listen for incoming transactions is to periodically query the blockchain by using the following HTTP endpoint:

`/cosmos/bank/v1beta1/balances/{address}`
