package main

import (
	"bufio"
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"log"

	// "math/big"
	"os"

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
	// readABI()

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

	abiJSON := readABI()

	var block *types.Block
	var err error

	block, err = client.BlockByNumber(context.Background(), nil) //big.NewInt(28880676))

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
			fmt.Println("tx hash: ", tx.Hash())
			// fmt.Println("method: ", combine.Method)
			// fmt.Println("raw data: ", hexutil.Encode(tx.Data()))
			extractParam(tx.Data(), combine, abiJSON, client)
		}
	}
}

func extractParam(_data []byte, _combine signature.Combine, _abiJSON abi.ABI, client *ethclient.Client) {

	method := _abiJSON.Methods[_combine.Name]

	// input := method.Inputs
	// data := common.Hex2Bytes(_data[2:])

	arguments, err := method.Inputs.Unpack(_data[4:])
	if err != nil {
		log.Fatal(err)
	}

	// var token_addr []string
	for _, position := range _combine.ParaPos {
		fmt.Println(arguments[position])

		addr, ok := arguments[position].([]common.Address)

		if !ok {
			fmt.Println("addresses error")
		}

		// TODO:
		token, err := NewErc20(addr[0], client)

		if err != nil {
			fmt.Println(err)
		}

		name, err := token.Name(nil)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(name)
	}

}

func readABI() abi.ABI {
	abiFile, err := os.ReadFile("signature/abi.json")
	if err != nil {
		log.Fatal(err)
	}

	var abiJSON abi.ABI
	err = json.Unmarshal(abiFile, &abiJSON)
	if err != nil {
		log.Fatal(err)
	}

	return abiJSON
}
