package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"time"
)

type Block struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
	// 为计算合适的hash值
	Nonce int
	// 交易
	Transactions []*Transaction
}

// 替代为创建新区块时查找合适的哈希值
// func (b *Block) SetHash() {
// 	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
// 	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, timestamp}, []byte{})
// 	hash := sha256.Sum256(headers)

// 	b.Hash = hash[:]
// }

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	defer func(begin time.Time) {
		fmt.Printf("Create block, cost time: %.2f s.\n", time.Now().Sub(begin).Seconds())
	}(time.Now())

	block := &Block{
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
		Transactions:  transactions,
	}

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

// 计算区块中所有交易的hash
func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
