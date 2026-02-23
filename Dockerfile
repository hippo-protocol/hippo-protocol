FROM golang:1.23-alpine3.19 AS build-env

# Install minimum necessary dependencies
# See https://github.com/CosmWasm/wasmvm/releases for wasmvm compatibility
RUN set -eux; apk add --no-cache ca-certificates build-base git

# Add 'hippo-protocol' source files
COPY . /src/hippo-protocol

# Set working directory for the 'hippo-protocol' build
WORKDIR /src/hippo-protocol

# Download and verify libwasmvm static libraries
# See https://github.com/CosmWasm/wasmvm/releases/tag/v2.2.4
ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.2.4/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v2.2.4/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 27fb13821dbc519119f4f98c30a42cb32429b111b0fdc883686c34a41777488f
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 70c989684d2b48ca17bbd55bb694bbb136d75c393c067ef3bdbca31d2b23b578

# Build with static linking for CosmWasm support
RUN LEDGER_ENABLED=false COSMOS_BUILD_OPTIONS="muslc" LINK_STATICALLY=true make build

# Verify binary is statically linked
RUN echo "Ensuring binary is statically linked ..." \
  && (file /src/hippo-protocol/build/hippod | grep "statically linked")

# Final image
FROM alpine:3.19

# Copy over binaries from the build-env
COPY --from=build-env /src/hippo-protocol/build/hippod /usr/bin/hippod

RUN chmod +x /usr/bin/hippod

EXPOSE 26656 26657 1317 9090

CMD ["/usr/bin/hippod", "version"]
