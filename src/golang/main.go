package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Dial address: ganache in localhost
	conn, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal("Error reaching RCP.", err)
	}

	ctx := context.Backgraoung()
	tx, pending, _ := conn.TranstactionByHash(ctx, common.HexToHash(""))

	if !pending {
		fmt.Println(tx)
	}
}
