package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"os"
	"strings"
	"testing"
)

var (
	testOwnerAddrStr = "0xFDb59BC058eFde421AdF049F27d3A03a4cedea2f"
	testOwnerPrivKey     = "ad92041b60126af952f8320b473ccb555d7274a53f1c27e12d2f1ea8aaecda7b"
	testOwnerUsername    = "0wn3r"

	testAdmin2AddrStr  = "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5"
	testAdmin2PrivKey  = "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"
	testAdmin2Username = "4dmin"

	testUser1AddrStr  = "0x12b3C6913a8eE35D1e0462f16Ac0aA6B6205a91a"
	testUser1PrivKey  = "e7c911fedc61cc1fd1a7a1cb84fd449562709cfa16a39f228cca07158c7307fb"
	testUser1Username = "User 0ne"

	testUser2AddrStr  = "0xD03A8E7E2265CD8239F34909324F98F00496EA31"
	testUser2PrivKey  = "7aa5be5263617d40346f8dc8d32f59a6cc6443bbf8d164bc1b89170f2d0679af"
	testUser2Username = "User Tw0"

	testChainId     uint16 = 5777
	testRpcEndpoint        = "http://localhost:7545"
)

// region Auxiliary/stub functions

// Auxiliary function for client stub
func initializeValidClient(endpoint string,chainId uint16,privkey string)(tops *bind.TransactOpts, c *ethclient.Client){
	os.Setenv("RPCENDPOINT", endpoint)
	chainIdStr := fmt.Sprintf("%d",chainId)
	os.Setenv("CHAINID", chainIdStr)
	os.Setenv("PRIVKEY", privkey)
	
	tops, c = setupClient(privkey)

	return
}

func testEvents(contractAddress common.Address, c *ethclient.Client) {
	query := ethereum.FilterQuery{ // Ethereum query for events in a given address
		Addresses: []common.Address{contractAddress},
	}

	logChannel := make(chan types.Log) // Channel to receive messages from
	// Subscribe for events
	sub, subErr := c.SubscribeFilterLogs(context.Background(), query, logChannel)
	if !strings.Contains(subErr.Error(), "notifications not supported") {
		// Testing with Ganache (HTTP) does not support subscribing
		for {
			select {
			case err := <-sub.Err():
				log.Fatal(err)
			case vLog := <-logChannel:
				fmt.Println(vLog) // pointer to event log
			}
		}
	} else {
		// Print about notifications not being supported
		fmt.Println("(Notifications for events not supported)")
	}
}

// endregion



func TestDeployOk(t *testing.T){
	/* Arrange: create valid client */
	to,c := initializeValidClient(testRpcEndpoint, testChainId,testOwnerPrivKey)

	/* Act: deploy and initialize general contract */
	addr, deployTrans, contract, deployError, initTrans, initError :=deployAndInitialize(to,c, testOwnerAddrStr,testOwnerUsername)

	/* Assert contract created properly */
	addrStr := addr.String()
	zeroStr := string2Address("0x0").Hex()
	zeroAddress := strings.Compare(addrStr,zeroStr) == 0
	if  zeroAddress ||
		deployTrans == nil ||
		contract == nil ||
		deployError != nil ||
		initTrans == nil ||
		initError != nil{
			t.Errorf("Failed to instantiate Contract")
	}

	// Check owner
	addrGet,err :=  contract.GetContractAddress(nil)
	addrGetStr := addrGet.String()

	ownerAddr, _ := contract.GetOwner(nil)

	isTheOwnerAdmin,_ :=contract.GetUserAdmin(nil,string2Address(testOwnerAddrStr))

	if !isTheOwnerAdmin {
		t.Errorf("Contract owner was not made admin")
	}

	if ownerAddr.Hex() != testOwnerAddrStr {
		t.Errorf("Contract owner not properly set")
	}

	if addrGetStr != addrStr || err != nil{
		t.Errorf("Created contract address differs from getter "+err.Error())
	}
}


