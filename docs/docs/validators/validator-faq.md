---
title: Validator FAQ
order: 3
---

:::warning

### Disclaimer

This is work in progress. Mechanisms and values are susceptible to change.
:::

## General Concepts

### What is a Hippo validator?

The Hippo Protocol is based on [CometBFT](https://docs.cometbft.com/v0.37/introduction/what-is-cometbft) that relies on a set of validators to secure the network. The role of validators is to run a full node and participate in consensus by broadcasting votes that contain cryptographic signatures signed by the validator's private key. Validators commit new blocks in the blockchain and receive revenue in exchange for their work. Validators must also participate in governance by voting on proposals. Validators are weighted according to their total stake.

### What is staking?

The Hippo Protocol is a public Proof-Of-Stake (PoS) blockchain, meaning that the weight of validators is determined by the amount of staking tokens (HP) bonded as collateral. These HP tokens can be self-delegated directly by the validator or delegated to the validator by other HP holders.

Any user in the system can declare their intention to become a validator by sending a `create-validator` transaction to become validator candidates.

The weight (i.e. voting power) of a validator determines whether they are an active validator. The active validator set is limited to [an amount](https://river.hippoprotocol.ai/hippo-protocol/staking) that changes over time.

### What is a full node?

A full node is a server running a chain's _binary_ (its software) that fully validates transactions and blocks of a blockchain and keeps a full record of all historic activity. A full node is distinct from a pruned node that processes only block headers and a small subset of transactions. Running a full node requires more resources than a pruned node. Validators can decide to run either a full node or a pruned node, but they need to make sure they retain enough blocks to be able to validate new blocks.

Of course, it is possible and encouraged for users to run full nodes even if they do not plan to be validators.

You can find more details about the requirements in the [Joining Mainnet Tutorial](../hub-tutorials/join-mainnet.md).

### What is a delegator?

Delegators are HP holders who cannot, or do not want to, run a validator themselves. HP holders can delegate HP to a validator and obtain a part of their revenue in exchange. For details on how revenue is distributed, see [What is the incentive to stake?](./validator-faq#what-is-the-incentive-to-stake) and [What are validators commission?](./validator-faq#what-is-a-validator-commission) in this document.

Because delegators share revenue with their validators, they also share risks. If a validator misbehaves, each of their delegators are partially slashed in proportion to their delegated stake. This penalty is one of the reasons why delegators must perform due diligence on validators before delegating. Spreading their stake over multiple validators is another layer of protection.

Delegators play a critical role in the system, as they are responsible for choosing validators. Being a delegator is not a passive role. Delegators must actively monitor the actions of their validators and participate in governance. For details on being a delegator, read the [Delegator FAQ](../delegators/delegator-faq.md).

## Becoming a Validator

### How to become a validator?

Any participant in the network can signal that they want to become a validator by sending a `create-validator` transaction, where they must fill out the following parameters:

- **Validator's `PubKey`:** The private key associated with this Tendermint/CometBFT `PubKey` is used to sign _prevotes_ and _precommits_.
- **Validator's Address:** Application level address that is used to publicly identify your validator. The private key associated with this address is used to delegate, unbond, claim rewards, and participate in governance.
- **Validator's name (moniker)**
- **Validator's website (Optional)**
- **Validator's description (Optional)**
- **Initial commission rate**: The commission rate on block rewards and fees charged to delegators.
- **Maximum commission:** The maximum commission rate that this validator can charge. This parameter is fixed and cannot be changed after the `create-validator` transaction is processed.
- **Commission max change rate:** The maximum daily increase of the validator commission. This parameter is fixed cannot be changed after the `create-validator` transaction is processed.

After a validator is created, HP holders can delegate HP to them, effectively adding stake to the validator's pool. The total stake of an address is the combination of HP bonded by delegators and HP self-bonded by the validator.

From all validator candidates that signaled themselves, the 22 validators with the most total stake are the designated **validators**. If a validator's total stake falls below the top 22, then that validator loses its validator privileges. The validator cannot participate in consensus or generate rewards until the stake is high enough to be in the top 22. Over time, the maximum number of validators may be increased via on-chain governance proposal.

## Additional Concepts

### What are the different types of keys?

There are two types of keys:

- **Tendermint/CometBFT key**: A unique key that is used to sign consensus votes.
  - It is associated with a public key `hippovalconspub` (To get this value, run `hippod tendermint show-validator`)
  - It is generated when the node is created with `hippod init`.
- **Application key**: This key is created from the `hippod` binary and is used to sign transactions. Application keys are associated with a public key that is prefixed by `hippopub` and an address that is prefixed by `hippo`.

The Tendermint/CometBFT key and the application key are derived from account keys that are generated by the `hippod keys add` command.

**Note:** A validator's operator key is directly tied to an application key and uses the `hippovaloper` and `hippovaloperpub` prefixes that are reserved solely for this purpose.

### What are the different states a validator can be in?

After a validator is created with a `create-validator` transaction, the validator is in one of three states:

- `in validator set`: Validator is in the active set and participates in consensus. The validator is earning rewards and can be slashed for misbehavior.
- `jailed`: Validator misbehaved and is in jail, i.e. outside of the validator set.

  - If the jailing is due to being offline for too long (i.e. having missed more than `25%` out of the last `10,000` blocks), the validator can send an `unjail` transaction in order to re-enter the validator set.
  - If the jailing is due to double signing, the validator cannot unjail.

- `unbonded`: Validator is not in the active set, and therefore not signing blocks. The validator cannot be slashed and does not earn any reward. It is still possible to delegate HP to an unbonded validator. Undelegating from an `unbonded` validator is immediate, meaning that the tokens are not subject to the unbonding period.

### What is self-delegation? How can I increase my self-delegation?

Self-delegation is a delegation of HP from a validator to themselves. The delegated amount can be increased by sending a `delegate` transaction from your validator's `application` application key.

### What is validator bond? How can I increase my validator bond?

Validator bond is a delegation of HP from a delegator to a validator. Validator operators can validator bond to themselves. The validator bond amount can be increased by sending a `ValidatorBond` transaction from any account delegated to your validator. Validator bond is required before a validator can accept delegations from liquid staking providers. As such it forces validators to put “skin in the game” in order to be entrusted with delegations from liquid staking providers. This disincentivizes malicious behavior and enables the validator to negotiate its relationship with liquid staking providers.

### Is there a minimum amount of HP that must be delegated to be an active (bonded) validator?

The minimum is 1 HP. But the network is currently secured by much higher values. You can check the minimum required HP to become part of the active validator set on the [Hippo River validator page](https://river.hippoprotocol.ai/hippo-protocol/staking).

### How do delegators choose their validators?

Delegators are free to choose validators according to their own subjective criteria. Selection criteria includes:

- **Amount of validator-bonded HP:** Number of HP a validator validator-bonded to themselves. A validator with a higher amount of self-delegated HP indicates that the validator is sharing the risk and consequences for their actions, or has enough goodwill from the community so that others post validator bond on the validator's behalf.
- **Amount of delegated HP:** Total number of HP delegated to a validator. A high voting power shows that the community trusts this validator. Larger validators also decrease the decentralization of the network, so delegators are suggested to consider delegating to smaller validators.
- **Commission rate:** Commission applied on revenue by validators before the revenue is distributed to their delegators.
- **Track record:** Delegators review the track record of the validators they plan to delegate to. This track record includes past votes on proposals and historical average uptime.
- **Community contributions:** Another (more subjective) criteria is the work that validators have contributed to the community, such as educational content, participation in the community channels, contributions to open source software, etc.

Apart from these criteria, validators send a `create-validator` transaction to signal a website address to complete their resume. Validators must build reputation one way or another to attract delegators. For example, a good practice for validators is to have a third party audit their setup. Note though, that the CometBFT team does not approve or conduct any audits themselves. For more information on due diligence, see the [A Delegator’s Guide to Staking](https://medium.com/@interchain_io/3d0faf10ce6f) blog post.

## Responsibilities

### Do validators need to be publicly identified?

No, they do not. Each delegator can value validators based on their own criteria. Validators are able to register a website address when they nominate themselves so that they can advertise their operation as they see fit. Some delegators prefer a website that clearly displays the team operating the validator and their resume, while other validators might prefer to be anonymous validators with positive track records.

### What are the responsibilities of a validator?

Validators have two main responsibilities:

- **Be able to constantly run a correct version of the software:** Validators must ensure that their servers are always online and their private keys are not compromised.

- **Actively participate in governance:** Validators are required to vote on every proposal.

Additionally, validators are expected to be active members of the community. Validators must always be up-to-date with the current state of the ecosystem so that they can easily adapt to any change.

### What does 'participate in governance' entail?

Validators and delegators on the Hippo Protocol can vote on proposals to change operational parameters (such as the block gas limit), coordinate upgrades, or make a decision on any given matter.

Validators play a special role in the governance system. As pillars of the system, validators are required to vote on every proposal. It is especially important since delegators who do not vote inherit the vote of their validator.

### What does staking imply?

Staking HP can be thought of as a safety deposit on validation activities. When a validator or a delegator wants to retrieve part or all of their deposit, they send an `unbonding` transaction. Then, HP undergoes a **3-week unbonding period** during which they are liable to being slashed for potential misbehaviors committed by the validator before the unbonding process started.

Validators, and by association delegators, receive block rewards, fees, and have the right to participate in governance. If a validator misbehaves, a certain portion of their total stake is slashed. This means that every delegator that bonded HP to this validator gets penalized in proportion to their bonded stake. Delegators are therefore incentivized to delegate to validators that they anticipate will function safely.

### Can a validator run away with their delegators' HP?

By delegating to a validator, a user delegates voting power. The more voting power a validator have, the more weight they have in the consensus and governance processes. This does not mean that the validator has custody of their delegators' HP. **A validator cannot run away with its delegator's funds**.

Even though delegated funds cannot be stolen by their validators, delegators' tokens can still be slashed by a small percentage if their validator suffers a [slashing event](#what-are-the-slashing-conditions), which is why we encourage due diligence when [selecting a validator](#how-do-delegators-choose-their-validators).

### How often is a validator chosen to propose the next block? Does frequency increase with the quantity of bonded HP?

The validator that is selected to propose the next block is called the proposer. Each proposer is selected deterministically. The frequency of being chosen is proportional to the voting power (i.e. amount of bonded HP) of the validator. For example, if the total bonded stake across all validators is 100 HP and a validator's total stake is 10 HP, then this validator is the proposer ~10% of the blocks.

### How can a validator safely quit validating on the Hippo Protocol?

If a validator simply shuts down their node, this would result in the validator and their delegators getting slashed for being offline. The only way to safely exit a validator node running on the Hippo Protocol is by unbonding the validator with the `UnbondValidator` message. As a result, the validator gets jailed and kicked out of the active set of validators, without getting slashed. They can then proceed to shut down their node without risking their tokens.

It's highly advised to inform your delegators when doing this, as they will still be bonded to your validator after it got jailed. They will need to manually unbond and they might not have been made aware of this via their preferred wallet application.

## Incentives

### What is the incentive to stake?

Each member of a validator's staking pool earns different types of revenue:

- **Block rewards:** Native tokens of applications (e.g. HP on the Hippo Protocol) run by validators are inflated to produce block provisions. These provisions exist to incentivize HP holders to bond their stake. Non-bonded HP are diluted over time.
- **Transaction fees:** The Hippo Protocol maintains an allow list of tokens that are accepted as fee payment. The initial fee token is the `ahp`.

This total revenue is divided among validators' staking pools according to each validator's weight. Then, within each validator's staking pool the revenue is divided among delegators in proportion to each delegator's stake. A commission on delegators' revenue is applied by the validator before it is distributed.

### What is a validator commission?

Revenue received by a validator's pool is split between the validator and their delegators. The validator can apply a commission on the part of the revenue that goes to their delegators. This commission is set as a percentage. Each validator is free to set their initial commission, maximum daily commission change rate, and maximum commission. The Hippo Protocol enforces the parameter that each validator sets. The maximum commission rate is fixed and cannot be changed. However, the commission rate itself can be changed after the validator is created as long as it does not exceed the maximum commission.

### What is the incentive to run a validator?

Validators earn proportionally more revenue than their delegators because of the commission they take on the staking rewards from their delegators.

Validators also play a major role in governance. If a delegator does not vote, they inherit the vote from their validator. This voting inheritance gives validators a major responsibility in the ecosystem.

### How are block rewards distributed?

Block rewards are distributed proportionally to all validators relative to their voting power. This means that even though each validator gains HP with each reward, all validators maintain equal weight over time.

For example, 10 validators have equal voting power and a commission rate of 1%. For this example, the reward for a block is 1000 HP and each validator has 20% of self-bonded HP. These tokens do not go directly to the proposer. Instead, the tokens are evenly spread among validators. So now each validator's pool has 100 HP. These 100 HP are distributed according to each participant's stake:

- Commission: `100*80%*1% = 0.8 HP`
- Validator gets: `100\*20% + Commission = 20.8 HP`
- All delegators get: `100\*80% - Commission = 79.2 HP`

Then, each delegator can claim their part of the 79.2 HP in proportion to their stake in the validator's staking pool.

### How are fees distributed?

Fees are similarly distributed with the exception that the block proposer can get a bonus on the fees of the block they propose if the proposer includes more than the strict minimum of required precommits.

When a validator is selected to propose the next block, the validator must include at least 2/3 precommits of the previous block. However, an incentive to include more than 2/3 precommits is a bonus. The bonus is linear: it ranges from 1% if the proposer includes 2/3rd precommits (minimum for the block to be valid) to 5% if the proposer includes 100% precommits. Of course the proposer must not wait too long or other validators may timeout and move on to the next proposer. As such, validators have to find a balance between wait-time to get the most signatures and risk of losing out on proposing the next block. This mechanism aims to incentivize non-empty block proposals, better networking between validators, and mitigates censorship.

### What are the slashing conditions?

If a validator misbehaves, their delegated stake is partially slashed. Two faults can result in slashing of funds for a validator and their delegators:

- **Double signing:** If someone reports on chain A that a validator signed two blocks at the same height on chain A and chain B, and if chain A and chain B share a common ancestor, then this validator gets slashed by 5% on chain A.
- **Downtime:** If a validator misses more than `95%` of the last `10,000` blocks (roughly ~19 hours), they are slashed by 0.01%.

### Are validators required to self-delegate HP?

No, they do not need to self-delegate. Even though there is no obligation for validators to self-delegate, delegators may want their validator to have self-delegated HP in their staking pool. In other words, validators share the risk.

Note however that it's possible that some validators decide to self-delegate via a different address for security reasons.

### How to prevent concentration of stake in the hands of a few top validators?

The community is expected to behave in a smart and self-preserving way. When a mining pool in Bitcoin gets too much mining power the community usually stops contributing to that pool. The Hippo Protocol relies on the same effect. Additionally, when delegators switch to another validator, they are not subject to the unbonding period, which removes any barrier to quickly redelegating tokens in service of improving decentralization.

### Who can validator bond?

The validator themselves, but also any other address delegated to the validator.

### How can I validator bond?

Once delegated to a validator, a delegator (or validator operator) can convert their delegation to a validator into Validator Bond by signing a ValidatorBond message.

The ValidatorBond message is exposed by the staking module and can be executed as follows:

```
hippod tx staking validator-bond hippovaloper13h5xdxhsdaugwdrkusf8lkgu406h8t62jkqv3h <delegator> --from mykey
```

There are no partial Validator Bonds: when a delegator or validator converts their shares to a particular validator into Validator Bond, their entire delegation to that validator is converted to Validator Bond. If a validator or delegator wishes to convert only some of their delegation to Validator Bond, they should transfer those funds to a separate address and Validator Bond from that address, or redelegate the funds that they do not wish to validator bond to another validator before converting their delegation to validator bond.

To convert Validator Bond back into a standard delegation, simply unbond the shares.

### How does a delegator or validator mark their delegation as a validator bond?

Once delegated to a validator, sign a `ValidatorBond` message.

### Are validator bonds subject to additional slashing conditions?

No, in the event of a slash, a validator bond is slashed at the same rate as a regular bond.

### Can I validator bond some of my tokens and delegate the remaining portion normally?

The `ValidatorBond` message converts the full balance delegated to a validator into validator bond. To validator bond some tokens and delegate the remaining portion normally, use two addresses: the first will delegate + ValidatorBond, and the second will just delegate.

## Technical Requirements

### What are hardware requirements?

A modest level of hardware specifications is initially required and rises as network use increases. Participating in the testnet is the best way to learn more. You can find the current hardware recommendations in the [Joining Mainnet documentation](../hub-tutorials/join-mainnet.md).

Validators are recommended to set up [sentry nodes](https://docs.cometbft.com/v0.37/core/validators) to protect your validator node from DDoS attacks.

### What are software requirements?

In addition to running a Hippo Protocol node, validators are expected to implement monitoring, alerting, and management solutions. There are [several tools](https://medium.com/solar-labs-team/cosmos-how-to-monitoring-your-validator-892a46298722) that you can use.

### What are bandwidth requirements?

The Hippo network has the capacity for very high throughput relative to chains like Ethereum or Bitcoin.

We recommend that the data center nodes connect only to trusted full nodes in the cloud or other validators that know each other socially. This connection strategy relieves the data center node from the burden of mitigating denial-of-service attacks.

Ultimately, as the network becomes more heavily used, multigigabyte per day bandwidth is very realistic.

### How to handle key management?

Validators are expected to run an HSM that supports ed25519 keys. Here are potential options:

- YubiHSM 2
- Ledger Nano S
- Ledger BOLOS SGX enclave
- Thales nShield support

The Interchain Foundation does not recommend one solution above the other. The community is encouraged to bolster the effort to improve HSMs and the security of key management.

### What can validators expect in terms of operations?

Running an effective operation is key to avoiding unexpected unbonding or slashing. Operations must be able to respond to attacks and outages, as well as maintain security and isolation in the data center.

### What are the maintenance requirements?

Validators are expected to perform regular software updates to accommodate chain upgrades and bug fixes. It is suggested to consider using [Cosmovisor](https://docs.cosmos.network/v0.45/run-node/cosmovisor.html) to partially automate this process.

### How can validators protect themselves from denial-of-service attacks?

Denial-of-service attacks occur when an attacker sends a flood of internet traffic to an IP address to prevent the server at the IP address from connecting to the internet.

An attacker scans the network, tries to learn the IP address of various validator nodes, and disconnects them from communication by flooding them with traffic.

One recommended way to mitigate these risks is for validators to carefully structure their network topology using a sentry node architecture.

Validator nodes are expected to connect only to full nodes they trust because they operate the full nodes themselves or the trust full nodes are run by other validators they know socially. A validator node is typically run in a data center. Most data centers provide direct links to the networks of major cloud providers. The validator can use those links to connect to sentry nodes in the cloud. This mitigation shifts the burden of denial-of-service from the validator's node directly to its sentry nodes, and can require that new sentry nodes are spun up or activated to mitigate attacks on existing ones.

Sentry nodes can be quickly spun up or change their IP addresses. Because the links to the sentry nodes are in private IP space, an internet-based attack cannot disturb them directly. This strategy ensures that validator block proposals and votes have a much higher chance to make it to the rest of the network.

For more sentry node details, see the [CometBFT Documentation](https://docs.cometbft.com/v0.37/core/validators) or the [Sentry Node Architecture Overview](https://forum.cosmos.network/t/sentry-node-architecture-overview/454) on the forum.
