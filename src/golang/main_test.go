package main

import (
	"os"
	"testing"
)

func TestSetupClientOk(t *testing.T) {
	os.Setenv("RPCENDPOINT","http://localhost:7545")
	os.Setenv("CHAINID", "5777")
	privkey := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"
	pubkey := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5"

	myTrOps, myClient := setupClient(privkey)

	// Check TransactOps (myTrOps)
	FromAddrStr := myTrOps.From.Hex()
	if FromAddrStr != pubkey {
		t.Errorf("[Error] testSetupClientNotOk:\tNull account (expected not null)")
	}

	// Check Client object
	if myClient == nil {
		t.Errorf("[Error] testSetupClientNotOk:\tnil client (expected not nil)")
	}
}

// Wrong Endpoint IP/port/protocol
func TestSetupClientWrongEndpoint(t *testing.T) {
	os.Setenv("RPCENDPOINT","0.0.0.0:0")
	os.Setenv("CHAINID", "5777")
	privkey := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"
	pubkey := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5"

	myTrOps, myClient := setupClient(privkey)

	// Check TransactOps (myTrOps)
	FromAddrStr := myTrOps.From.Hex()
	if FromAddrStr != pubkey {
		t.Errorf("[Error] TestSetupClientWrongEndpoint:\tNull account (expected not null)")
	}

	// Check Client object
	if myClient != nil {
		t.Errorf("[Error] TestSetupClientWrongEndpoint:\tnot nil client (expected nil)")
	}
}

// Wrong Private key
func TestSetupClientWrongTransOps(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("[Error] TestSetupClientWrongTransOps:\t expected panic")
		}
	}()

	os.Setenv("RPCENDPOINT","0.0.0.0:0")
	privkey := "private key"

	// It should get a panic
	_, _ = setupClient(privkey)
}