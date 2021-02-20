#!/bin/bash
[ -z $1 ] && (echo "no node datadir";exit)
geth attach ipc:$1/geth.ipc
exit
