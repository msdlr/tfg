#!/bin/bash

# Please execute this script from /path/to/repo/src/sh (its location)

# remove build/ directory if existing
rm -r build/ 2>/dev/null

for c in ../solidity/contracts/*sol
do

#echo "solc --abi $c --overwrite -o build"
solc --abi $c --overwrite -o build

# Generate the filename with .abi extension 
c_abi=$(echo $c | sed -e 's/\.sol/\.abi/ ' -e 's|../solidity/contracts/||g')
# Generate the filename with .go extension
c_go=$(echo $c | sed -e 's/\.sol/\.go/ ' -e 's|../solidity/contracts/||g')
# Generate the filename with .bin extension
c_bin=$(echo $c | sed -e 's/\.sol/\.bin/ ' -e 's|../solidity/contracts/||g')

# Convert .abi file to .go
#echo "abigen --abi=./build/$c_abi --pkg=main --out=$c_go"
abigen --abi=./build/$c_abi --pkg=main --out=solc_$c_go

#echo "solc --bin $c --overwrite -o build"
solc --bin $c --overwrite -o build

# Compile go contract with included deploy methods
#echo "abigen --bin=./build/$c_abi --abi=./build/$c_abi --pkg=main --out=$c_go"
abigen --bin=./build/$c_bin --abi=./build/$c_abi --pkg=main --out=solc_$c_go  

done

# Post: remove build/ directory and move .go files to src/go
rm -r build/ 2>/dev/null
mv -v *go ../golang
