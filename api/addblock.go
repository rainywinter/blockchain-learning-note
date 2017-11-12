package api

import (
	"fmt"
	"net/http"
	"rw/blockchain/blockchain"
	"time"
)

const (
	data_key = "data"
	wait_msg = "create new block, please wait a while."
)

func init() {
	http.HandleFunc("/addblock", addBlock)
}

func addBlock(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	data := r.FormValue(data_key)
	if data == "" {
		w.Write([]byte("error: empty parameter."))
		return
	}

	// done := make(chan struct{})
	// defer close(done)
	// go func() {
	// FOR_MARK:
	// 	for {
	// 		select {
	// 		case <-done:
	// 			break FOR_MARK
	// 		default:
	// 			w.Write([]byte(wait_msg))
	// 			time.Sleep(time.Second)
	// 		}
	// 	}
	// }()

	go func() {
		fmt.Printf("begin create new block,data:%s\n", data)
		begin := time.Now()
		chain := blockchain.GetBlockchain()
		block := chain.AddBlock(data)

		str := fmt.Sprintf("prev hash:%x\nblock data:%s\nblock hash:%x\ncost time:%.2fs.\n", block.PrevBlockHash, block.Data, block.Hash, time.Now().Sub(begin).Seconds())
		fmt.Println(str)
	}()

	w.Write([]byte(wait_msg))
}
