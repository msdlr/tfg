package main

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// Test keystore
var testKs keystore.KeyStore
var ksDir = "/tmp/aaaa"

func TestNewAccountOK(t *testing.T) {
	testKs := keystore.NewKeyStore(ksDir, keystore.StandardScryptN, keystore.StandardScryptP)

	newPubkey := CreateNewAccount(testKs, "passw0rd")
	if newPubkey == "" {
		t.Errorf("[Error] TestNewAccountOK:\tAccount not created")
	}
}

func TestNewAccountNilKS(t *testing.T) {
	// Give a location for the KS where user has no permission
	testKs := keystore.NewKeyStore("/", keystore.StandardScryptN, keystore.StandardScryptP)
	var newPubkey string = ""

	newPubkey = CreateNewAccount(testKs, "passw0rd")
	if newPubkey != "" {
		t.Errorf("[Error] TestNewAccountNilKS:\tAccount was created somehow (?)")
	}
}

func TestOpenAccount(t *testing.T) {
	testKs := keystore.NewKeyStore(ksDir, keystore.StandardScryptN, keystore.StandardScryptP)
	// We need to create new wallets to have public keys to unlock
	successfulPub := CreateNewAccount(testKs, "passw0rd")
	failPub := CreateNewAccount(testKs, "passw0rd")

	// Close them
	CloseAccount(testKs, successfulPub)
	CloseAccount(testKs, failPub)

	// Once we have the public key we test the method
	shouldBe1 := OpenAccount(testKs, successfulPub, "passw0rd") // Successful
	shouldBe2 := OpenAccount(testKs, failPub, "passw0rd!")      // Wrong password
	shouldBe3 := OpenAccount(testKs, "newPubkey", "passw0rd")   // Pub key not found

	if shouldBe1 != 1 || shouldBe2 != 2 || shouldBe3 != 3 {
		t.Errorf("[Error] TestOpenAccount:\tExpected 1 (was %d)\nExpected 2 (was %d)\nExpected 3 (was %d)\n", shouldBe1, shouldBe2, shouldBe3)
	}
}

func TestCloseAccountOK(t *testing.T) {
	testKs := keystore.NewKeyStore(ksDir, keystore.StandardScryptN, keystore.StandardScryptP)
	// We need to create new wallets to have public keys to unlock
	pub := CreateNewAccount(testKs, "passw0rd")
	// Accounts unlocked by default
	//OpenAccount(testKs, pub, "passw0rd")

	if !CloseAccount(testKs, pub) {
		t.Error("[Error] TestCloseAccountOK:\tAccount not unlocked")
	}
}

// endregion