func TestAddUserOk(t *testing.T){
	/* Arrange: We need an initialized contract */

	to,c := initializeValidClient(testRpcEndpoint, testChainId,testOwnerPrivKey)
	_, _, contract, _, _, _ :=deployAndInitialize(to,c, testOwnerAddrStr,testOwnerUsername)

	registeredBeforeAdd, _ := contract.GetUserRegistered(nil, string2Address(testUser1AddrStr))

	userAddr := string2Address(testUser1AddrStr)
	_, addUserErr := contract.AddUser(to, userAddr, testUser1Username)

	/* Assert */

	if registeredBeforeAdd {
		t.Errorf("User is already registered")
	}

	if addUserErr != nil {
		t.Errorf("Could not add user."+ addUserErr.Error())
	}

	// Check if registered via getter
	isUserRegistered, userRegisteredErr := contract.GetUserRegistered(nil, userAddr)

	if userRegisteredErr != nil  || !isUserRegistered {
		t.Errorf("Error retrieving if user is registered (false, expected true) "+ userRegisteredErr.Error())
	}

	// Check if admin via getter
	isUserAdmin, userAdminErr := contract.GetUserAdmin(nil, userAddr)

	if userAdminErr != nil || isUserAdmin{
		t.Errorf("User is admin (true, expected false) "+ userAdminErr.Error())
	}

	// Check user attempts via getter
	attempts, attemptsErr := contract.GetUserAttempts(nil, userAddr)

	if attemptsErr != nil || attempts > 0 {
		t.Errorf("Error, attempts="+string(attempts) +", expected 0"+ userAdminErr.Error())
	}

	// Check getters for relating pubkey and username
	getUserIdStr, getUserIdErr :=contract.GetUserId(nil, userAddr)
	getUserAddress, getUserAddressErr := contract.GetUserAddress(nil,testUser1Username)

	if getUserIdErr != nil {
		t.Errorf("Error relating user public address to username "+ getUserIdErr.Error())
	}

	if getUserAddressErr != nil {
		t.Errorf("Error relating username to public address"+ getUserAddressErr.Error())
	}

	getUserStr:= getUserAddress.String()

	if getUserIdStr != testUser1Username || getUserStr != testUser1AddrStr {
		t.Errorf("User ID and public key do not match")
	}

	// Check AuthContract

	userAuthContract,_ := contract.GetUserAuthContract(nil, userAddr)

	userAuthStr := userAuthContract.String()
	if userAuthStr == string2Address("0x0").String() {
		t.Errorf("Error retrieving user's AuthContract")
	}
}

