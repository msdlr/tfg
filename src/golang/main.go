package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"os"
	"strconv"
	"time"
)

/*
	Global vars and functions
*/

/* GLOBAL VARS */
// keystore to use ethereum keys
var ks *keystore.KeyStore

func main() {
	/*
		Check&set envvars
	*/

	var envHOME = os.Getenv("HOME") // This is not to be modified
	envHOME=envHOME
	// ETHereum KeyStore path
	if os.Getenv("ETHKS") == "" {
		os.Setenv("ETHKS", os.Getenv("HOME")+"/eth/node1/keystore")
	}

	// Peer for connecting to the blockchain
	if os.Getenv("RPCENDPOINT") == "" {
		os.Setenv("RPCENDPOINT","http://localhost:7545")
	}

	// Interval for checking blobkchain events
	if os.Getenv("EVNTITV") == "" {
		os.Setenv("EVNTITV","5")
	}


	// Set up keystore with the correct path
	ks = keystore.NewKeyStore(os.Getenv("ETHKS"), keystore.StandardScryptN, keystore.StandardScryptP)


	// Launch events checking
	go checkEvents()

	loadTestAccount("0xFDb59BC058eFde421AdF049F27d3A03a4cedea2f", "ad92041b60126af952f8320b473ccb555d7274a53f1c27e12d2f1ea8aaecda7b")
	// Main routine stuck in inf loop
	for true {
		//fmt.Println("s")
	}

}

// Start a goroutine to check for events in the Blockchain
func checkEvents() {
	secs, _ := strconv.Atoi(os.Getenv("EVNTITV")) // Get seconds as number

	interval := time.Duration(secs) * time.Second
	interval=interval
	fmt.Println("Checking events...")

	ticker := time.NewTicker(interval)

	for range ticker.C {
		// Code here is executed every second
		fmt.Println("Checking events...", time.Now().UTC())
	}

}
