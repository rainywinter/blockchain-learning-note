package blockchain

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	// 为计算合适的hash值
	Nonce int
}

// 替代为创建新区块时查找合适的哈希值
// func (b *Block) SetHash() {
// 	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
// 	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
// 	hash := sha256.Sum256(headers)

// 	b.Hash = hash[:]
// }

func NewBlock(data string, prevBlockHash []byte) *Block {
	defer func(begin time.Time) {
		fmt.Printf("Create block,data: %s, cost time: %.2f s.\n", data, time.Now().Sub(begin).Seconds())
	}(time.Now())
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash
	block.Nonce = nonce

	return block
}

//Serialize  将区块信息序列化为字节码
func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		fmt.Println("Serialize error ", err)
		return nil
	}
	return result.Bytes()
}

// DeserializeBlock 反序列化，解出区块信息
func DeserializeBlock(src []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(src))
	err := decoder.Decode(&block)
	if err != nil {
		fmt.Println("Serialize error ", err)
		return nil
	}

	return &block
}
