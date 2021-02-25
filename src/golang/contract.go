package main

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

/*
// AddUser(opts *bind.TransactOpts, _addr common.Address, _id string) (*types.Transaction, error)
func contractAddUser (opts *bind.TransactOpts, _addr common.Address, _id string) (trans *types.Transaction,err error) {
	trans, err = mainObj.AddUser(opts, _addr, _id)
	return
}
*/