// Try to register a user with a username taken, and a username that is already registered
func TestAddUserAlreadyTaken(t *testing.T) {
	/* Arrange: We need an initialized contract with a user in it */
	to,c := initializeValidClient(testRpcEndpoint, testChainId, testOwnerPrivKey)
	_, _, contract, _, _, _ :=deployAndInitialize(to,c, testOwnerAddrStr,testOwnerUsername)

	userAddr := string2Address(testUser1AddrStr) // Call contract method AddUser
	_, _ = contract.AddUser(to, userAddr, testUser1Username)


	/* Act: Try to add another different user with the same username */

	user1Addr := string2Address(testUser2AddrStr)
	_, err1 := contract.AddUser(to, user1Addr, testUser1Username)

	/* Act: Try to add a user already on the system */

	_, err2 := contract.AddUser(to, userAddr, testUser1Username+"1")

	/* Assert: is the new user registered with the same username as the existing one? */

	user1Registered, _ := contract.GetUserRegistered(nil, user1Addr)

	if user1Registered {
		t.Errorf("A second user was created with the same username")
	}

	/* Assert: check if a user is registered twice with to different usernames, it should be address 0 (username not taken) */

	shouldBeAddress0, _:= contract.GetUserAddress(nil, testUser1Username+"1")

	if shouldBeAddress0.String() == userAddr.String() {
		t.Errorf("Username was registered twice")
	}

	if err1 == nil || err2 == nil {
		t.Errorf(err1.Error()+"\n"+err2.Error())
	}
}
func TestAddUserWithoutPermission(t *testing.T) {
	/* Arrange: 2 clients and 3 users in the system, one is the contract owner, and a user tries to delete another non-admin */
	adminTransOps, adminClient := initializeValidClient(testRpcEndpoint, testChainId, testOwnerPrivKey)
	userTransOps, userClient := initializeValidClient(testRpcEndpoint, testChainId, testUser1PrivKey)
	contractAddress, _, contractAsAdmin, _, _, _ := deployAndInitialize(adminTransOps, adminClient, testOwnerPrivKey, testOwnerUsername)
	contractAsUser, _ := NewMain(contractAddress, userClient)

	contractAsAdmin.AddUser(adminTransOps, string2Address(testUser1AddrStr), testUser1Username)

	/* Act: the new user tries to remove a user from the system */
	_, addErr := contractAsUser.AddUser(userTransOps, string2Address(testUser2AddrStr), testUser2Username)

	// Third user is not expected to be registered
	user2IsRegistered, _ := contractAsUser.GetUserRegistered(nil, string2Address(testUser2AddrStr))

	if addErr == nil || user2IsRegistered {
		t.Errorf("Non-admin user was able to remove another user")
	}
}
func TestRemoveUserOk(t *testing.T ){
	/* Arrange: We need an initialized contract with a user in it */
	to,c := initializeValidClient(testRpcEndpoint, testChainId, testOwnerPrivKey)
	_, _, contract, _, _, _ :=deployAndInitialize(to,c, testOwnerAddrStr,testOwnerUsername)

	userAddr := string2Address(testUser1AddrStr) // Call contract method AddUser
	_, _ = contract.AddUser(to, userAddr, testUser1Username)

	/* Act: Call rmUser */
	_, rmError := contract.RmUser(to,userAddr, testUser1Username)

	/* Assert: check if user is not registered now */
	registeredAfterRemoving, _ := contract.GetUserRegistered(nil,userAddr)
	isAdmin, _ := contract.GetUserAdmin(nil,userAddr)
	loggedIn, _ := contract.GetUserLoggedIn(nil,userAddr)
	addressForUserName, _ := contract.GetUserAddress(nil, testUser1Username) // Address for the username after removing the user

	if registeredAfterRemoving ||
		rmError != nil ||
		isAdmin == true ||
		loggedIn == true ||
		addressForUserName.String() != string2Address("0x0").String()	{
		t.Errorf("User was not properly removed")
	}
}
func TestRemoveOwner(t *testing.T ){
	/* Arrange: We need an initialized contract */
	to,c := initializeValidClient(testRpcEndpoint, testChainId, testOwnerPrivKey)
	_, _, contract, _, _, _ :=deployAndInitialize(to,c, testOwnerAddrStr,testOwnerUsername)

	registeredBeforeRemoving, _ := contract.GetUserRegistered(nil, string2Address(testOwnerAddrStr))

	/* Act: try to remove the only user in the system */
	_, removeErr := contract.RmUser(to, string2Address(testOwnerAddrStr),testOwnerUsername)

	/* Assert: error ocurred and user is still registered */
	registeredAfterRemoving, _ := contract.GetUserRegistered(nil, string2Address(testOwnerAddrStr))

	if removeErr == nil || registeredBeforeRemoving != registeredAfterRemoving {
		t.Errorf("The owner was removed from the contract (?)")
	}

}
func TestRemoveMismatching(t *testing.T ){
	/* Arrange: We need an initialized contract */
	to,c := initializeValidClient(testRpcEndpoint, testChainId, testOwnerPrivKey)
	_, _, contract, _, _, _ :=deployAndInitialize(to,c, testOwnerAddrStr,testOwnerPrivKey)


	_, _ = contract.AddUser(to, string2Address(testUser1AddrStr), testUser1Username)

	/* Act: call rmUser on whoever with mismatching username and address, should fail */
	_, rmErr := contract.RmUser(to, string2Address(testUser1AddrStr),testUser2Username) // testuser2 is address 0x0...0 (unregistered)

	/* Assert: error ocurred, user is not removed */
	userRegistered,_  := contract.GetUserRegistered(nil, string2Address(testUser1AddrStr))

	if rmErr == nil || !userRegistered {
		t.Errorf("User was removed from the system: "+rmErr.Error())
	}
}
func TestRemoveWithoutPermission(t *testing.T) {
	/* Arrange: 2 clients and 3 users in the system, one is the contract owner, and a user tries to delete another non-admin */
	adminTransOps,adminClient := initializeValidClient(testRpcEndpoint, testChainId,testOwnerPrivKey)
	userTransOps,userClient := initializeValidClient(testRpcEndpoint, testChainId, testUser1PrivKey)
	contractAddress, _, contractAsAdmin, _, _, _ :=deployAndInitialize(adminTransOps,adminClient, testOwnerAddrStr,testOwnerPrivKey)
	contractAsUser, _ := NewMain(contractAddress,userClient)

	contractAsAdmin.AddUser(adminTransOps,string2Address(testUser1AddrStr), testUser1AddrStr)
	contractAsAdmin.AddUser(adminTransOps,string2Address(testUser2AddrStr), testUser2AddrStr)

	/* Act: the new user tries to remove a user from the system */
	_, rmErr := contractAsUser.RmUser(userTransOps,string2Address(testUser2AddrStr), testUser2AddrStr)

	// Third user is expected to be registered
	user2IsRegistered,_ := contractAsUser.GetUserRegistered(nil,string2Address(testUser2AddrStr))

	if rmErr == nil || !user2IsRegistered {
		t.Errorf("Non-admin user was able to remove another user")
	}


}
func TestPromoteUserOk(t *testing.T) {
	/* Arrange: We need an initialized contract with a user in it */
	to,c := initializeValidClient(testRpcEndpoint, testChainId, testOwnerPrivKey)
	contractAddress, _, contract, _, _, _ :=deployAndInitialize(to,c, testOwnerAddrStr,testOwnerUsername)

	userAddr := string2Address(testUser1AddrStr) // Call contract method AddUser
	_, _ = contract.AddUser(to, userAddr, testUser1Username)

	/* Act: promotion of user to admin */
	_, promoteErr := contract.PromoteUser(to,string2Address(testUser1AddrStr))

	/* Assert: user is admin now */
	isUserAdminNow, _ := contract.GetUserAdmin(nil,string2Address(testUser1AddrStr))

	if !isUserAdminNow || promoteErr != nil {
		t.Errorf("User was not promoted to admin")
	}

	/* Assert event */

	testEvents(contractAddress, c)

}
func TestPromoteUserNonExisting(t *testing.T) {
	/* Arrange: We need an initialized contract with a user in it */
	to,c := initializeValidClient(testRpcEndpoint, testChainId, testOwnerPrivKey)
	_, _, contract, _, _, _ :=deployAndInitialize(to,c, testOwnerAddrStr,testOwnerUsername)

	// User is not even on the system
	wasUserAdmin,_ := contract.GetUserAdmin(nil,string2Address(testUser1AddrStr))

	/* Act: promotion of user to admin */
	_, promoteErr := contract.PromoteUser(to,string2Address(testUser1AddrStr))

	/* Assert: user is admin now */
	isUserAdminNow, _ := contract.GetUserAdmin(nil,string2Address(testUser1AddrStr))

	if isUserAdminNow != wasUserAdmin || promoteErr == nil {
		t.Errorf("User not in the system was promoted to admin")
	}
}
func TestPromoteUserWithoutPermission(t *testing.T) {
	/* Arrange: 2 clients and 3 users in the system, one is the contract owner, and a user tries to delete another non-admin */
	adminTransOps,adminClient := initializeValidClient(testRpcEndpoint, testChainId,testOwnerPrivKey)
	userTransOps,userClient := initializeValidClient(testRpcEndpoint, testChainId, testUser1PrivKey)
	contractAddress, _, contractAsAdmin, _, _, _ :=deployAndInitialize(adminTransOps,adminClient, testOwnerAddrStr,testOwnerUsername)
	contractAsUser, _ := NewMain(contractAddress,userClient)

	contractAsAdmin.AddUser(adminTransOps,string2Address(testUser1AddrStr),testUser1Username)
	contractAsAdmin.AddUser(adminTransOps,string2Address(testUser2AddrStr),testUser2Username)

	/* Act: the new user tries to remove a user from the system */
	_, promoteErr := contractAsUser.PromoteUser(userTransOps,string2Address(testUser2AddrStr))

	/* Assert Third user is not expected to be admin */
	wasUserPromoted,_ := contractAsUser.GetUserAdmin(nil,string2Address(testUser2AddrStr))

	if promoteErr == nil || wasUserPromoted {
		t.Errorf("Non-admin user was able to remove another user")
	}
}
func TestDemoteAdminOk(t *testing.T) {
	/* Arrange: Valid contract with two administrators */
	adminTransOps,adminClient := initializeValidClient(testRpcEndpoint, testChainId,testOwnerPrivKey)
	contractAddress, _, contractAsAdmin, _, _, _ :=deployAndInitialize(adminTransOps,adminClient, testOwnerAddrStr, testOwnerUsername)

	contractAsAdmin.AddUser(adminTransOps,string2Address(testUser1AddrStr),testUser1Username)
	contractAsAdmin.PromoteUser(adminTransOps,string2Address(testUser1AddrStr)) // Promote user, now two admins in the system

	/* Act: owner demotes new admin */
	_, demoteErr := contractAsAdmin.DemoteAdmin(adminTransOps,string2Address(testUser1AddrStr))

	/* Assert: the other user should no longer be admin */
	isAdmin, _ := contractAsAdmin.GetUserAdmin(nil,string2Address(testUser1AddrStr))

	if demoteErr != nil || isAdmin {
		t.Errorf("User is admin (expected: not admin)")
	}

	/* Assert event */
	testEvents(contractAddress, adminClient)
}
func TestDemoteNonExistentUser(t *testing.T) {
	// User that tries to promote a non-admin user
	unregisteredUserAddrStr := testUser1AddrStr

	/* Arrange: Valid contract, no users */
	adminTransOps,adminClient := initializeValidClient(testRpcEndpoint, testChainId,testOwnerPrivKey)
	_, _, contractAsAdmin, _, _, _ :=deployAndInitialize(adminTransOps,adminClient, testOwnerAddrStr,testOwnerUsername)
	// Mini-assert: user is not registered -> isAdmin = false
	wasAdmin, _ := contractAsAdmin.GetUserAdmin(nil,string2Address(unregisteredUserAddrStr))

	/* Act: owner demotes a user that is not in the system */
	_, demoteErr := contractAsAdmin.DemoteAdmin(adminTransOps,string2Address(unregisteredUserAddrStr))

	/* Assert */
	isAdmin, _ := contractAsAdmin.GetUserAdmin(nil,string2Address(unregisteredUserAddrStr))

	if demoteErr == nil || wasAdmin != isAdmin {
		t.Errorf("User not in the system is admin (expected: not admin)")
	}
}
func TestDemoteWithoutPermission(t *testing.T) {
	/* Arrange: 2 clients and 3 users in the system, one is the contract owner, and a user tries to delete another non-admin */
	ownerTransOps, ownerClient := initializeValidClient(testRpcEndpoint, testChainId,testOwnerPrivKey)
	userTransOps,userClient := initializeValidClient(testRpcEndpoint, testChainId, testUser1PrivKey)
	contractAddress, _, contractAsAdmin, _, _, _ :=deployAndInitialize(ownerTransOps, ownerClient, testOwnerAddrStr,testOwnerUsername)
	contractAsUser, _ := NewMain(contractAddress,userClient)

	contractAsAdmin.AddUser(ownerTransOps,string2Address(testAdmin2AddrStr),testAdmin2Username)
	contractAsAdmin.PromoteUser(ownerTransOps,string2Address(testAdmin2AddrStr)) // Promoted to admin
	contractAsAdmin.AddUser(ownerTransOps,string2Address(testUser1AddrStr),testUser1Username)

	/* Act: the new user tries to demote a user from the system */
	_, demoteErr := contractAsUser.PromoteUser(userTransOps,string2Address(testUser1AddrStr))

	/* Assert Third user is not expected to be admin */
	wasAdminDemoted,_ := contractAsUser.GetUserAdmin(nil,string2Address(testUser1AddrStr))

	if demoteErr == nil || wasAdminDemoted {
		t.Errorf("Non-admin user was able to remove another user")
	}
}

