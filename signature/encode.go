package signature

import (
	"bufio"
	"bytes"
	"container/list"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"
)

type Combine struct {
	Sigs   string
	Method string
}

func readMethod(address string) []Combine {
	file, err := os.Open(address)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// lis := list.New()
	lis := []Combine{}
	for scanner.Scan() {
		line := scanner.Text()
		transferSignature := []byte(line)
		hash := sha3.NewLegacyKeccak256()
		hash.Write(transferSignature)
		methodId := hash.Sum(nil)[:4]

		var combine Combine
		combine.Method = line
		combine.Sigs = hexutil.Encode(methodId)

		// lis.PushBack(combine)
		lis = append(lis, combine)
		// fmt.Println(hexutil.Encode(methodId))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return lis
}

// Element de-duplication
func deDuplicate(lis *list.List) *list.List {
	seen := make(map[interface{}]bool)
	newL := list.New()
	for e := lis.Front(); e != nil; e = e.Next() {
		if _, ok := seen[e.Value]; !ok {
			seen[e.Value] = true
			newL.PushBack(e.Value)
		}
	}

	return newL
}

func printSig(address string, lis []Combine) {
	// gob.RegisterName("Combine", Combine{}) //注册gob

	outFile, err := os.Create(address)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// var buf bytes.Buffer
	encoder := gob.NewEncoder(outFile)

	if err := encoder.Encode(lis); err != nil {
		log.Fatal(err)
		return
	}

	// if err := ioutil.WriteFile(address, buf.Bytes(), 0644); err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// for item := lis.Front(); item != nil; item = item.Next() {
	// 	fmt.Fprintln(outFile, item.Value)
	// }
}

func readStruct(address string) {
	file, err := ioutil.ReadFile(address)
	if err != nil {
		log.Fatal(err)
		return
	}

	var result *list.List
	decoder := gob.NewDecoder(bytes.NewReader(file))
	if err := decoder.Decode(&result); err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(result)
	// scanner := bufio.NewScanner(file)

	// for scanner.Scan() {
	// 	var combine Combine
	// 	combine = scanner.Text()
	// }
}

func Main() {
	lis := readMethod("signature/methods.txt")

	// newL := deDuplicate(lis)

	printSig("signature/struct.gob", lis)

	// readStruct("signature/struct.txt")
}
