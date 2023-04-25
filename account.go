package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math"
	"math/big"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

func accountTest() {
	client, err := ethclient.Dial(providerUrl)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")

	fmt.Println("we are getting account balance.")
	getAccountBalance("0xe9A9fb0554af4FF167F2b33a64a358Ea4C2D6Aad", client)

	fmt.Println("we are generating wallet.")
	generateWallet()

	fmt.Println("we are verifing address.")
	addressVerify("0x592aA75097598Bca17978006583de06Dd0477768", client)
	_ = client
}

func addressVerify(address string, client *ethclient.Client) {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	isAddress := re.MatchString(address)
	if isAddress {
		addr := common.HexToAddress(address)
		bytecode, err := client.CodeAt(context.Background(), addr, nil)
		if err != nil {
			log.Fatal(err)
		}

		isContract := len(bytecode) > 0

		fmt.Printf("is contract: %v\n", isContract)
	}
}

func getAccountBalance(address string, client *ethclient.Client) {
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	fmt.Println(ethValue)
}

func generateWallet() {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:])

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:])

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))
}
