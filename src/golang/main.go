package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
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

	// ETHereum KeyStore path
	if os.Getenv("ETHKS") == "" {
		os.Setenv("ETHKS", os.Getenv("HOME")+"/eth/node1/keystore")
	}

	// Peer for connecting to the blockchain
	if os.Getenv("RPCENDPOINT") == "" {
		os.Setenv("RPCENDPOINT","http://localhost:7545")
	}

	// Interval for checking blockchain events
	if os.Getenv("EVNTITV") == "" {
		os.Setenv("EVNTITV","5")
	}

	// Chain ID
	if os.Getenv("CHAINID") == "" {
		os.Setenv("CHAINID","1337")
	}

	// Ethereum private key
	if os.Getenv("PRIVKEY") == "" {
		os.Setenv("PRIVKEY","ad92041b60126af952f8320b473ccb555d7274a53f1c27e12d2f1ea8aaecda7b")
	}

	// Set up keystore with the correct path
	ks = keystore.NewKeyStore(os.Getenv("ETHKS"), keystore.StandardScryptN, keystore.StandardScryptP)

	privateKey, _ := crypto.HexToECDSA( os.Getenv("PRIVKEY") )
	chainId,_ :=strconv.Atoi(os.Getenv("CHAINID"))

	transactOps, err := bind.NewKeyedTransactorWithChainID(privateKey,big.NewInt(int64(chainId)))
	if err != nil {
		log.Fatalf("Error getting TransactionOps: %v", err)
	}

	client, err := ethclient.Dial("HTTP://127.0.0.1:7545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	client = client

	// DeployMain(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Main, error)
	addr, trans,main,err :=DeployMain(transactOps,client)

	fmt.Println(addr.Hex(),trans,main)


	/*

	// Launch events checking
	go checkEvents()

	// Main routine stuck in inf loop
	for true {
		//fmt.Println("s")
	}

	*/
}

// Start a goroutine to check for events in the Blockchain
func checkEvents() {
	secs, _ := strconv.Atoi(os.Getenv("EVNTITV")) // Get seconds as number

	interval := time.Duration(secs) * time.Second
	fmt.Println("Checking events...")

	ticker := time.NewTicker(interval)

	for range ticker.C {
		// Code here is executed every second
		fmt.Println("Checking events...", time.Now().UTC())
	}

}
