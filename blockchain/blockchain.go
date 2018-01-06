package blockchain

import (
	"encoding/hex"
	"log"
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

func (bc *Blockchain) Iterator() *BlockchainIterator {
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

func (bc *Blockchain) FindUnspentTransactions(address string) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIDx, out := range tx.Vout {
				// 交易输出被花费了
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIDx {
							continue Outputs
						}
					}
				}
				// 可以被解锁
				if out.CanBeUnlockedWith(address) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}

		}

		if bci.End() {
			break
		}
	}
	return unspentTXs
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(address)
	accumulated := 0

Mark:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				if accumulated >= amount {
					break Mark
				}
			}
		}
	}
	return accumulated, unspentOutputs
}

func (bc *Blockchain) MineBlock(transactions []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})
}

// 创建创世区块以初始化区块链
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
