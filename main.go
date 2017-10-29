package main

import (
	"fmt"
	"rw/blockchain/blockchain"
	"strconv"
	"time"
)

func main() {
	defer func(begin time.Time) {
		fmt.Printf("total cost time: %.2f s\n", time.Now().Sub(begin).Seconds())
	}(time.Now())

	chain := blockchain.NewBlockchain()
	chain.AddBlock("test1")
	chain.AddBlock("test2")

	for _, block := range chain.Blocks {
		fmt.Printf("PrevHash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Valid()))
		fmt.Println()
	}
}
