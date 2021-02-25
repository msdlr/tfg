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

var (
	ownerAddressStr = "0xFDb59BC058eFde421AdF049F27d3A03a4cedea2f" // [0]
	ownerPrivKey = "ad92041b60126af952f8320b473ccb555d7274a53f1c27e12d2f1ea8aaecda7b"

	admin2AddressStr = "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5" // [1]
	admin2PrivKey = "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"

	user1AddressStr = "0x12b3C6913a8eE35D1e0462f16Ac0aA6B6205a91a" // [2]
	user1PrivKey = "e7c911fedc61cc1fd1a7a1cb84fd449562709cfa16a39f228cca07158c7307fb"

	user2AddressStr = "0xD03A8E7E2265CD8239F34909324F98F00496EA31" // [3]
	user2PrivKey = "7aa5be5263617d40346f8dc8d32f59a6cc6443bbf8d164bc1b89170f2d0679af"

	chainIdStr = "chainId"
	chainId uint16 = 5777
	rpcEndpoint = "http://localhost:7545"
)

func TestDeployOk(t *testing.T){
	/* Arrange: create valid client */
	to,c := initializeValidClient(rpcEndpoint,chainId,"b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4")

	/* Act: deploy and initialize general contract */
	addr, deployTrans, main, deployError, initTrans, initError :=deployAndInitialize(to,c,"0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5","TestOwner")

	/* Assert contract created properly */
	addrStr := addr.String()
	zeroStr := string2Address("0x0").Hex() // Address is not 0x00...0
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

	if addrGetStr != addrStr || err != nil{
		t.Errorf("Created contract address differs from getter "+err.Error())
	}
}

// region function: addUser

