package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	// "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func txTest() {
	client := getClient()

	getBlockInfo(client)
	// getAllTxInfo(client, *big.NewInt(8882896))
	// txEth(client)
	// describeBlock()
}

func getClient() *ethclient.Client {
	client, err := ethclient.Dial(providerUrl)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func getBlockInfo(client *ethclient.Client) { //用于获取当前最新区块的head number，以及其中包含交易的数量
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

func txEth(client *ethclient.Client) {
	privateKey, err := crypto.HexToECDSA(myPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("from address: ", fromAddress)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("nonce at: ", nonce)

	value := big.NewInt(100000000000000000)
	gasLimit := uint64(21000)
	toAddress := common.HexToAddress("0xA738f13354ADaf4969aE7e8C8E5a975eee20a4A9")

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("suggest gas price: ", gasPrice)

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx sent: ", signedTx.Hash().Hex())
}

func txToken(client *ethclient.Client) *types.Transaction {
	privateKey, err := crypto.HexToECDSA(myPrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("from address: ", fromAddress)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("nonce at: ", nonce)

	value := big.NewInt(0)
	toAddress := common.HexToAddress("0xA738f13354ADaf4969aE7e8C8E5a975eee20a4A9")
	tokenAddress := common.HexToAddress("0x326C977E6efc84E512bB9C30f76E30c160eD06FB")

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("suggest gas price: ", gasPrice)

	transferSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferSignature)
	methodId := hash.Sum(nil)[:4]
	fmt.Println("methodId: ", hexutil.Encode(methodId))

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Println("padded address: ", hexutil.Encode(paddedAddress))

	amount := new(big.Int)
	amount.SetString("10000000000000000000", 10)

	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println("padded amount: ", hexutil.Encode(paddedAmount))

	var data []byte
	data = append(data, methodId...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit := uint64(50000)
	fmt.Println(gasLimit)

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("tx sent: ", signedTx.Hash().Hex())

	return signedTx
}

func describeBlock() {
	client, err := ethclient.Dial("wss://quaint-hardworking-breeze.ethereum-goerli.discover.quiknode.pro/31cc07938198cf8aa72ef7364dfc58e0578f8708/")
	if err != nil {
		log.Fatal(err)
	}

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex())
		}
	}
}

func rawTx(client *ethclient.Client, signedTx *types.Transaction) {

}
