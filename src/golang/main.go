package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

func dial(url string) {
	// Dial address: ganache in localhost
	conn, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal("Error reaching RCP.", err)
		fmt.Println("Connection to provided URL failed")
	}


	fmt.Println("Connection to "+ url+" successful!")
	_ = conn // we'll use this in the upcoming sections
}

func main(){
	/*
	Param checking:
	If none: default (localhost)
	1: url of node to stablich connection to.
	 */

	// No parameters provided
	var url string
	if len(os.Args) == 1 {
		url = "http://localhost:8545"
	} else {
		url = os.Args[1]
	}

	// Stablish connection to the blockchain
	dial(url)
}
