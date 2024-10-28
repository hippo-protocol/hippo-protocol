FROM golang:1.22.0-bullseye AS build-env

# Install minimum necessary dependencies,
ENV PACKAGES make git gcc
RUN apt-get update -y
RUN apt-get install -y $PACKAGES

# Add 'hippo-protocol' source files
COPY . /src/hippo-protocol

# Set working directory for the 'hippo-protocol' build
WORKDIR /src/hippo-protocol

# Install hippo-protocol
RUN make clean && make build

# Final image
FROM debian:bullseye-slim

# Copy over binaries from the build-env
COPY --from=build-env /src/hippo-protocol/build/hippod /usr/bin/hippod

RUN chmod +x /usr/bin/hippod

EXPOSE 26656 26657 1317 9090
