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

	//region Check&set envvars

	// Ethereum KeyStore path (Unused when usin ganache as infrastructure)
	/*
	if os.Getenv("ETHKS") == "" {
		os.Setenv("ETHKS", os.Getenv("HOME")+"/eth/node1/keystore")
	}*/

	// Peer for connecting to the blockchain
	if os.Getenv("RPCENDPOINT") == "" {
		os.Setenv("RPCENDPOINT","http://localhost:7545")
	}

	// Interval for checking blockchain events
	if os.Getenv("EVNTITV") == "" {
		os.Setenv("EVNTITV","5")
	}

	if os.Getenv("CHAINID") == "" {
		os.Setenv("CHAINID","5777")
	}

	// Ethereum private key
	if os.Getenv("PRIVKEY") == "" {
		os.Setenv("PRIVKEY","ad92041b60126af952f8320b473ccb555d7274a53f1c27e12d2f1ea8aaecda7b")
	}
	// Ethereum public key
	if os.Getenv("PUBKEY") == "" {
		os.Setenv("PUBKEY","FDb59BC058eFde421AdF049F27d3A03a4cedea2f")
	}

	// endregion

	/* Set up keystore with the correct path */
	//ks = keystore.NewKeyStore(os.Getenv("ETHKS"), keystore.StandardScryptN, keystore.StandardScryptP)

	// Retrieve TransactionOps and client object
	myTrOps, myClient := setupClient(os.Getenv("PRIVKEY"))

	// DeployMain(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Main, error)
	addr, trans, main, err := DeployMain(myTrOps, myClient)

	if err == nil {
		fmt.Println("## NEW CONTRACT DEPLOYED ##")
		fmt.Println("Address:\t\t", addr.Hex())
		fmt.Println("Transaction hash:\t", trans.Hash())
		fmt.Println("Gas Used:\t\t", trans.Gas(), "(price:", trans.GasPrice(), ")")
		fmt.Println("Nonce:\t\t\t", trans.Nonce())
	}

	myPubKey := publicAddressFromString(os.Getenv("PUBKEY"))

	main.Initialize(myTrOps, myPubKey, "msdlr")

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
func setupClient(privKeyStr string) (*bind.TransactOpts, *ethclient.Client) {
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

// publicAddressFromString generates a common.Address object from a string which represents a public key
func publicAddressFromString(pubKeyStr string) common.Address {
	return common.HexToAddress(pubKeyStr)
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

// region Keystore functions


// openAccount loads an eth keypair with a password, from the keystore path
func OpenAccount(ks *keystore.KeyStore, _pubKey string, _password string) uint8 {

	// This returns an array with the keys stored in the keystore path
	var ethAccArray = ks.Accounts()

	// Iterate through every account to find which pub key coincides
	for i := 0; i < len(ethAccArray); i++ {
		// Check if the pub key is the same as the one provided
		if _pubKey == ethAccArray[i].Address.Hex() {
			err := ks.Unlock(ethAccArray[i], _password)
			if err == nil {
				os.Setenv("ETHACC", ethAccArray[i].Address.Hex())
				// Account is unlocked, we can get out of the iterating loop
				fmt.Println("Account " + _pubKey + " unlocked!")
				return 1
			} else {
				// There has been an error unlocking the account
				fmt.Println("Account was found but not unlocked")
				return 2
			}
		}
	}
	// If no account was unlocked
	return 3
}

// CreateNewAccount: New keypair
func CreateNewAccount(ks *keystore.KeyStore, _pass string) string {
	// Encrypt private key
	account, err := ks.NewAccount(_pass)

	// Error checking
	if err != nil {
		//log.Fatal(err)
		return ""
	}
	// fmt.Println(account.Address.Hex())
	fmt.Printf("Created keypair w/ pub address " + account.Address.Hex() + "\n")
	return account.Address.Hex()
}

// CloseAccount closes the ethereum wallet when we are finished
func CloseAccount(ks *keystore.KeyStore, _pub string) bool {
	// This returns an array with the keys stored in the keystore path
	var acc string = os.Getenv("ETHACC")

	// Iterate through every account to find which pub key coincides
	for _, account := range ks.Accounts() {
		// Search for the public key to lock
		if _pub == account.Address.Hex() {
			// If found, it is locked so that noone can access it.
			err := ks.Lock(account.Address)
			if err == nil {
				os.Setenv("ETHACC", account.Address.Hex())
				// Account is unlocked, we can get out of the iterating loop
				fmt.Println("Account " + acc + " locked")
				return true
			} else {
				fmt.Println("Account was found but not locked")
			}
		}
	}
	return false
}
// endregion
