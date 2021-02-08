#!/bin/sh

# Script to create an Ethereum wallet
# $1 : location directory for the keypair

[ -z "$1" ] && geth -dev account new && exit

geth -datadir "$1" -dev account new
exit
