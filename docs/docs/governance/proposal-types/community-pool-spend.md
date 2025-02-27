---
order: 3
parent:
  order: 1
---

# Community Pool Spend

Hippo Protocol launched with community-spend capabilities, effectively unlocking the potential for token-holders to vote to approve spending from the Community Pool.

## Learn About the Community Pool

### How is the Community Pool funded?

92% of all staking rewards generated (via block rewards & transaction fees) are continually transferred to and accrue within the Community Pool.

### How can funding for the Community Pool change?

Though the rate of funding is currently fixed at 92% of staking rewards, the effective rate is dependent upon the Hippo Protocol's staking rewards, which can change with inflation and block times.

The current parameter `Community Tax` parameter of 92% may be modified with a governance proposal and enacted immediately after the proposal passes.

### How much money is in the Community Pool?

You may directly query the Hippo Protocol for the balance of the Community Pool:

`hippod q distribution community-pool --chain-id hippo-protocol-1 --node <rpc-node-address> `

Alternatively, popular Hippo explorers such as [Hippo River](https://river.hippoprotocol.ai/) display the ongoing Community Pool balance.

### How can funds from the Community Pool be spent?

Funds from the Hippo Community Pool may be spent via successful governance proposal.

### How are funds disbursed after a community-spend proposal is passed?

If a community-spend proposal passes successfully, the number of HP encoded in the proposal will be transferred from the community pool to the address encoded in the proposal, and this will happen immediately after the voting period ends.

## Why create a proposal to use Community Pool funds?

Why create a community-spend proposal?

**As a strategy: funding is fast.** Besides the time it takes to push your proposal on-chain, the only other limiting factor is a fixed 14-day voting period. As soon as the proposal passes, your account will be credited the full amount of your proposal request.

**To build rapport.** Engaging publicly with the community is the opportunity to develop relationships with stakeholders and to educate them about the importance of your work. Unforeseen partnerships could arise, and overall the community may value your work more if they are involved as stakeholders.

**To be more independent.** Having a more consistently funded source and having a report with its stakeholders means you can use your rapport to have confidence in your ability to secure funding independent from any.
