#!/usr/bin/env bash

#== Requirements ==
#
## make sure your `go env GOPATH` is in the `$PATH`
## Install:
## + latest buf (v1.0.0-rc11 or later)
## + protobuf v3
#
## All protoc dependencies must be installed not in the module scope
## currently we must use grpc-gateway v1
# cd ~
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
# go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0
# go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
# go install github.com/cosmos/gogoproto/protoc-gen-gocosmos@latest
# go get github.com/regen-network/cosmos-proto@latest # doesn't work in install mode


set -eo pipefail

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find ./hippo -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    buf generate --template buf.gen.gogo.yaml $file
  done
done

cd ..

# move proto files to the right places
cp -r github.com/hippocrat-dao/hippo-protocol/* ./
rm -rf github.com