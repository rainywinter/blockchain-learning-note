package api

import (
	"fmt"
	"net/http"
	"rw/blockchain/blockchain"
	"strconv"
)

func init() {
	http.HandleFunc("/dumpchain", dumpChain)
}

func dumpChain(w http.ResponseWriter, r *http.Request) {
	chain := blockchain.GetBlockchain()

	iter := chain.Iterator()

	for {
		block := iter.Next()
		w.Write([]byte(fmt.Sprintf("PrevHash: %x\n", block.PrevBlockHash)))
		w.Write([]byte(fmt.Sprintf("Data: %s\n", block.Data)))
		w.Write([]byte(fmt.Sprintf("Hash: %x\n", block.Hash)))
		pow := blockchain.NewProofOfWork(block)
		w.Write([]byte(fmt.Sprintf("PoW: %s\n", strconv.FormatBool(pow.Valid()))))
		w.Write([]byte("\n\n"))

		if iter.End() {
			break
		}
	}
}
