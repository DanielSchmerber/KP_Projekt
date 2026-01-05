package bitutil

import (
	"errors"
	"fmt"
)

type BitArray struct {
	data []bool
}

func NewBitArrayFromBytes(bytes []byte) *BitArray {
	b := &BitArray{
		data: make([]bool, 0, len(bytes)*8),
	}

	for _, v := range bytes {
		b.AppendByte(v)
	}

	return b
}

func (b *BitArray) AppendByte(v byte) {
	for i := 7; i >= 0; i-- {
		b.AppendBit((v>>i)&1 == 1)
	}
}

func (b *BitArray) Length() uint64 {
	return uint64(len(b.data))
}

func (b *BitArray) AppendBit(value bool) {
	b.data = append(b.data, value)
}

func (b *BitArray) GetBit(index uint64) (bool, error) {
	if index >= uint64(len(b.data)) {
		return false, errors.New("bit index out of range")
	}
	return b.data[index], nil
}

func (b *BitArray) ToUint32Array() ([]uint32, error) {
	if len(b.data)%32 != 0 {
		return nil, fmt.Errorf("bit length not divisible by 32: bits=%d", len(b.data))
	}

	count := len(b.data) / 32
	result := make([]uint32, count)

	readByte := func(bitOffset int) byte {
		var v byte
		for i := 0; i < 8; i++ {
			if b.data[bitOffset+i] {
				v |= 1 << (7 - i)
			}
		}
		return v
	}

	for i := 0; i < count; i++ {
		base := i * 32
		b0 := uint32(readByte(base + 0))
		b1 := uint32(readByte(base + 8))
		b2 := uint32(readByte(base + 16))
		b3 := uint32(readByte(base + 24))

		result[i] = b0 | (b1 << 8) | (b2 << 16) | (b3 << 24)
	}

	return result, nil
}

func (b *BitArray) Split(size uint64) ([]*BitArray, error) {
	if size == 0 {
		return nil, errors.New("size must be > 0")
	}

	var result []*BitArray

	for i := uint64(0); i < uint64(len(b.data)); i += size {
		end := i + size
		if end > uint64(len(b.data)) {
			end = uint64(len(b.data))
		}

		part := &BitArray{
			data: append([]bool(nil), b.data[i:end]...),
		}
		result = append(result, part)
	}

	return result, nil
}

func (b *BitArray) Clone() BitArray {
	copyData := append([]bool(nil), b.data...)
	return BitArray{data: copyData}
}
func LeftRotate(x uint32, c uint32) uint32 {
	return (x << c) | (x >> (32 - c))
}
