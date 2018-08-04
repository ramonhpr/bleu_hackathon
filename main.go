package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const TID = "0xb1ed364e4333aae1da4a901d5231244ba6a35f9421d4607f7cb90d60bf45578a"
const URL = "https://mainnet.infura.io"

func main() {
	rpcCli, errRpcClient := rpc.Dial(URL)
	if errRpcClient != nil {
		log.Fatal("Dial erro: Something went wrong")
	}
	var cli = ethclient.NewClient(rpcCli)
	var ctx = context.Background()

	tx, isPending, err := cli.TransactionByHash(ctx, common.HexToHash(TID))

	if err != nil {
		log.Fatalf("TransactionByHash error: %v\n", err)
	} else if isPending == false {
		fmt.Println(string(tx.Data()[:len(tx.Data())]))
	}
}
