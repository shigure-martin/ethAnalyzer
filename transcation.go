package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func txTest() {
	client := getClient()

	// getBlockInfo(client)
	getAllTxInfo(client, *big.NewInt(8882896))
}

func getClient() *ethclient.Client {
	client, err := ethclient.Dial(providerUrl)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func getBlockInfo(client *ethclient.Client) {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("block head number: %s\n", header.Number.String())

	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx number: %d\n", block.Transactions().Len())
}

func getAllTxInfo(client *ethclient.Client, number big.Int) {
	var block *types.Block
	var err error
	if number.Cmp(big.NewInt(0)) > 0 {
		block, err = client.BlockByNumber(context.Background(), &number)
	} else {
		block, err = client.BlockByNumber(context.Background(), nil)
	}

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("block hash: ", block.Hash().Hex())

	for _, tx := range block.Transactions()[:1] {
		fmt.Println("transaction hash: ", tx.Hash().Hex())
		// fmt.Println(hexutil.Encode(tx.Data()))
		chainId, err := client.ChainID(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("chain id: %d\n", chainId)

		if from, err := types.Sender(types.NewEIP155Signer(chainId), tx); err == nil {
			fmt.Printf("tx from: %s\n", from.Hex())
		}

		if receipt, err := client.TransactionReceipt(context.Background(), tx.Hash()); err == nil {
			fmt.Println("receipt status: ", receipt.Status)
			fmt.Println("receipt logs: ", receipt.Logs)
		}
	}

	blockHash := common.HexToHash("0x84f233a0b8ea9b506e552122cb00e55ed80c0082d39edd43e7f41c26ae498e66")
	count, err := client.TransactionCount(context.Background(), blockHash)
	fmt.Println(count)
	if err != nil {
		log.Fatal(err)
	}

	for idx := uint(0); idx < 4; idx++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(tx.Hash().Hex())
	}
}
