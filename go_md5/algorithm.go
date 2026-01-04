package main

import (
	"encoding/binary"
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

func md5(array bitutil.BitArray) (string, error) {

	messageChunks, err := array.Split(512)

	if err != nil {
		return "", errors.New("could not split message in 512 bit chunks. (this shouldnt happen)")
	}

	// Initialisiere die Variablen: (laut RFC 1321)
	var a0 uint32 = 0x67452301
	var b0 uint32 = 0xEFCDAB89
	var c0 uint32 = 0x98BADCFE
	var d0 uint32 = 0x10325476

	for _, chunk := range messageChunks {

		submessages := chunk.ToUint32Array()

	}
}

func preprocessBytes(bitArray *bitutil.BitArray) *bitutil.BitArray {

	byteArray := bitArray.Clone()

	//Append a 1 at the end
	byteArray.AppendBit(true)

	messageLength := byteArray.Length()
	messageLengthbuf := make([]byte, 4)
	binary.LittleEndian.PutUint64(messageLengthbuf, messageLength)

	var requiredPaddingBits = (448 - (messageLength % 512)) % 512
	requiredPaddingBits = (requiredPaddingBits + 512) % 512

	requiredPaddingBits -= 1

	for i := uint64(0); i < requiredPaddingBits; i++ {
		byteArray.AppendBit(false)
	}

	for _, v := range messageLengthbuf {
		byteArray.AppendByte(v)
	}
	return &byteArray
}