func TestSetOwnerOk (t *testing.T) {
	/* Arrange: We need an initialized contract with a user in it */
	to,c := initializeValidClient(testRpcEndpoint, testChainId, testOwnerPrivKey)
	_, _, contract, _, _, _ :=deployAndInitialize(to,c, testOwnerAddrStr,testOwnerUsername)

	adminAddress := string2Address(testAdmin2AddrStr) // Call contract method AddUser
	_, _ = contract.AddUser(to, adminAddress, testAdmin2Username)
	_, _ = contract.PromoteUser(to,string2Address(testAdmin2AddrStr)) // Promote that user to admin

	/* Act: set that admin as the new owner */
	_, setOwnerErr := contract.SetOwner(to, adminAddress)

	/* Assert: the new owner has changed */
	newOwnerAddr,_ := contract.GetOwner(nil)


	if setOwnerErr !=nil  && newOwnerAddr.Hex() == testAdmin2AddrStr {
		t.Errorf("New owner was not set properly ("+setOwnerErr.Error()+")")
	}

}
func TestSetOwnerNotOwner (t *testing.T) {
	/* Arrange: We need an initialized contract with a user in it */
	ownerTops, ownerClient := initializeValidClient(testRpcEndpoint, testChainId, testOwnerPrivKey)
	contractAddress, _, contract, _, _, _ :=deployAndInitialize(ownerTops, ownerClient, testOwnerAddrStr,testOwnerUsername)

	// Retrieve contract's methods as new user
	userTops, userClient := initializeValidClient(testRpcEndpoint, testChainId, testUser1PrivKey)
	contractAsUser, _ := NewMain(contractAddress, userClient)

	/* Act: call setOwner as the new user */
	_, setOwnerErr := contractAsUser.SetOwner(userTops, string2Address(testUser1AddrStr))

	/* Assert: the new owner has changed */
	newOwnerAddr,_ := contract.GetOwner(nil)

	if setOwnerErr == nil  || newOwnerAddr.Hex() == testUser1AddrStr {
		t.Errorf("User was able to change the contract owner")
	}
}
// TestGeneratePass: for hashing using SoliditySha3 and check the result is the same
func TestGeneratePass(t *testing.T) {
	/* Arrange: deploy new contract */
	ownerTops, ownerClient := initializeValidClient(testRpcEndpoint, testChainId, testOwnerPrivKey)
	_, _, contract, _, _, _ :=deployAndInitialize(ownerTops, ownerClient, testOwnerAddrStr,testOwnerUsername)

	/* Act: call generatePass (client method) */
	pass,goPassHash := generatePass()

	/* Assert: call contract method generateHash and compare the result */
	solPassHash, _:= contract.GenerateHash(nil,pass)

	if pass < 10000 {
		t.Errorf("Password with less than 5 digits")
	}

	if !equalHash(goPassHash,solPassHash) {
		t.Errorf("Hashes don't match")
	}
}

func TestGetOTPOk(t *testing.T){

}
func TestGetOTPUserNotRegistered(t *testing.T){

}
func TestGetGetOTPUserOnline(t *testing.T){

}
func TestGetGetOTPUserLocked(t *testing.T){

}

func TestTryLoginOk(t *testing.T){

}
func TestTryLoginUserNotRegistered(t *testing.T){

}
func TestGetTryLoginUserOnline(t *testing.T){

}
func TestGetTryLoginUserLocked(t *testing.T){

}

func TestTryLogoutOk(t *testing.T){

}
func TestTryLogoutUserNotRegistered(t *testing.T){

}
func TestGetTryLogoutUserOffline(t *testing.T){

}