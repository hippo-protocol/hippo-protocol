---
title: Hippo Protocol Archives
order: 2
---

With each breaking upgrade of the Hippo Protocol, the network is restarted at height 0. During this process, an export of the last state of the previous network is made to produce the genesis state of the new one.

As a result, the blocks of the previous networks are not downloaded by new clients (as they sync from the new genesis state), and may be deleted by existing full-nodes.

In an effort to maintain transparency, the interchain hosts archives of the previous versions of the Hippo Protocol network. These archives can be found [here](https://archive.interchain.io/).

If you want to make archives available to the community, feel free to open a PR to this file and add them.
