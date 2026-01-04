package bitutil

import "errors"

type BitArray struct {
	data []byte
	bits uint64
}

func (b *BitArray) Length() uint64 {
	return b.bits
}

func (b *BitArray) AppendBit(value bool) {
	byteIndex := b.bits / 8
	bitIndex := b.bits % 8

	if byteIndex >= uint64(len(b.data)) {
		b.data = append(b.data, 0)
	}

	if value {
		b.data[byteIndex] |= 1 << bitIndex
	}

	b.bits++
}

func (b *BitArray) AppendByte(v byte) {
	for i := 7; i >= 0; i-- {
		b.AppendBit((v>>i)&1 == 1)
	}
}

func (b *BitArray) GetBit(index uint64) (bool, error) {
	if index >= b.bits {
		return false, errors.New("bit index out of range")
	}

	byteIndex := index / 8
	bitIndex := index % 8

	return (b.data[byteIndex] & (1 << bitIndex)) != 0, nil
}

func (b *BitArray) ToUint32Array() ([]uint32, error) {
	if b.bits%32 != 0 {
		return nil, errors.New("bit length not divisible by 32")
	}

	count := b.bits / 32
	result := make([]uint32, count)

	for i := uint64(0); i < count; i++ {
		var value uint32
		for bit := uint64(0); bit < 32; bit++ {
			v, _ := b.GetBit(i*32 + bit)
			if v {
				value |= 1 << bit
			}
		}
		result[i] = value
	}

	return result, nil
}

func (b *BitArray) Split(size uint64) ([]*BitArray, error) {
	if size == 0 {
		return nil, errors.New("size must be > 0")
	}

	var result []*BitArray

	for i := uint64(0); i < b.bits; i += size {
		part := &BitArray{}
		end := i + size
		if end > b.bits {
			end = b.bits
		}

		for j := i; j < end; j++ {
			v, _ := b.GetBit(j)
			part.AppendBit(v)
		}

		result = append(result, part)
	}

	return result, nil
}

func (b *BitArray) Clone() BitArray {
	dataCopy := make([]byte, len(b.data))
	copy(dataCopy, b.data)

	return BitArray{
		data: dataCopy,
		bits: b.bits,
	}
}
