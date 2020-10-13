#!/bin/sh

command -v entr || (echo "install entr"; exit)

ls ./*sol | entr truffle test