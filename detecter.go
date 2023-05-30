package main

import (
	"bufio"
	"container/list"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	// "github.com/ethereum/go-ethereum"
	// "github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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

func checkSig(sigs []signature.Combine, data string) (bool, signature.Combine) {
	result := false

	result_combine := signature.Combine{}

	for _, combine := range sigs {
		if combine.Sigs == data {
			result_combine = combine
			result = true
			break
		}
	}

	return result, result_combine
}

func getTxs(client *ethclient.Client) {
	sigs := signature.GetCombines("signature/struct.gob") //readSigs()

	var block *types.Block
	var err error

	block, err = client.BlockByNumber(context.Background(), big.NewInt(9082721))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("block number: ", block.Number())

	for _, tx := range block.Transactions() {
		if len(tx.Data()) == 0 {
			continue
		}
		data := hexutil.Encode(tx.Data())[:10]
		isDetected, combine := checkSig(sigs, data)
		if isDetected {
			fmt.Println("tx: ", tx.Hash())
			fmt.Println("method: ", combine.Method)
			fmt.Println("raw data: ", hexutil.Encode(tx.Data()))
		}
	}
}

// func extractParam(_data string, _method string) {
// 	abiJSON := `[{"inputs": [{"internalType": "uint256","name": "amountIn","type": "uint256"},{"internalType": "uint256","name": "amountOutMin","type": "uint256"},{"internalType": "address[]","name": "path","type": "address[]"},{"internalType": "address","name": "to","type": "address"},{"internalType": "uint256","name": "deadline","type": "uint256"}],"name": "swapExactTokensForETHSupportingFeeOnTransferTokens","outputs": [],"stateMutability": "payable","type": "function"}]`

// 	abiObj, err := abi.JSON(strings.NewReader(abiJSON))
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	hexData := "0x19bea24b00000000000000000000000005569489c489a392e94090077ac3cc2d1a9e6a45000000000000000000000000006903471486b8b25aed908a347482e8e6e56c56bf0000000000000000000000000000000000000000000000000000000066d3d1200000000000000000000000000000000000000000000000018d4abe6d6c0000000000000000000000000000000000000000000000000000000000000002"

// 	// bytesData, _ := hexutil.Decode(hexData)
// 	method := abiObj.Methods["swapExactTokensForETHSupportingFeeOnTransferTokens"]

// 	input := method.Inputs
// 	data := common.Hex2Bytes(hexData)

// 	arguments, err := method.Inputs.Unpack(data)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
