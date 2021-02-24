package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
	os.Setenv("PRIVKEY", privkey)

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

	// Assert
	addrStr := addr.String()
	zeroStr := publicAddressFromString("0x0").Hex() // Address is not 0x00...0
	zeroAddress := strings.Compare(addrStr,zeroStr) == 0
	if  zeroAddress ||
		deployTrans == nil ||
		main == nil ||
		deployError != nil ||
		initTrans == nil ||
		initError != nil{
			t.Errorf("Failed to instantiate Contract")
	}

	// Check owner
	addrGet,err :=  main.GetContractAddress(nil)
	addrGetStr := addrGet.String()

	if addrGetStr != addrStr {
		t.Errorf("Created contract address differs from getter")
	}

	if err != nil {
		fmt.Println("Error:"+err.Error())
	}
}

func TestAddUserOk(t *testing.T){
	// Arrange: We need an initialized contract

	ownerPub := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5" // [1]
	ownerPriv := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"

	newUserPub := "0x12b3C6913a8eE35D1e0462f16Ac0aA6B6205a91a" // [2]
	newUserName := "testUser"

	to,c := initializeValidClient("http://localhost:7545",5777,ownerPriv)
	_, _, main, _, _, _ :=deployAndInitialize(to,c,ownerPub,"TestOwner")

	// Act: Call contract method AddUser
	userAddr := publicAddressFromString(newUserPub)
	_, addUserErr := main.AddUser(to, userAddr, newUserName)

	// Assert
	if addUserErr != nil {
		t.Errorf("Could not add user."+ addUserErr.Error())
	}

	// Check if registered via getter
	isUserRegistered, userRegisteredErr := main.GetUserRegistered(nil, userAddr)

	if userRegisteredErr != nil {
		t.Errorf("Error retrieving if user is registered "+ userRegisteredErr.Error())
	}

	if !isUserRegistered {
		t.Errorf("User is not registered (false, expected true)")
	}

	// Check if admin via getter
	isUserAdmin, userAdminErr := main.GetUserAdmin(nil, userAddr)

	if userAdminErr != nil {
		t.Errorf("Error retrieving if user is registered "+ userAdminErr.Error())
	}

	if isUserAdmin {
		t.Errorf("User is admin (true, expected false)")
	}

	// Check user attempts via getter
	attempts, attemptsErr := main.GetUserAttempts(nil, userAddr)

	if attemptsErr != nil {
		t.Errorf("Error getting user attempts"+ userAdminErr.Error())
	}

	if attempts > 0 {
		t.Errorf("Error, attempts="+string(attempts) +", expected 0"+ userAdminErr.Error())
	}

	// Check getters for relating pubkey and username
	getUserIdStr, getUserIdErr :=main.GetUserId(nil, userAddr)
	getUserAddress, getUserAddressErr := main.GetUserAddress(nil,newUserName)

	if getUserIdErr != nil {
		t.Errorf("Error relating user public address to username "+ getUserIdErr.Error())
	}

	if getUserAddressErr != nil {
		t.Errorf("Error relating username to public address"+ getUserAddressErr.Error())
	}

	getUserStr:= getUserAddress.String()

	if getUserIdStr != newUserName || getUserStr != newUserPub {
		t.Errorf("User ID and public key do not match")
	}

	// Check AuthContract

	userAuthContract,_ := main.GetUserAuthContract(nil, userAddr)

	userAuthStr := userAuthContract.String()
	if userAuthStr == common.BytesToAddress([]byte("0")).String() {
		t.Errorf("Error retrieving user's AuthContract")
	}
}
