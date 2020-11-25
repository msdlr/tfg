package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

/*
	Global vars and functions
*/
/* ENVVARS */
// $HOME
var envHOME = os.Getenv("HOME")
// New envvar: ETHKS, where keystores are saved
var _ = os.Setenv("ETHKS", envHOME+"/eth/")
var envETHKS string = os.Getenv("ETHKS")
var _ = os.Setenv("RPCENDPOINTG","http://localhost:7545")
var envRPCENDPOINT = "http://localhost:7545"

// Pub key of the unlocked string if needed. env: ETHACC
var envUnlockedAccount string

/* GLOBAL VARS */
// keystore to use ethereum keys
var ks = keystore.NewKeyStore(envETHKS, keystore.StandardScryptN, keystore.StandardScryptP)

func main() {
	/*
		Param checking:
		If none: default (localhost)
		1: url of node to stablish connection to.
	*/

	// No parameters provided
	var url string
	if len(os.Args) == 1 {
		url = envRPCENDPOINT
	} else {
		url = os.Args[1]
	}

	// Stablish connection to the blockchain
	// contactBlockchain returns an error, nil if none
	if contactBlockchain(url) != nil {
		fmt.Println("Error connecting to the blockchain")
		os.Exit(1)
	}

	// Create and setup the new address
	//newAccount(envHOME+"eth", "prueba")
	useAccount("0x0DDB3d979973A0288F4832676d2e6Aa29bC1d42d", "prueba")
}
