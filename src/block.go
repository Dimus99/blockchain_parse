package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"time"
)

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, block *Block) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(fmt.Sprintf("No response from request, %+v", err))
	}
	//fmt.Printf("%s", resp.Body)
	err = json.NewDecoder(resp.Body).Decode(&block)
	if err != nil {
		panic(err)
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
	new_block.Bits = BE_TO_LE(zfill(fmt.Sprintf("%x", block.Bits)))
	new_block.Time = BE_TO_LE(zfill(fmt.Sprintf("%x", block.Time)))
	new_block.Nonce = BE_TO_LE(zfill(fmt.Sprintf("%x", block.Nonce)))
	return new_block
}
func hash256(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

func check(block Block_string) {
	data := block.Ver + block.Prev_block + block.Mrkl_root + block.Time + block.Bits + block.Nonce
	//println(data)
	head, _ := hex.DecodeString(data)
	hash := hash256(hash256(head))

	res := Reverse_binary(hash)
	h, _ := hex.DecodeString(block.Hash)
	if len(h) < 32 {
		panic(fmt.Sprintf("%+v\n", block.Hash))
	}
	hh := [32]byte(h)
	if res == hh {
		fmt.Println("good block")
	} else {
		fmt.Printf("%+v\n", block.Hash)
		fmt.Printf("%+v\n", res)
		panic("bad request")
	}
}

func main() {

	db, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS blocks_go(id INT PRIMARY KEY, block_index INTEGER, hash TEXT, version TEXT, merkle_root TEXT, bits TEXT, time TEXT, nonce TEXT, prev_block TEXT);")
	if err != nil {
		panic(err)
	}
	result, err := db.Query("SELECT max(block_index) FROM blocks_go;")
	start := START

	if err == nil {
		result.Next()
		result.Scan(&start)
	}
	quote_chan := grab(start)
	for {
		block_i := <-quote_chan
		block := block_i.block
		ind := block_i.index
		db.Exec("INSERT INTO blocks_go(block_index, hash, version, merkle_root, bits, time, nonce, prev_block) VALUES($1, $2, $3, $4, $5, $6, $7, $8);",
			ind, block.Hash, block.Ver, block.Mrkl_root, block.Bits, block.Time, block.Nonce, block.Prev_block)
		if ind%1000 == 0 {
			check(block)
			fmt.Println(ind)
		}

	}
}
