package main

import (
	"bufio"
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
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

	block, err = client.BlockByNumber(context.Background(), big.NewInt(28880676))

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
			extractParam(hexutil.Encode(tx.Data()), combine, abiJSON)
		}
	}
}

func extractParam(_data string, _combine signature.Combine, _abiJSON abi.ABI) {

	method := _abiJSON.Methods[_combine.Name]

	// input := method.Inputs
	data := common.Hex2Bytes(_data[2:])

	arguments, err := method.Inputs.Unpack(data[4:])
	if err != nil {
		log.Fatal(err)
	}

	// var token_addr []string
	for _, position := range _combine.ParaPos {
		// token_addr = append(token_addr, arguments[position])
		fmt.Println(arguments[position])
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

	// for _, method := range abiJSON.Methods {
	// 	fmt.Println(method.Name)
	// }
	// method := abiJSON.Methods["swapExactTokensForETHSupportingFeeOnTransferTokens"]
	// fmt.Println("method: ", method.Inputs.NonIndexed())

	// params := method.Inputs.NonIndexed()[2].Type
	// fmt.Println("param type: ", params.String() == "address[]")

	// rawData := "0x791ac947000000000000000000000000000000000000000000000000016345785d8a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000007cf196415cdd1ef08ca2358a8282d33ba089b9f300000000000000000000000000000000000000000000000000000000647264c40000000000000000000000000000000000000000000000000000000000000002000000000000000000000000b4fbf271143f4fbf7b91a5ded31805e42b2208d6000000000000000000000000efadb5c4d4fc46c51e6639214aa95057d25a2573"
	// data := common.Hex2Bytes(rawData[2:])              // remove '0x'
	// paramValues, err := method.Inputs.Unpack(data[4:]) // remove the signatrue of function
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(paramValues)

	return abiJSON
}
