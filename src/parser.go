package main

import (
	"fmt"
)

var (
	WORKERS int = 8
	START       = 407401
)

type BLOCK_INDEX struct {
	index int
	block Block_string
}

func grab(start int) <-chan BLOCK_INDEX {
	c := make(chan BLOCK_INDEX)
	for i := 0; i < WORKERS; i++ {
		go func(num int) {
			ii := num
			for j := start + ii; j < 800000; j += WORKERS {
				block := Block{}
				getJson(fmt.Sprintf("https://blockchain.info/rawblock/%+v", j), &block)
				new := block_to_LE(block)
				if len(block.Hash) == 0 {
					panic(fmt.Sprintf("%+v\n, %+v\n", block, j))
				}
				c <- BLOCK_INDEX{j, new}
			}
		}(i)
	}
	fmt.Println("Запущено потоков: ", WORKERS)
	return c
}
