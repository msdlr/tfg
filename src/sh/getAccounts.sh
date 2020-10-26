#!/bin/sh

# Script to print ethereum addresses in a node
# $1 : path to keypairs

[ -z "$1" ] && (echo "no node datadir";exit)

geth -datadir "$dir" -dev account list 2>/dev/null | awk ' {/^Account/; gsub(/[{}]/,""); gsub(/#/,""); print $3 }' | sed -e 's/://'
exit
