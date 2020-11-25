package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"os"
)

func contactBlockchain(url string) error {
	// Dial address: ganache in localhost
	conn, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal("Error reaching RCP.", err)
		fmt.Println("Connection to provided URL failed")
	} else {
		fmt.Println("Connection to " + url + " successful!")
		_ = conn
	}
	return err
}

// Load an eth keypair with a password, from the keystore path
func openAccount(_pubKey string, _password string) {

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
				os.Setenv("ETHACC",ethAccArray[i].Address.Hex())
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

// Close the ethereum wallet when we are finished
func closeAccount(){
	// This returns an array with the keys stored in the keystore path
	var ethAccArray = ks.Accounts()
	var acc string = os.Getenv("ETHACC")

	// Iterate through every account to find which pub key coincides
	for i:= 0; i< len(ethAccArray); i++ {
		// Search for the public key to lock
		if acc == ethAccArray[i].Address.Hex() {
			// If found, it is locked so that noone can access it.
			err := ks.Lock(ks.Accounts()[i].Address)
			if err == nil {
				os.Setenv("ETHACC", ethAccArray[i].Address.Hex())
				// Account is unlocked, we can get out of the iterating loop
				fmt.Println("Account " + acc + " locked")
				break
			} else {
				fmt.Println("Account was found but not locked")
			}
		}
	}
}