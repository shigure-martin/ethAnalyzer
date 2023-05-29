package main

import (
	"bufio"
	"container/list"
	"context"
	"fmt"
	"log"
	"os"

	// "github.com/ethereum/go-ethereum"
	// "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"example/analyzer/signature"
)

func detect() {
	client := getClient(true)
	getTxs(client)
}

func readSigs() *list.List {
	file, err := os.Open("signature/sigs.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lis := list.New()
	for scanner.Scan() {
		line := scanner.Text()
		lis.PushBack(line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return lis
}

func checkSig(sigs []signature.Combine, data string) bool {
	result := false

	for _, combine := range sigs {
		if combine.Sigs == data {
			result = true
			break
		}
	}

	return result
}

func getTxs(client *ethclient.Client) {
	sigs := signature.GetCombines("signature/struct.gob") //readSigs()

	var block *types.Block
	var err error

	block, err = client.BlockByNumber(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("block number: ", block.Number())

	for _, tx := range block.Transactions() {
		if len(tx.Data()) == 0 {
			continue
		}
		data := hexutil.Encode(tx.Data())[:10]
		isDetected := checkSig(sigs, data)
		if isDetected {
			fmt.Println("tx: ", tx.Hash())
		}
	}
}
