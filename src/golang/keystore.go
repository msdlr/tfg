package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

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
