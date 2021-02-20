#!/bin/sh
sudo apt-get install software-properties-common & 
sudo add-apt-repository -y ppa:ethereum/ethereum &&
sudo add-apt-repository -y ppa:ethereum/ethereum-dev
sudo apt update
sudo apt install ethereum solc -y
exit
