package signature

import (
	"bufio"
	"container/list"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/sha3"
)

type Combine struct {
	Sigs    string
	Method  string
	paraPos []int
}

func readMethod(address string) []Combine {
	file, err := os.Open(address)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	lis := []Combine{}
	for scanner.Scan() {
		line := scanner.Text()
		subLines := strings.Split(line, " ")

		paraPos := []int{}
		ints := strings.Split(subLines[0], ",")
		for _, pos := range ints {
			para, _ := strconv.Atoi(pos)
			paraPos = append(paraPos, para)
		}

		transferSignature := []byte(subLines[1])
		hash := sha3.NewLegacyKeccak256()
		hash.Write(transferSignature)
		methodId := hash.Sum(nil)[:4]

		var combine Combine
		combine.Method = subLines[1]
		combine.Sigs = hexutil.Encode(methodId)
		combine.paraPos = paraPos

		lis = append(lis, combine)
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
}

func readStruct(address string) {
	file, err := os.Open(address)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lis []Combine
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&lis); err != nil {
		log.Fatal(err)
		return
	}

	for _, combine := range lis {
		fmt.Println(combine.Sigs + " " + combine.Method)
	}
}

func GetCombines(address string) []Combine {
	file, err := os.Open(address)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var result []Combine
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&result); err != nil {
		log.Fatal(err)
		return []Combine{}
	}

	return result
}

func Main() {
	lis := readMethod("signature/methods.txt")

	// newL := deDuplicate(lis)
	for _, combine := range lis {
		fmt.Println(combine)
	}

	// printSig("signature/struct.gob", lis)

	// readStruct("signature/struct.gob")
}
