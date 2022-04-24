package main

import (
	"crypto/sha256"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block *Block
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork(b *Block) *ProofOfWork {
	pow := &ProofOfWork{b}

	return pow
}

// Run performs a proof-of-work
// implement
func (pow *ProofOfWork) Run() (int, []byte) {
	nonce := 0
	hashInt := new(big.Int)
	target := big.NewInt(1)
	target.Lsh(target, 256-pow.block.Bits)

	temp := []byte{}
	temp = append(temp, pow.block.PrevBlockHash...)
	temp = append(temp, pow.block.HashData()...)
	temp = append(temp, IntToHex(pow.block.Timestamp)...)
	temp = append(temp, IntToHex(int64(pow.block.Bits))...)

	temp1 := append(temp, IntToHex(int64(nonce))...)
	temp2 := sha256.Sum256(temp1)
	hashInt.SetBytes(temp2[:])
	for target.Cmp(hashInt) <= 0 {
		nonce++
		temp1 = append(temp, IntToHex(int64(nonce))...)
		temp2 = sha256.Sum256(temp1)
		hashInt.SetBytes(temp2[:])
	}
	nonce--
	return nonce, temp2[:]
}

// Validate validates block's PoW
// implement
func (pow *ProofOfWork) Validate() bool {
	target := big.NewInt(1)
	target.Lsh(target, 256-pow.block.Bits)
	hashInt := new(big.Int)
	hashInt.SetBytes(pow.block.Hash)
	if hashInt.Cmp(target) < 0 {
		return true
	} else {
		return false
	}
}
