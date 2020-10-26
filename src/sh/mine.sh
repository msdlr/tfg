#!/bin/sh

# Script to connect the node to a bootnode and start mining
# Params:
# $1: directory where the keypair is stored
# $2: public key of the account to use
# $3: bootnode to connect to (enode://...)
# $4: --console for interactive prompt

[ -z "$1" ] && (echo "no datadir";exit) || datadir=$1
[ -z "$2" ] && (echo "no account";exit) || account=$2
[ -z "$3" ] && (echo "no bootnode";exit) || bootnode=$3

geth \
 --datadir $datadir \
 --syncmode 'full' \
 --port 30311 \
 --rpc --rpcaddr '0.0.0.0' \
 --rpccorsdomain '*' --rpcport 8501 \
 --rpcapi 'personal,eth,net,web3,txpool,miner' \
 --bootnodes "$bootnode" \
 --networkid 1010 \
 --gasprice '0' \
 #--unlock "$account" \
 -dev -unlock 0 \
 #--password password.txt --mine
 #--allow-insecure-unlock
[ -z $4 ] $4
