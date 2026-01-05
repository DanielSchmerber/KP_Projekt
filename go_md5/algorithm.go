package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"go_md5/bitutil"
	"math"
)

func generateKtable() []uint32 {

	K := make([]uint32, 64)
	for i := 0; i < 64; i++ {
		val := math.Floor(math.Pow(2, 32) * math.Abs(math.Sin(float64(i)+1)))
		K[i] = uint32(val)
	}

	return K

}

var s = [64]uint32{
	7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22,
	5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20,
	4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23,
	6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21,
}

func Md5(array *bitutil.BitArray) (string, error) {

	messageChunks, err := preprocessBytes(array).Split(512)
	kTable := generateKtable()
	if err != nil {
		return "", errors.New("could not split message in 512 bit chunks. (this shouldnt happen)")
	}

	// Initialisiere die Variablen: (laut RFC 1321)
	var a0 uint32 = 0x67452301
	var b0 uint32 = 0xEFCDAB89
	var c0 uint32 = 0x98BADCFE
	var d0 uint32 = 0x10325476

	for _, chunk := range messageChunks {

		submessages, err := chunk.ToUint32Array()

		if err != nil {
			return "", errors.New("could not convert message chunk to uint32 array: " + err.Error())
		}

		a := a0
		b := b0
		c := c0
		d := d0

		for i := range 64 {
			f, g := uint32(0), uint32(0)
			if i >= 0 && i <= 15 {
				f = (b & c) | (^b & d)
				g = uint32(i)
			}

			if i >= 16 && i <= 31 {
				f = (d & b) | (^d & c)
				g = uint32((5*i + 1) % 16)
			}

			if i >= 32 && i <= 47 {
				f = b ^ c ^ d
				g = uint32((3*i + 5) % 16)
			}

			if i >= 48 && i <= 63 {
				f = c ^ (b | ^d)
				g = uint32((7 * i) % 16)
			}
			f = f + a + kTable[i] + submessages[g]
			a = d
			d = c
			c = b
			b = b + bitutil.LeftRotate(f, s[i])
		}

		a0 += a
		b0 += b
		c0 += c
		d0 += d

	}

	digest := make([]byte, 16)

	binary.LittleEndian.PutUint32(digest[0:4], a0)
	binary.LittleEndian.PutUint32(digest[4:8], b0)
	binary.LittleEndian.PutUint32(digest[8:12], c0)
	binary.LittleEndian.PutUint32(digest[12:16], d0)

	return hex.EncodeToString(digest), nil
}

func preprocessBytes(bitArray *bitutil.BitArray) *bitutil.BitArray {
	ba := bitArray.Clone()

	origLen := ba.Length()

	ba.AppendBit(true)

	curLen := ba.Length()
	padZeros := (448 - (curLen % 512) + 512) % 512
	for i := uint64(0); i < padZeros; i++ {
		ba.AppendBit(false)
	}

	lenBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(lenBuf, origLen)

	for _, v := range lenBuf {
		ba.AppendByte(v)
	}

	return &ba
}
