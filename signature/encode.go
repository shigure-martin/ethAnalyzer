package signature

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"
)

func Main() {
	file, err := os.Open("signature/methods.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lis := list.New()
	for scanner.Scan() {
		line := scanner.Text()
		transferSignature := []byte(line)
		hash := sha3.NewLegacyKeccak256()
		hash.Write(transferSignature)
		methodId := hash.Sum(nil)[:4]
		lis.PushBack(hexutil.Encode(methodId))
		fmt.Println(hexutil.Encode(methodId))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	seen := make(map[interface{}]bool)
	newL := list.New()
	for e := lis.Front(); e != nil; e = e.Next() {
		if _, ok := seen[e.Value]; !ok {
			seen[e.Value] = true
			newL.PushBack(e.Value)
		}
	}

	outFile, err := os.Create("signature/sigs.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	for item := newL.Front(); item != nil; item = item.Next() {
		fmt.Fprintln(outFile, item.Value)
	}
}
