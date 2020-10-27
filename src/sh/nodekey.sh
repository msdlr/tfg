#!/bin/sh

[ -z "$1" ] && (echo "Select directory"; exit)
touch $1/bootnode.key
bootnode -genkey $1/bootnode.key
bootnode -nodekey $1/bootnode.key
exit
