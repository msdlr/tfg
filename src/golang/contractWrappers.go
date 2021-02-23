package main

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Deploys the an instance of generalContract and inicializes it with the public address of the creator
func deployAndInitialize(myTrOps *bind.TransactOpts, myClient *ethclient.Client, ownerPubKeyStr string, ownerUserName string)(addr common.Address, deployTrans *types.Transaction, main *Main, deployError error, initTrans *types.Transaction, initError error){
	addr, deployTrans, main, deployError = DeployMain(myTrOps, myClient)
	/*
	if err == nil {
		fmt.Println("## NEW CONTRACT DEPLOYED ##")
		fmt.Println("Address:\t\t", addr.Hex())
		fmt.Println("Transaction hash:\t", trans.Hash())
		fmt.Println("Gas Used:\t\t", trans.Gas(), "(price:", trans.GasPrice(), ")")
		fmt.Println("Nonce:\t\t\t", trans.Nonce())
	}
	 */
	myPubKey := publicAddressFromString(ownerPubKeyStr)
	initTrans,initError = main.Initialize(myTrOps, myPubKey, ownerUserName)
	return
}
