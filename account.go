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

	// fmt.Println("we are generating wallet...")
	// for i := int(0); i < 1000; i++ {
	// 	finded := generateWalletAndTest(client)
	// 	fmt.Println(i)
	// 	if finded {
	// 		break
	// 	}
	// }

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

func getAccountBalance(address string, client *ethclient.Client) big.Float {
	account := common.HexToAddress(address)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	fmt.Println(ethValue)

	return *ethValue
}

func generateWallet() (string, string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyStr := hexutil.Encode(privateKeyBytes)
	fmt.Println("private key: ", privateKeyStr)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println("public key: ", hexutil.Encode(publicKeyBytes)[4:])

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println("address: ", address)

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:]))

	return privateKeyStr, address
}

func generateWalletAndTest(client *ethclient.Client) bool {
	privateKey, address := generateWallet()
	ethValue := getAccountBalance(address, client)

	if ethValue.Cmp(big.NewFloat(0)) > 0 {
		fmt.Println("private key: ", privateKey)
		fmt.Println("address: ", address)
		fmt.Println("eth value: ", ethValue)
		return true
	}
	return false
}
