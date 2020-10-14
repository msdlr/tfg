#!/bin/sh

command -v entr >/dev/null || (echo "install entr"; exit)

while true 
do
	ls contracts/*sol test/*sol | entr truffle test
done
