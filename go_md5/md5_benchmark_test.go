package main

import (
	"bytes"
	"go_md5/bitutil"
	"testing"
)

var md5BenchInputs = []struct {
	name string
	data []byte
}{
	{"Empty", []byte("")},
	{"Small", []byte("test")},
	{"Block-64B", bytes.Repeat([]byte("a"), 64)},
	{"Medium-1KB", bytes.Repeat([]byte("a"), 1024)},
	{"Large-1MB", bytes.Repeat([]byte("a"), 1024*1024)},
}

func BenchmarkMd5(b *testing.B) {
	for _, input := range md5BenchInputs {
		b.Run(input.name, func(b *testing.B) {
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				bits := bitutil.NewBitArrayFromBytes(input.data)
				_, err := Md5(bits)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
