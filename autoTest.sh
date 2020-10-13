#!/bin/sh

command -v entr || (echo "install entr"; exit)

ls contracts/*sol test/*sol | entr truffle test