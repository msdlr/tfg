#!/bin/sh

# Script to create an Ethereum wallet
# $1 : location directory for the keypair

[ -z "$1" ] && (echo "no datadir" ; exit)

geth -datadir "$1" -dev account new
exit
