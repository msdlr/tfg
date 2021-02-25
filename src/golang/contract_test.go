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

// region function: addUser

func TestAddUserOk(t *testing.T){
	/* Arrange: We need an initialized contract */

	userAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5" // [1]
	ownerPrivStr := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"

	newuserPubStr := "0x12b3C6913a8eE35D1e0462f16Ac0aA6B6205a91a" // [2]
	newUserName := "testUser"

	to,c := initializeValidClient("http://localhost:7545",5777,ownerPrivStr)
	_, _, main, _, _, _ :=deployAndInitialize(to,c,userAddrStr,"TestOwner")

	/* Act: Call contract method AddUser */
	userAddr := publicAddressFromString(newuserPubStr)
	_, addUserErr := main.AddUser(to, userAddr, newUserName)

	/* Assert */
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

	if getUserIdStr != newUserName || getUserStr != newuserPubStr {
		t.Errorf("User ID and public key do not match")
	}

	// Check AuthContract

	userAuthContract,_ := main.GetUserAuthContract(nil, userAddr)

	userAuthStr := userAuthContract.String()
	if userAuthStr == common.BytesToAddress([]byte("0")).String() {
		t.Errorf("Error retrieving user's AuthContract")
	}
}

// Try to register a user with a username taken, and a username that is already registered
func TestAddUserNotOk(t *testing.T) {
	/* Arrange: We need an initialized contract with a user in it */

	userAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5" // [1]
	ownerPrivStr := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"

	userPubStr := "0x12b3C6913a8eE35D1e0462f16Ac0aA6B6205a91a" // [2]
	userNickname := "testUser"

	newuserPubStr := "0xa0A9e0409f8A0e03f41e1AAd5Bb19E86C4fE5Acc" // [3]


	to,c := initializeValidClient("http://localhost:7545",5777,ownerPrivStr)
	_, _, main, _, _, _ :=deployAndInitialize(to,c,userAddrStr,"TestOwner")

	userAddr := publicAddressFromString(userPubStr) // Call contract method AddUser
	_, _ = main.AddUser(to, userAddr, userNickname)


	/* Act: Try to add another different user with the same username */

	newUserAddr := publicAddressFromString(newuserPubStr)
	_, err1 := main.AddUser(to, newUserAddr, userNickname)

	/* Act: Try to add a user already on the system */

	_, err2 := main.AddUser(to, userAddr, userNickname+"1")

	/* Assert: is the new user registered with the same username as the existing one? */

	newUserRegistered, _ := main.GetUserRegistered(nil, newUserAddr)

	if newUserRegistered {
		t.Errorf("A second user was created with the same username")
	}

	/* Assert: check if a user is registered twice with to different usernames, it should be address 0 (username not taken) */

	shouldBeAddress0, _:= main.GetUserAddress(nil, userNickname+"1")

	if shouldBeAddress0.String() == userAddr.String() {
		t.Errorf("Username was registered twice")
	}

	if err1 == nil || err2 == nil {
		t.Errorf(err1.Error()+"\n"+err2.Error())
	}
}

//endregion

// region function: rmUser
func TestRemoveUserOk(t *testing.T ){
	/* Arrange: We need an initialized contract with a user in it */

	userAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5" // [1]
	ownerPrivStr := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"

	userPubStr := "0x12b3C6913a8eE35D1e0462f16Ac0aA6B6205a91a" // [2]
	userNickname := "testUser"

	to,c := initializeValidClient("http://localhost:7545",5777,ownerPrivStr)
	_, _, main, _, _, _ :=deployAndInitialize(to,c,userAddrStr,"TestOwner")

	userAddr := publicAddressFromString(userPubStr) // Call contract method AddUser
	_, _ = main.AddUser(to, userAddr, userNickname)

	registeredBeforRemoving, _ := main.GetUserRegistered(nil,userAddr)

	/* Act: Call rmUser */
	_, rmError := main.RmUser(to,userAddr,userNickname)

	/* Assert: check if user is not registered now */
	registeredAfterRemoving, _ := main.GetUserRegistered(nil,userAddr)
	isAdmin, _ := main.GetUserAdmin(nil,userAddr)
	loggedIn, _ := main.GetUserLoggedIn(nil,userAddr)
	addressForUserName, _ := main.GetUserAddress(nil,userNickname) // Address for the username after removing the user

	if registeredBeforRemoving == registeredAfterRemoving ||
		rmError != nil ||
		isAdmin == true ||
		loggedIn == true ||
		addressForUserName.String() != publicAddressFromString("0x0").String()	{
		t.Errorf("User was not properly removed")
	}
}

func TestRemoveOwner(t *testing.T ){

}

func TestRemoveMismatching(t *testing.T ){

}
// endregion

// region function: promoteUser
func TestPromoteUserOk(t *testing.T) {

}

func TestPromoteUserFailed(t *testing.T) {

}

// endregion

// region func demoteAdmin
func TestDemoteAdminOk(t *testing.T) {

}

func TestDemoteAdminFailed(t *testing.T) {

}
// endregion