package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"os"
	"strings"
	"testing"
)

// region Auxiliary/stub functions

// Auxiliary function for client stub
func initializeValidClient(endpoint string,chainId uint16,privkey string)(tops *bind.TransactOpts, c *ethclient.Client){
	os.Setenv("RPCENDPOINT", endpoint)
	chainIdStr := fmt.Sprintf("%d",chainId)
	os.Setenv("CHAINID", chainIdStr)

	// [1] in Ganache
	tops, c = setupClient(privkey)

	return
}

// endregion

func TestDeployOk(t *testing.T){
	// Arrange
	to,c := initializeValidClient("http://localhost:7545",5777,"b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4")

	// Act
	// (addr common.Address, deployTrans *types.Transaction, main *Main, deployError error, initTrans *types.Transaction, initError error)
	addr, deployTrans, main, deployError, initTrans, initError :=deployAndInitialize(to,c,"0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5","TestOwner")
	// Address is not 0x00...0

	addrStr := addr.String()
	zeroStr := publicAddressFromString("0x0").Hex()
	zeroAddress := strings.Compare(addrStr,zeroStr) == 0
	if  zeroAddress ||
		deployTrans == nil ||
		main == nil ||
		deployError != nil ||
		initTrans == nil ||
		initError != nil{
			t.Errorf("Failed to instantiate Contract")
	}

	// Assert

}
