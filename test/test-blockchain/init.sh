#!/bin/bash

for node in ./n*
do
	geth --datadir "$node" init ./bc/test.json 
done
