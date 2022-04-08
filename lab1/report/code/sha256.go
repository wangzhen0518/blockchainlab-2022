package main

import (
	"bytes"
	"encoding/binary"
)

func mySha256(message []byte) [32]byte {
	//前八个素数平方根的小数部分的前面32位
	// h0 := uint32(0x6a09e667)
	// h1 := uint32(0xbb67ae85)
	// h2 := uint32(0x3c6ef372)
	// h3 := uint32(0xa54ff53a)
	// h4 := uint32(0x510e527f)
	// h5 := uint32(0x9b05688c)
	// h6 := uint32(0x1f83d9ab)
	// h7 := uint32(0x5be0cd19)
	h := [8]uint32{
		0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a, 0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19}

	//自然数中前面64个素数的立方根的小数部分的前32位
	k := [64]uint32{
		0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5, 0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
		0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3, 0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
		0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc, 0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
		0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7, 0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
		0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13, 0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
		0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3, 0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
		0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5, 0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
		0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208, 0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2}

	tolMessage, N := fillMessage(message)
	N /= 64
	div512Message := [][]byte{}
	for i := 0; i < N; i++ {
		div512Message = append(div512Message, tolMessage[i*64:i*64+64])
	}

	for i := 0; i < N; i++ {
		div32Message := seperate32Message(div512Message[i])

		loopVar := h
		// a-0, b-1, c-2, d-3, e-4, f-5, g-6, h-7
		for p := 0; p < 64; p++ {
			t1 := loopVar[7] + _SIGMA1(loopVar[4]) + _Ch(loopVar[4], loopVar[5], loopVar[6]) + k[p] + div32Message[p]
			t2 := _SIGMA0(loopVar[0]) + _Maj(loopVar[0], loopVar[1], loopVar[2])
			for q := 7; q > 0; q-- {
				loopVar[q] = loopVar[q-1]
			}
			loopVar[4] += t1
			loopVar[0] = t1 + t2
		}
		for p := 0; p < 8; p++ {
			h[p] += loopVar[p]
		}
	}

	tempDate := []byte{}
	sha256data := [32]byte{}
	for i := 0; i < 8; i++ {
		tempDate = append(tempDate, uint32ToBytes(h[i])...)
	}
	copy(sha256data[:], tempDate)
	return sha256data
}

func int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}
func uint32ToBytes(i uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}

func fillMessage(message []byte) (plus_message []byte, num8 int) {
	by0 := []byte{0b00000000}
	by1 := []byte{0b10000000}
	// 需要补k个0，满足 8*length + 1 + k = 512m + 448
	// 从而 k + 1 + 8*length mod 512 =448
	// 即 k = 448 - 8*length mod 512 - 1
	// 所以8位的1Byte 0有 (448 - 8*length mod 512 - 1)/8 个
	// 额外有一列 '1' + '0'*(63 - length mod 64) mod 8
	length := len(message)
	k := (448 - (8*length)%512 - 1) / 8
	app := bytes.Repeat(by0, k)
	temp := int64ToBytes(int64(8 * length))
	plus_message = bytes.Join([][]byte{message, by1, app, temp}, []byte{})
	num8 = len(plus_message)
	return
}

// seperate32Message 将数据分为32位一组的包，使用uint32是为了后续的位运算，否则继续使用byte，位运算结果将出错
func seperate32Message(div512Message []byte) (div32Message []uint32) {
	for i := 0; i < 16; i++ {
		div32Message = append(div32Message, binary.BigEndian.Uint32(div512Message[i*4:i*4+4]))
	}
	for i := 16; i < 64; i++ {
		div32Message = append(div32Message, _sigma1(div32Message[i-2])+div32Message[i-7]+
			_sigma0(div32Message[i-15])+div32Message[i-16])
	}
	return
}

func rightRotate(n uint32, d uint) uint32 {
	return ((n >> d) | (n << (32 - d))) & ((1 << 32) - 1)
}
func _Ch(x, y, z uint32) uint32 {
	return (x & y) ^ ((^x) & z)
}
func _Maj(x, y, z uint32) uint32 {
	return (x & y) ^ (x & z) ^ (y & z)
}
func _SIGMA0(n uint32) uint32 {
	return rightRotate(n, 2) ^ rightRotate(n, 13) ^ rightRotate(n, 22)
}
func _SIGMA1(n uint32) uint32 {
	return rightRotate(n, 6) ^ rightRotate(n, 11) ^ rightRotate(n, 25)
}
func _sigma0(n uint32) uint32 {
	return rightRotate(n, 7) ^ rightRotate(n, 18) ^ (n >> 3)
}
func _sigma1(n uint32) uint32 {
	return rightRotate(n, 17) ^ rightRotate(n, 19) ^ (n >> 10)
}
