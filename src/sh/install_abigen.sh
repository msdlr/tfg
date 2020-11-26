#!/bin/sh
# Install: solc & protoc
sudo apt install solc protobuf-compiler 
# Install: abigen
go get -u github.com/ethereum/go-ethereum &&
cd $GOPATH/src/github.com/ethereum/go-ethereum/ &&
make &&
make devtools
