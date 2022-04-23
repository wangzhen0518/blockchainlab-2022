package main

import (
	"bytes"
	"math"
)

// MerkleTree represent a Merkle tree
type MerkleTree struct {
	RootNode *MerkleNode
}

// MerkleNode represent a Merkle tree node
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// NewMerkleTree creates a new Merkle tree from a sequence of data
// implement
func NewMerkleTree(data [][]byte) *MerkleTree {
	num := len(data)
	depth := math.Ceil(math.Log2(float64(num)))
	node, _ := dfs(data, int(depth), 0, num)
	var mTree = MerkleTree{node}

	return &mTree
}

func dfs(data [][]byte, depth, num, numLimit int) (*MerkleNode, int) {
	if depth == 0 {
		if num < numLimit {
			return NewMerkleNode(nil, nil, data[num]), num + 1
		}
		return nil, num
	}
	left, numLeft := dfs(data, depth-1, num, numLimit)
	right, numRight := dfs(data, depth-1, numLeft, numLimit)
	node, numTol := NewMerkleNode(left, right, nil), numRight
	return node, numTol
}

// NewMerkleNode creates a new Merkle tree node
// implement
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}
	node.Left = left
	node.Right = right
	temp := [32]byte{}
	if left == nil {
		temp = mySha256(data)
	} else {
		if right == nil { //奇数情况
			right = &MerkleNode{nil, nil, left.Data}
		}
		temp = mySha256(bytes.Join([][]byte{left.Data, right.Data}, []byte{}))
	}
	node.Data = make([]byte, 32)
	copy(node.Data, temp[:])
	return &node
}
