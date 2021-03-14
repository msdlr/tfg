package main

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	solsha3 "github.com/msdlr/go-solidity-sha3"
	"math/rand"
	"time"
)

var mainObj *Main // Main object used for invoking contract's methods

// Deploys the an instance of generalContract and inicializes it with the public address of the creator
// For clients that don't deploy a contract but just want to access: NewMain(address common.Address, backend bind.ContractBackend) (*Main, error)
func deployAndInitialize(myTrOps *bind.TransactOpts, myClient *ethclient.Client, ownerPubKeyStr string, ownerUserName string)(addr common.Address, deployTrans *types.Transaction, main *Main, deployError error, initTrans *types.Transaction, initError error){
	addr, deployTrans, main, deployError = DeployMain(myTrOps, myClient)
	mainObj = main // We need this main object to be preserved!!!
	/*
	if err == nil {
		fmt.Println("## NEW CONTRACT DEPLOYED ##")
		fmt.Println("Address:\t\t", addr.Hex())
		fmt.Println("Transaction hash:\t", trans.Hash())
		fmt.Println("Gas Used:\t\t", trans.Gas(), "(price:", trans.GasPrice(), ")")
		fmt.Println("Nonce:\t\t\t", trans.Nonce())
	}
	 */
	myPubKey := string2Address(ownerPubKeyStr)
	initTrans,initError = mainObj.Initialize(myTrOps, myPubKey, ownerUserName)
	return
}

// generatePass generates a hashed OTP number (10000-65535) to store on a smart contract
func generatePass() (pass uint16, passHash32 [32]byte) {
	// Generate the random password
	seed := rand.NewSource(time.Now().UnixNano())
	pass = 1

	// Generate a random number bigger than 1E5 (so that there are 5 digits)
	for true {
		num := uint16(rand.New(seed).Uint32())
		if num >= 10000 {
			// If found we've finished
			pass = num
			break
		}
	}

	// Hash it
	passHash := solsha3.SoliditySHA3([]string{"uint16"}, pass)
	return pass, sliceToArray32(passHash)
}

// sliceToArray32 is the opposite of array[:] which transforms an array into a slice
func sliceToArray32(in []byte) (out [32]byte){
	for i:=0;i<32;i++{
		out[i]= in[i]
	}
	return
}

func generateAndStorePassHash(){

}

// genPassAndStoreHash is a wrapper that generates a pass Hash and writes it into a contract. OTP pass generation and writing on a contract should be performed together
func genPassAndStoreHash(contractMain *Main, transactionOpts *bind.TransactOpts) (pass uint16, passHash [32]byte, setHashPassError error) {
	pass, passHash = generatePass()
	_, setHashPassError = contractMain.SetHashPass(transactionOpts, passHash)
	return
}

/*
// AddUser(opts *bind.TransactOpts, _addr common.Address, _id string) (*types.Transaction, error)
func contractAddUser (opts *bind.TransactOpts, _addr common.Address, _id string) (trans *types.Transaction,err error) {
	trans, err = mainObj.AddUser(opts, _addr, _id)
	return
}
*/