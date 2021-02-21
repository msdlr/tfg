package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
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
var ks *keystore.KeyStore		//Keystore
var myTrOps *bind.TransactOpts 	// Transaction Ops
var myClient *ethclient.Client	// Client retrieved after dialing the RPC endpoint

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


	/* Set up keystore with the correct path */
	ks = keystore.NewKeyStore(os.Getenv("ETHKS"), keystore.StandardScryptN, keystore.StandardScryptP)

	// Retrieve TransactionOps and client object
	myTrOps, myClient := setupClient()

	// DeployMain(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Main, error)
	addr, trans, main, err := DeployMain(myTrOps, myClient)
	main = main

	if err == nil {
		fmt.Println("## NEW CONTRACT DEPLOYED ##")
		fmt.Println("Address: ", addr.Hex())
		fmt.Println("Transaction hash: ", trans.Hash())
		fmt.Println("Gas Used: ", trans.Gas(), "(price:", trans.GasPrice(), ")")
		fmt.Println("Nonce: ", trans.Nonce())
	}

	myPubKey := getPubKeyFromPrivKey(os.Getenv("PRIVKEY"))

	main.Initialize(myTrOps, *myPubKey, "msdlr")

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
func setupClient() (*bind.TransactOpts, *ethclient.Client) {
	/* Set-up client */
	privateKey, _ := crypto.HexToECDSA(os.Getenv("PRIVKEY"))
	chainId, _ := strconv.Atoi(os.Getenv("CHAINID"))
	transactOps, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(chainId)))
	if err != nil {
		log.Fatalf("Error getting TransactionOps: %v", err)
	}

	client, err := ethclient.Dial(os.Getenv("RPCENDPOINT"))
	return transactOps, client
}

// getPubKeyFromPrivKey derives a public address corresponding to the private one passed as an argument
func getPubKeyFromPrivKey(privKeyStr string)(pubk *common.Address ){
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
