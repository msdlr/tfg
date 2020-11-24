package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
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

/* GLOBAL VARS */
// keystore to search for ethereum keys
var ks = keystore.NewKeyStore(envETHKS, keystore.StandardScryptN, keystore.StandardScryptP)

func dial(url string) {
	// Dial address: ganache in localhost
	conn, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal("Error reaching RCP.", err)
		fmt.Println("Connection to provided URL failed")
	} else {
		fmt.Println("Connection to " + url + " successful!")
		_ = conn
	}
}

// Load an eth keypair with a password, from the keystore path
func useAccount(_pubKey string, _password string) {
	// Encrypted keypair
	//file := _pubKey
	// Create a new keypair

	// This returns an array with the keys stored in the keystore path
	var ethAccArray = ks.Accounts()
	var success bool = false

	// Iterate through every account to find which pub key coincides
	for i:= 0; i< len(ethAccArray); i++ {
		// Check if the pub key is the same as the one provided
		if _pubKey == ethAccArray[i].Address.Hex()  {
			err := ks.Unlock(ethAccArray[i], _password)
			if err == nil {
				success = true
				// Account is unlocked, we can get out of the iterating loop
				fmt.Println("Account " + _pubKey + " unlocked!")
				break
			} else {
				// There has been an error unlocking the account
				success = false
				fmt.Println("Account was found but not unlocked")
			}
		}
	}
	// If no account was unlocked
	if !success {
		fmt.Println("Account was not found")
	}
}

// New keypair
func newAccount(_pass string) string {
	// Encrypt private key
	account, err := ks.NewAccount(_pass)

	// Error checking
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(account.Address.Hex())
	fmt.Printf("Created keypair w/ pub address " + account.Address.Hex() + "\n")
	return account.Address.Hex()
}

func main() {
	/*
		Param checking:
		If none: default (localhost)
		1: url of node to stablich connection to.
	*/

	// No parameters provided
	var url string
	if len(os.Args) == 1 {
		url = "http://localhost:7545"
	} else {
		url = os.Args[1]
	}

	// Stablish connection to the blockchain
	dial(url)

	// Create and setup the new address
	//newAccount(envHOME+"eth", "prueba")
	useAccount("0x0DDB3d979973A0288F4832676d2e6Aa29bC1d42d", "prueba")
}
