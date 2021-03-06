package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
	Global vars and functions
*/

/* GLOBAL VARS */
var ks *keystore.KeyStore      //Keystore for ethereum keys
var myTrOps *bind.TransactOpts // Transaction Ops
var myClient *ethclient.Client // Client retrieved after dialing the RPC endpoint

func main() {

	//region Check&set envvars

	// Ethereum KeyStore path (Unused when usin ganache as infrastructure)
	/*
		if os.Getenv("ETHKS") == "" {
			os.Setenv("ETHKS", os.Getenv("HOME")+"/eth/node1/keystore")
		}*/

	// Peer for connecting to the blockchain
	if os.Getenv("RPCENDPOINT") == "" {
		os.Setenv("RPCENDPOINT", "http://localhost:7545")
	}

	// Interval for checking blockchain events
	if os.Getenv("EVNTITV") == "" {
		os.Setenv("EVNTITV", "5")
	}

	if os.Getenv("CHAINID") == "" {
		os.Setenv("CHAINID", "5777")
	}

	// Ethereum private key
	if os.Getenv("PRIVKEY") == "" {
		os.Setenv("PRIVKEY", "ad92041b60126af952f8320b473ccb555d7274a53f1c27e12d2f1ea8aaecda7b")
	}
	// Ethereum public key
	if os.Getenv("PUBKEY") == "" {
		os.Setenv("PUBKEY", "FDb59BC058eFde421AdF049F27d3A03a4cedea2f")
	}

	// Ethereum public key
	if os.Getenv("USERNAME") == "" {
		os.Setenv("USERNAME", "msdlr")
	}

	// endregion

	/* Set up keystore with the correct path */
	//ks = keystore.NewKeyStore(os.Getenv("ETHKS"), keystore.StandardScryptN, keystore.StandardScryptP)

	// Retrieve TransactionOps and client object
	myTrOps, myClient := setupClient(os.Getenv("PRIVKEY"))
	deployAndInitialize(myTrOps, myClient,os.Getenv("PUBKEY"),os.Getenv("USERNAME"))
	myClient.Close()

	/*

		// Launch events checking
		go checkEvents()

		// Main routine stuck in inf loop
		for true {
			//fmt.Println("s")
		}

	*/
}


// setupClient retrieves the Transaction Ops and dials the RPC endpoint to establish the ethclient.client object
func setupClient(privKeyStr string) (tops *bind.TransactOpts, c *ethclient.Client) {
	/* Set-up client */
	privateKey, _ := crypto.HexToECDSA(privKeyStr)
	chainId, _ := strconv.Atoi(os.Getenv("CHAINID"))
	transactOps, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(chainId)))
	if err != nil {
		log.Fatalf("Error getting TransactionOps: %v", err)
	}

	client, err := ethclient.Dial(os.Getenv("RPCENDPOINT"))
	return transactOps, client
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

// string2Address parses a common.Address object from a string
func string2Address(pubKeyStr string) common.Address {
	return common.HexToAddress(pubKeyStr)
}

// getPubKeyFromPrivKey derives a public address corresponding to the private one passed as an argument
/*
func getPubKeyFromPrivKey(privKeyStr string) (pubk *common.Address) {
	// String -> private key object
	privateKey, err := crypto.HexToECDSA(privKeyStr)
	if err != nil {
		log.Fatal(err)
	}
	// Derive ECDSA public key
	publicKeyDerived := privateKey.Public()
	publicKeyECDSA, ok := publicKeyDerived.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	// Convert ECDSA public key to common.Address object (Ethereum)
	pubAddr := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &pubAddr
}
*/
