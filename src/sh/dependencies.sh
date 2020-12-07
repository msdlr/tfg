#!/bin/bash

# This scripts installs the dependencies (go and packages used) for the project

sudo apt install golang-go -y # The one in ubuntu/debian? repos

# Golang packages:
go get github.com/ethereum/go-ethereum # Ethereum
go get github.com/gorilla/websocket
go get github.com/shirou/gopsutil/cpu