func TestAddUserOk(t *testing.T){
	/* Arrange: We need an initialized contract */

	ownerAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5" // [1]
	ownerPrivStr := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"

	newuserPubStr := "0x12b3C6913a8eE35D1e0462f16Ac0aA6B6205a91a" // [2]
	newUserName := "testUser"

	to,c := initializeValidClient(rpcEndpoint,chainId,ownerPrivStr)
	_, _, main, _, _, _ :=deployAndInitialize(to,c, ownerAddrStr,"TestOwner")

	registeredBeforeAdd, _ := main.GetUserRegistered(nil, string2Address(newuserPubStr))

	userAddr := string2Address(newuserPubStr)
	_, addUserErr := main.AddUser(to, userAddr, newUserName)

	/* Assert */

	if registeredBeforeAdd {
		t.Errorf("User is already registered")
	}

	if addUserErr != nil {
		t.Errorf("Could not add user."+ addUserErr.Error())
	}

	// Check if registered via getter
	isUserRegistered, userRegisteredErr := main.GetUserRegistered(nil, userAddr)

	if userRegisteredErr != nil  || !isUserRegistered {
		t.Errorf("Error retrieving if user is registered (false, expected true) "+ userRegisteredErr.Error())
	}

	// Check if admin via getter
	isUserAdmin, userAdminErr := main.GetUserAdmin(nil, userAddr)

	if userAdminErr != nil || isUserAdmin{
		t.Errorf("User is admin (true, expected false) "+ userAdminErr.Error())
	}

	// Check user attempts via getter
	attempts, attemptsErr := main.GetUserAttempts(nil, userAddr)

	if attemptsErr != nil || attempts > 0 {
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
func TestAddUserAlreadyTaken(t *testing.T) {
	/* Arrange: We need an initialized contract with a user in it */

	ownerAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5" // [1]
	ownerPrivStr := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"

	userPubStr := "0x12b3C6913a8eE35D1e0462f16Ac0aA6B6205a91a" // [2]
	userNickname := "testUser"

	newuserPubStr := "0xa0A9e0409f8A0e03f41e1AAd5Bb19E86C4fE5Acc" // [3]


	to,c := initializeValidClient(rpcEndpoint,chainId,ownerPrivStr)
	_, _, main, _, _, _ :=deployAndInitialize(to,c, ownerAddrStr,"TestOwner")

	userAddr := string2Address(userPubStr) // Call contract method AddUser
	_, _ = main.AddUser(to, userAddr, userNickname)


	/* Act: Try to add another different user with the same username */

	newUserAddr := string2Address(newuserPubStr)
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

func TestAddUserWithoutPermission(t *testing.T) {
	// ContractOwner
	ownerPrivKey := "ad92041b60126af952f8320b473ccb555d7274a53f1c27e12d2f1ea8aaecda7b"
	userPrivKey := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"
	// Non-admin user that tries to add a new user
	userAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5"
	ownerAddrStr := "0xFDb59BC058eFde421AdF049F27d3A03a4cedea2f"

	// User that is tried to be added
	otherUserAddrStr := "0x045C24525C46DBfaA8CfF3EA6C48a0e877777bFF"

	/* Arrange: 2 clients and 3 users in the system, one is the contract owner, and a user tries to delete another non-admin */
	adminTransOps, adminClient := initializeValidClient(rpcEndpoint, chainId, ownerPrivKey) // [0]
	userTransOps, userClient := initializeValidClient(rpcEndpoint, chainId, userPrivKey)    // [1]
	contractAddress, _, adminMain, _, _, _ := deployAndInitialize(adminTransOps, adminClient, ownerAddrStr, "TestOwner")
	userMain, _ := NewMain(contractAddress, userClient)

	adminMain.AddUser(adminTransOps, string2Address(userAddrStr), "newUser")

	/* Act: the new user tries to remove a user from the system */
	_, addErr := userMain.AddUser(userTransOps, string2Address(otherUserAddrStr), "newUser9")

	// Third user is not expected to be registered
	otherUserIsRegistered, _ := userMain.GetUserRegistered(nil, string2Address(otherUserAddrStr))

	if addErr == nil || otherUserIsRegistered {
		t.Errorf("Non-admin user was able to remove another user")
	}
}
//endregion

// region function: rmUser
func TestRemoveUserOk(t *testing.T ){
	/* Arrange: We need an initialized contract with a user in it */

	ownerAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5" // [1]
	ownerPrivStr := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"

	userPubStr := "0x12b3C6913a8eE35D1e0462f16Ac0aA6B6205a91a" // [2]
	userNickname := "testUser"

	to,c := initializeValidClient(rpcEndpoint,chainId,ownerPrivStr)
	_, _, main, _, _, _ :=deployAndInitialize(to,c, ownerAddrStr,"TestOwner")

	userAddr := string2Address(userPubStr) // Call contract method AddUser
	_, _ = main.AddUser(to, userAddr, userNickname)

	/* Act: Call rmUser */
	_, rmError := main.RmUser(to,userAddr,userNickname)

	/* Assert: check if user is not registered now */
	registeredAfterRemoving, _ := main.GetUserRegistered(nil,userAddr)
	isAdmin, _ := main.GetUserAdmin(nil,userAddr)
	loggedIn, _ := main.GetUserLoggedIn(nil,userAddr)
	addressForUserName, _ := main.GetUserAddress(nil,userNickname) // Address for the username after removing the user

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

	ownerAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5" // [1]
	ownerPrivStr := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"
	to,c := initializeValidClient(rpcEndpoint,chainId,ownerPrivStr)
	_, _, main, _, _, _ :=deployAndInitialize(to,c, ownerAddrStr,"TestOwner")

	registeredBeforeRemoving, _ := main.GetUserRegistered(nil, string2Address(ownerAddrStr))

	/* Act: try to remove the only user in the system */
	_, removeErr := main.RmUser(to, string2Address(ownerAddrStr),"TestOwner")

	/* Assert: error ocurred and user is still registered */
	registeredAfterRemoving, _ := main.GetUserRegistered(nil, string2Address(ownerAddrStr))

	if removeErr == nil || registeredBeforeRemoving != registeredAfterRemoving {
		t.Errorf("The owner was removed from the contract (?)")
	}

}

func TestRemoveMismatching(t *testing.T ){
	/* Arrange: We need an initialized contract */

	ownerAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5" // [1]
	ownerPrivStr := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"
	to,c := initializeValidClient(rpcEndpoint,chainId,ownerPrivStr)
	_, _, main, _, _, _ :=deployAndInitialize(to,c, ownerAddrStr,"TestOwner")

	userPubStr := "0x12b3C6913a8eE35D1e0462f16Ac0aA6B6205a91a" // [2]
	userNickname := "testUser"

	_, _ = main.AddUser(to, string2Address(userPubStr), userNickname)

	/* Act: call rmUser on whoever with mismatching username and address, should fail */
	_, rmErr := main.RmUser(to, string2Address(userPubStr),"testUser2") // testUser2 is address 0x0...0 (unregistered)

	/* Assert: error ocurred, user is not removed */
	userRegistered,_  := main.GetUserRegistered(nil, string2Address(userPubStr))

	if rmErr == nil || !userRegistered {
		t.Errorf("User was removed from the system: "+rmErr.Error())
	}
}

func TestRemoveWithoutPermission(t *testing.T) {
	// ContractOwner
	ownerPrivKey := "ad92041b60126af952f8320b473ccb555d7274a53f1c27e12d2f1ea8aaecda7b"
	userPrivKey := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"
	// User that tries to remove a non-admin user
	userAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5"
	ownerAddrStr := "0xFDb59BC058eFde421AdF049F27d3A03a4cedea2f"

	// User that is tried to be removed
	otherUserAddrStr := "0x045C24525C46DBfaA8CfF3EA6C48a0e877777bFF"

	/* Arrange: 2 clients and 3 users in the system, one is the contract owner, and a user tries to delete another non-admin */
	adminTransOps,adminClient := initializeValidClient(rpcEndpoint,chainId,ownerPrivKey) // [0]
	userTransOps,userClient := initializeValidClient(rpcEndpoint,chainId,userPrivKey) // [1]
	contractAddress, _, adminMain, _, _, _ :=deployAndInitialize(adminTransOps,adminClient, ownerAddrStr,"TestOwner")
	userMain, _ := NewMain(contractAddress,userClient)

	adminMain.AddUser(adminTransOps,string2Address(userAddrStr),"newUser")
	adminMain.AddUser(adminTransOps,string2Address(otherUserAddrStr),"newUser9") // [9]

	/* Act: the new user tries to remove a user from the system */
	_, rmErr := userMain.RmUser(userTransOps,string2Address(otherUserAddrStr),"newUser9")

	// Third user is expected to be registered
	otherUserIsRegistered,_ := userMain.GetUserRegistered(nil,string2Address(otherUserAddrStr))

	if rmErr == nil || !otherUserIsRegistered {
		t.Errorf("Non-admin user was able to remove another user")
	}


}
// endregion

// region function: promoteUser
func TestPromoteUserOk(t *testing.T) {
	/* Arrange: We need an initialized contract with a user in it */

	ownerAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5" // [1]
	ownerPrivStr := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"

	userPubStr := "0x12b3C6913a8eE35D1e0462f16Ac0aA6B6205a91a" // [2]
	userNickname := "testUser"

	to,c := initializeValidClient(rpcEndpoint,chainId,ownerPrivStr)
	_, _, main, _, _, _ :=deployAndInitialize(to,c, ownerAddrStr,"TestOwner")

	userAddr := string2Address(userPubStr) // Call contract method AddUser
	_, _ = main.AddUser(to, userAddr, userNickname)

	/* Act: promotion of user to admin */
	_, promoteErr := main.PromoteUser(to,string2Address(userPubStr))

	/* Assert: user is admin now */
	isUserAdminNow, _ := main.GetUserAdmin(nil,string2Address(userPubStr))

	if !isUserAdminNow || promoteErr != nil {
		t.Errorf("User was not promoted to admin")
	}

	/* Assert event */
}

func TestPromoteUserNonExisting(t *testing.T) {
	/* Arrange: We need an initialized contract with a user in it */

	ownerAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5" // [1]
	ownerPrivStr := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"

	userPubStr := "0x12b3C6913a8eE35D1e0462f16Ac0aA6B6205a91a" // [2]

	to,c := initializeValidClient(rpcEndpoint,chainId,ownerPrivStr)
	_, _, main, _, _, _ :=deployAndInitialize(to,c, ownerAddrStr,"TestOwner")

	// User is not even on the system
	wasUserAdmin,_ := main.GetUserAdmin(nil,string2Address(userPubStr))

	/* Act: promotion of user to admin */
	_, promoteErr := main.PromoteUser(to,string2Address(userPubStr))

	/* Assert: user is admin now */
	isUserAdminNow, _ := main.GetUserAdmin(nil,string2Address(userPubStr))

	if isUserAdminNow != wasUserAdmin || promoteErr == nil {
		t.Errorf("User not in the system was promoted to admin")
	}
}

func TestPromoteUserWithoutPermission(t *testing.T) {
	// ContractOwner
	ownerAddrStr := "0xFDb59BC058eFde421AdF049F27d3A03a4cedea2f"
	ownerPrivKey := "ad92041b60126af952f8320b473ccb555d7274a53f1c27e12d2f1ea8aaecda7b"
	// User that tries to remove a non-admin user
	userAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5"
	userPrivKey := "b7e6a03909b31f05c90354dd1a2bf61df5f223198c62551127250fdce2f6ffd4"

	// User that is tried to be promoted
	otherUserAddrStr := "0x045C24525C46DBfaA8CfF3EA6C48a0e877777bFF"

	/* Arrange: 2 clients and 3 users in the system, one is the contract owner, and a user tries to delete another non-admin */
	adminTransOps,adminClient := initializeValidClient(rpcEndpoint,chainId,ownerPrivKey) // [0]
	userTransOps,userClient := initializeValidClient(rpcEndpoint,chainId,userPrivKey) // [1]
	contractAddress, _, adminMain, _, _, _ :=deployAndInitialize(adminTransOps,adminClient, ownerAddrStr,"TestOwner")
	userMain, _ := NewMain(contractAddress,userClient)

	adminMain.AddUser(adminTransOps,string2Address(userAddrStr),"newUser")
	adminMain.AddUser(adminTransOps,string2Address(otherUserAddrStr),"newUser9") // [9]

	/* Act: the new user tries to remove a user from the system */
	_, promoteErr := userMain.PromoteUser(userTransOps,string2Address(otherUserAddrStr))

	/* Assert Third user is not expected to be admin */
	wasUserPromoted,_ := userMain.GetUserAdmin(nil,string2Address(otherUserAddrStr))

	if promoteErr == nil || wasUserPromoted {
		t.Errorf("Non-admin user was able to remove another user")
	}
}

// endregion

// region func demoteAdmin
func TestDemoteAdminOk(t *testing.T) {
	// ContractOwner
	ownerAddrStr := "0xFDb59BC058eFde421AdF049F27d3A03a4cedea2f"
	ownerPrivKey := "ad92041b60126af952f8320b473ccb555d7274a53f1c27e12d2f1ea8aaecda7b"
	// User that tries to remove a non-admin user
	userAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5"

	/* Arrange: Valid contract with two administrators */
	adminTransOps,adminClient := initializeValidClient(rpcEndpoint,chainId,ownerPrivKey) // [0]
	_, _, adminMain, _, _, _ :=deployAndInitialize(adminTransOps,adminClient, ownerAddrStr,"TestOwner")

	adminMain.AddUser(adminTransOps,string2Address(userAddrStr),"newUser")
	adminMain.PromoteUser(adminTransOps,string2Address(userAddrStr)) // Promote user, now two admins in the system

	/* Act: owner demotes new admin */
	_, demoteErr := adminMain.DemoteAdmin(adminTransOps,string2Address(userAddrStr))

	/* Assert: the other user should no longer be admin */
	isAdmin, _ := adminMain.GetUserAdmin(nil,string2Address(userAddrStr))

	if demoteErr != nil || isAdmin {
		t.Errorf("User is admin (expected: not admin)")
	}
}

func TestDemoteNonExistentUser(t *testing.T) {
	// ContractOwner
	ownerAddrStr := "0xFDb59BC058eFde421AdF049F27d3A03a4cedea2f"
	ownerPrivKey := "ad92041b60126af952f8320b473ccb555d7274a53f1c27e12d2f1ea8aaecda7b"
	// User that tries to promote a non-admin user
	userAddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5"

	/* Arrange: Valid contract, no users */
	adminTransOps,adminClient := initializeValidClient(rpcEndpoint,chainId,ownerPrivKey) // [0]
	_, _, adminMain, _, _, _ :=deployAndInitialize(adminTransOps,adminClient, ownerAddrStr,"TestOwner")
	// Mini-assert: user is not registered -> isAdmin = false
	wasAdmin, _ := adminMain.GetUserAdmin(nil,string2Address(userAddrStr))

	/* Act: owner demotes a user that is not in the system */
	_, demoteErr := adminMain.DemoteAdmin(adminTransOps,string2Address(userAddrStr))

	/* Assert */
	isAdmin, _ := adminMain.GetUserAdmin(nil,string2Address(userAddrStr))

	if demoteErr == nil || wasAdmin != isAdmin {
		t.Errorf("User not in the system is admin (expected: not admin)")
	}
}

func TestDemoteWithoutPermission(t *testing.T) {
	// ContractOwner
	ownerAddrStr := "0xFDb59BC058eFde421AdF049F27d3A03a4cedea2f"
	ownerPrivKey := "ad92041b60126af952f8320b473ccb555d7274a53f1c27e12d2f1ea8aaecda7b"
	// User that tries to remove a non-admin user
	admin2AddrStr := "0xe065fAE3BaF52ee871C956E55C88548E4d17F5A5"

	// User that is tried to be promoted
	userAddrStr := "0xa0A9e0409f8A0e03f41e1AAd5Bb19E86C4fE5Acc"
	userPrivkey := "81ed56f988ec5cf9c1adccc35ae4d0aa118dda78cfaa80ccfc29c4aec2ca35d0"

	/* Arrange: 2 clients and 3 users in the system, one is the contract owner, and a user tries to delete another non-admin */
	ownerTransOps, ownerClient := initializeValidClient(rpcEndpoint,chainId,ownerPrivKey) // [0]
	userTransOps,userClient := initializeValidClient(rpcEndpoint,chainId, userPrivkey) // [4]
	contractAddress, _, adminMain, _, _, _ :=deployAndInitialize(ownerTransOps, ownerClient, ownerAddrStr,"TestOwner")
	userMain, _ := NewMain(contractAddress,userClient)

	adminMain.AddUser(ownerTransOps,string2Address(admin2AddrStr),"newAdmin")
	adminMain.PromoteUser(ownerTransOps,string2Address(admin2AddrStr)) // Promoted to admin
	adminMain.AddUser(ownerTransOps,string2Address(userAddrStr),"newUser")

	/* Act: the new user tries to demote a user from the system */
	_, demoteErr := userMain.PromoteUser(userTransOps,string2Address(userAddrStr))

	/* Assert Third user is not expected to be admin */
	wasAdminDemoted,_ := userMain.GetUserAdmin(nil,string2Address(userAddrStr))

	if demoteErr == nil || wasAdminDemoted {
		t.Errorf("Non-admin user was able to remove another user")
	}
}
// endregion