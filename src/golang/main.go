package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
)

/*
	Global vars and functions
*/
var envHOME = os.Getenv("HOME")

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

// Load an eth keypair
func readKeystore(_pathToKeypair string, _password string) {
	// Encrypted keypair
	file := _pathToKeypair
	// Create a new keypair
	ks := keystore.NewKeyStore(envHOME+"/eth/", keystore.StandardScryptN, keystore.StandardScryptP)
	jsonBytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	password := _password
	account, err := ks.Import(jsonBytes, password, password)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex()) // 0x20F8D42FB0F667F2E53930fed426f225752453b3

	if err := os.Remove(file); err != nil {
		log.Fatal(err)
	}
}

// New keypair
func newKeystore(_path string, _pass string) {
	// Create the new kp file
	ks := keystore.NewKeyStore(_path, keystore.StandardScryptN, keystore.StandardScryptP)
	password := _pass
	// Encrypt private key
	account, err := ks.NewAccount(password)

	// Error checking
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(account.Address.Hex())
	fmt.Printf("Created keypair w/ pub address " + account.Address.Hex() + "\n")
}

// When a keypair is generated we don't know the exact filename for it
func locateGeneratedKeypair() {

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
		url = "http://localhost:8545"
	} else {
		url = os.Args[1]
	}

	// Stablish connection to the blockchain
	dial(url)

	// Create and setup the new address
	//newKeystore(envHOME+"eth", "hola")
	readKeystore(envHOME+"/eth/"+"UTC--2020-11-17T09-23-16.999473109Z--268c013964b50841fc534daa92954c2b049cb007", "hola")
}
