package blockchain

import (
	"sync"

	"github.com/boltdb/bolt"
)

// 持久化存储区块数据
const DBFile = "blockchain.db"
const BlocksBucket = "blocks"

var (
	blockchainInstance *Blockchain
	// 仅仅创建一个链,避免并发创建多个
	once sync.Once
)

type BlockchainIterator struct {
	curHash []byte
	db      *bolt.DB
}

// 逆向迭代区块数据
func (bci *BlockchainIterator) Next() *Block {
	var block *Block

	err := bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))

		buf := b.Get(bci.curHash)
		block = DeserializeBlock(buf)
		return nil
	})

	if err != nil {
		panic(err)
	}
	bci.curHash = block.PrevBlockHash

	return block
}

func (bci *BlockchainIterator) End() bool {
	return bci.curHash == nil
}

type Blockchain struct {
	// Blocks []*Block
	// 尾部hash
	tip []byte
	db  *bolt.DB
}

func (bc *Blockchain) Iteratore() *BlockchainIterator {
	return &BlockchainIterator{
		curHash: bc.tip,
		db:      bc.db,
	}
}

func (bc *Blockchain) AddBlock(data string) *Block {
	// prevBlock := bc.Blocks[len(bc.Blocks)-1]
	// newBlock := NewBlock(data, prevBlock.Hash)
	// bc.Blocks = append(bc.Blocks, newBlock)

	// 存储文件
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		panic(err)
	}

	// 创建新块
	newBlock := NewBlock(data, lastHash)

	// 将新块持久化
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		err = b.Put(newBlock.Hash, newBlock.Serialize())
		err = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash
		return nil
	})

	return newBlock
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func newBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(DBFile, 0600, nil)
	if err != nil {
		panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		if b == nil {
			genesis := NewGenesisBlock()

			b, err = tx.CreateBucket([]byte(BlocksBucket))
			if err != nil {
				panic(err)
			}

			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				panic(err)
			}

			err = b.Put([]byte("l"), genesis.Hash)
			if err != nil {
				panic(err)
			}

			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	// return &Blockchain{
	// 	Blocks: []*Block{NewGenesisBlock()},
	// }

	return &Blockchain{
		tip: tip,
		db:  db,
	}
}

func GetBlockchain() *Blockchain {
	if blockchainInstance == nil {
		once.Do(func() {
			blockchainInstance = newBlockchain()
		})
	}
	return blockchainInstance
}
