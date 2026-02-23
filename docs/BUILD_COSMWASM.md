# Building hippo-protocol with CosmWasm Support

This document describes how to build hippo-protocol with CosmWasm (x/wasm) module support.

## Prerequisites

- Go 1.23 or later
- Make
- GCC (for CGO)
- Docker (for containerized builds)

## CosmWasm Version

The current integration uses CosmWasm wasmvm v2.2.4, which must match the version in `go.mod`.

## Standard Build

For development and testing on your local machine:

```bash
make build
```

This produces a dynamically-linked binary at `build/hippod`.

## Docker Build

For production deployments using Docker, the Dockerfile is configured to:

1. Use Alpine Linux for a minimal, secure base image
2. Download and verify libwasmvm static libraries (muslc) for both amd64 and arm64
3. Build with static linking to ensure the binary is self-contained
4. Verify the binary is statically linked

```bash
docker build -t hippo-protocol:latest .
```

### Docker Build Process

The Dockerfile performs these steps:

1. **Downloads libwasmvm_muslc libraries**: Pre-built static libraries for CosmWasm
2. **Verifies SHA256 checksums**: Ensures library integrity and authenticity
3. **Builds with muslc tag**: Enables static linking with musl libc
4. **Sets LINK_STATICALLY=true**: Forces static linking of all dependencies
5. **Verifies static linking**: Confirms the binary has no dynamic dependencies

## Release Build

For creating release binaries using goreleaser:

```bash
make ci-release
```

This uses the goreleaser configuration in `.goreleaser.yaml` which:

- Builds for multiple architectures (amd64, arm64, darwin)
- Downloads appropriate libwasmvm libraries for each platform
- Applies proper build tags and linker flags
- Creates versioned binary archives

## Build Options

### COSMOS_BUILD_OPTIONS

- `muslc`: Enable muslc build tag for static linking with musl libc (required for Alpine-based containers)
- `nostrip`: Disable symbol stripping (useful for debugging)

### LINK_STATICALLY

- Set to `true` to enable static linking (required for Alpine containers)

### LEDGER_ENABLED

- Set to `false` to disable Ledger hardware wallet support (reduces dependencies)
- Default: `true`

## Example: Static Build for Alpine

```bash
LEDGER_ENABLED=false COSMOS_BUILD_OPTIONS="muslc" LINK_STATICALLY=true make build
```

## Troubleshooting

### "cannot find libwasmvm_muslc" error

Ensure you're using the correct build environment:
- For Alpine/musl: Use Docker build or set `COSMOS_BUILD_OPTIONS="muslc"`
- For Debian/Ubuntu/glibc: Use standard build

### Dynamic linking in Docker

If the Docker build produces a dynamically-linked binary, check:
1. `LINK_STATICALLY=true` is set in the Dockerfile
2. `muslc` tag is included in COSMOS_BUILD_OPTIONS
3. libwasmvm_muslc libraries are downloaded to `/lib/`

## References

- [CosmWasm Integration Guide](https://github.com/CosmWasm/wasmd/blob/v0.54.2/INTEGRATION.md)
- [wasmvm Releases](https://github.com/CosmWasm/wasmvm/releases)
- [CosmWasm Documentation](https://cosmwasm.github.io/)
