# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

Hippo Protocol is a Cosmos SDK v0.50 application chain for healthcare data, with CosmWasm
(wasmd v0.54) and IBC v8. Binary `hippod`, chain name `hippo`, node home `~/.hippo`, native
denom `ahp` at **18 decimals**. Released versions are audited by Oak Security (reports linked
from README.md).

## Commands

### Build

```bash
go build ./...          # fast compile check
make build              # → ./build/hippod with ldflags; needs gcc, or LEDGER_ENABLED=false
make build-linux        # cross-compile, ledger off
```

`make` derives `VERSION` from `git describe --tags`. In a clone with no tags that fails and
the version string compiles in empty — harmless locally, but use a tagged clone for releases.

### Test

```bash
go test $(go list ./... | grep -v "/test/e2e")   # what CI runs
go test ./app/... -run TestNewApp -v             # single test
make test                                        # ./... + coverage → artifacts/coverage.html
```

`make test` runs `./...` **including** `test/e2e`, which needs a live node and fails
standalone. Use the filtered command for normal work.

CI enforces a **75% total coverage floor** (`.github/workflows/go.yml`). Adding untested code
can turn CI red even when every test you wrote passes.

### Lint and format

```bash
make lint      # golangci-lint v1.56.2
make format    # gofumpt + golangci-lint --fix
```

Some committed files are currently unformatted and **no CI job checks formatting** — run
`gofmt -l .` before committing rather than assuming the tree is clean.

### Run a local node

Full working sequence is the `e2e-genesis-test` job in `.github/workflows/go.yml`. Short form:

```bash
hippod init <moniker> --chain-id <id>
hippod keys add alice --keyring-backend file
hippod genesis add-genesis-account alice <amount>ahp --keyring-backend file
hippod genesis gentx alice <amount>ahp --chain-id <id> --keyring-backend file
hippod genesis collect-gentxs
hippod start
```

`min-gas-prices` is non-zero, so **every tx needs an explicit `--fees`** — see
`test/e2e/cli_test.go` for values that work. The REST API (1317) is disabled by default;
enable it in `~/.hippo/config/app.toml` before running `go test ./test/e2e`.

### Protobuf

`make proto-gen`, `make proto-lint` (both run in docker). Note there are currently no `.proto`
files to generate from — see below.

## Architecture

### There are no custom modules

This chain defines **no `x/` modules and no custom protobuf types**. `proto/` holds buf config
and codegen scripts but no `.proto` files. Every module is stock Cosmos SDK, IBC-go, or wasmd.
All customization is app-level: wiring, genesis parameters, and one inflation function. Don't
look for business logic in a module directory — there isn't one.

The baseline is **simapp + wasmd**, and the code says so: comments in `app/app.go`,
`hippod/cmd/root.go` and `go.mod` cite the exact upstream files they were copied from. When
something looks odd, diff it against `cosmos-sdk/simapp` or `wasmd/app` at the pinned versions
before assuming it's intentional — most of it is verbatim.

### Entry point

`hippod/main.go` → `hippod/cmd/root.go:NewRootCmd` → `app.New` in `app/app.go`.

`NewRootCmd` builds a **complete throwaway app in a temp dir on every CLI invocation**, even
`hippod keys list`, just to extract the encoding config. The temp dir exists because wasm locks
its data directory. The startup cost is expected, not a bug.

### Where customization lives

1. **`types/consensus/`** — chain constants. `policy.go` holds every economic and consensus
   parameter (block limits, inflation bounds, slashing, gov, staking, evidence age);
   `units.go` the denom and precision; `wallet.go` the bech32 prefixes and BIP44 path.
   Change a chain parameter here, never inline at the use site.
2. **`hippod/cmd/init.go:overrideGenesis`** — applies those constants to the genesis doc. This
   runs **only at `hippod init`**. It has no effect on a running chain; parameters on a live
   chain change by governance proposal.
3. **`app/keepers/`** — `keys.go` declares store keys, `keepers.go` constructs every keeper.
   Every keeper's authority is the gov module account; there is no admin key.
4. **`app/inflation.go`** — `CustomInflationCalculationFn`, the only novel consensus logic.
   A height-based halving schedule that **ignores `bondedRatio` entirely**, unlike the SDK
   default. Wired into the mint module in `app/app.go`.

`app/export.go`, `app/genesis.go`, most of `app/app.go` and `hippod/cmd/genaccounts.go` are
stock copies.

### Invariants to preserve

**The global SDK config is sealed.** `consensus.SetWalletConfig()` sets bech32 prefixes and
BIP44, then calls `Seal()`. Calling it twice in one process **panics** with `Config is sealed`.
Tests that need it use a `sync.Once` wrapper — copy the pattern in
`app/upgrades/v2_0_0/upgrade_test.go`. It must run before anything that builds address codecs.

**18 decimals, not 6.** `sdk.DefaultPowerReduction` is overridden to 10^18 in `app/app.go`'s
`init()`. Use `sdk.TokensFromConsensusPower` for amounts; raw `int64` literals overflow.

**Adding a module touches five places**, all under `app/`: the store key (`keepers/keys.go`),
keeper construction (`keepers/keepers.go`), `maccPerms` if it owns a module account
(`app.go`), the module manager, and **all three ordering lists** — BeginBlockers, EndBlockers,
and InitGenesis (which also sets ExportGenesis). Omitting an ordering list is a silent runtime
failure, not a compile error.

**A new store key on a live chain needs a store upgrade.** Adding to `keys.go` alone is not
enough; register `StoreUpgrades{Added: [...]}` on an `upgrades.Upgrade` or the node fails to
load state. `app/upgrades/v2_0_0/constants.go` is the worked example.

**Ante and post handler changes are consensus-breaking.** `app/ante.go` builds the decorator
chain by hand — it is wasmd's chain minus the circuit breaker, which is deliberate and
commented in place. Any change needs a coordinated upgrade. Be careful here: the existing
`app/ante_test.go` never exercises the real chain, so tests will not catch a regression.

### Upgrades

Each `app/upgrades/<version>/` exports an `upgrades.Upgrade` (name, handler creator, optional
store upgrades). Register new ones in the `Upgrades` slice in `app/app.go` —
`setupUpgradeStoreLoaders` and `setupUpgradeHandlers` pick them up from there automatically.

Naming is inconsistent by convention and must stay that way: directory `v2_0_0`, package
`v_2_0_0`, upgrade name `"v2.0.0"`. The upgrade name must match the on-chain plan name exactly.

### CosmWasm is permissionless

`CodeUploadAccess = AllowEverybody` and `InstantiateDefaultPermission = Everybody`, set in the
v2.0.0 upgrade handler and asserted by `test/e2e/wasm_test.go`. This is deliberate, and it's
also wasmd's own default. Tightening it is a governance parameter change, not a code change.
