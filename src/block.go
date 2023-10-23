package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, block *Block) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("No response from request")
	}
	//fmt.Printf("%s", resp.Body)
	err = json.NewDecoder(resp.Body).Decode(&block)
	if err != nil {
		log.Fatal(err)
	}
}

type Block struct {
	Hash       string `json:"hash"`
	Ver        int64  `json:"ver"`
	Prev_block string `json:"prev_block"`
	Mrkl_root  string ` json:"mrkl_root"`
	Bits       int64  `json:"bits"`
	Time       int64  `json:"time"`
	Nonce      int64  `json:"nonce"`
}
type Block_string struct {
	Hash       string `json:"hash"`
	Ver        string `json:"ver"`
	Prev_block string `json:"prev_block"`
	Mrkl_root  string ` json:"mrkl_root"`
	Bits       string `json:"bits"`
	Time       string `json:"time"`
	Nonce      string `json:"nonce"`
}

func block_to_LE(block Block) Block_string {
	new_block := Block_string{}
	new_block.Hash = block.Hash
	new_block.Ver = BE_TO_LE(zfill(fmt.Sprintf("%x", block.Ver)))
	new_block.Prev_block = BE_TO_LE(block.Prev_block)
	new_block.Mrkl_root = BE_TO_LE(block.Mrkl_root)
	new_block.Bits = BE_TO_LE(fmt.Sprintf("%x", block.Bits))
	new_block.Time = BE_TO_LE(fmt.Sprintf("%x", block.Time))
	new_block.Nonce = BE_TO_LE(zfill(fmt.Sprintf("%x", block.Nonce)))
	return new_block
}

func unhexlify(str string) []byte {
	res := make([]byte, 0)
	for i := 0; i < len(str); i += 2 {
		x, _ := strconv.ParseInt(str[i:i+2], 16, 32)
		res = append(res, byte(x))
	}
	return res
}

func check(block Block_string) {
	data := block.Ver + block.Prev_block + block.Mrkl_root + block.Time + block.Bits + block.Nonce
	println(data)
	header_bin := unhexlify(data)
	hash := sha256.Sum256(header_bin)
	hash = sha256.Sum256(hash[:])
	var hashInt big.Int
	hashInt.SetBytes(hash[:])
	panic(fmt.Sprintf("$+v", hashInt))
	//if hash[:8] <0
	panic(fmt.Sprintf("%+v", Reverse_binary(hash)))

	res := hexlify(Reverse_binary(hash))

	if res == block.Hash {
		fmt.Println("good block")
	} else {
		fmt.Println("123")
		fmt.Printf("%+v", block.Hash)
		fmt.Println(" 123 ")
		fmt.Printf("%+v", res)
		panic("bad request")
	}
	print(res)
}

func main() {
	block := Block{}
	getJson("https://blockchain.info/rawblock/100000", &block)
	//d, _ := json.Marshal(block)
	//fmt.Printf("%+v", block)
	new := block_to_LE(block)
	fmt.Printf("%+v", new)
	check(new)
}
