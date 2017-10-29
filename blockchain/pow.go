package blockchain

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"strconv"
)

// 难度系数
const targetBits = 24

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}

	return pow
}

func (p *ProofOfWork) getHashData(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			p.block.PrevBlockHash,
			p.block.Data,
			[]byte(strconv.FormatInt(p.block.Timestamp, 10)),
			[]byte(strconv.FormatInt(targetBits, 10)),
			[]byte(strconv.FormatInt(int64(nonce), 10)),
		},
		[]byte{},
	)
}

// 暴力查找满足条件的nonce和hash
func (p *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for nonce < math.MaxInt64 {
		data := p.getHashData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(p.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	return nonce, hash[:]
}

func (p *ProofOfWork) Valid() bool {
	var hashInt big.Int

	data := p.getHashData(p.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(p.target) == -1
}
