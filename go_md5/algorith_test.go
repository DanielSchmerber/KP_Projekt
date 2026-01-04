package main

import (
	"go_md5/bitutil"
	"testing"
)

func toBitArray(s string) *bitutil.BitArray {
	return bitutil.NewBitArrayFromBytes([]byte(s))
}

func TestMd5(t *testing.T) {
	testHash(t, "The quick brown fox jumps over the lazy dog", "9e107d9d372bb6826bd81d3542a419d6")
	testHash(t, "The quick brown fox jumps over the lazy dog.", "e4d909c290d0fb1ca068ffaddf22cbd0")
	testHash(t, "", "d41d8cd98f00b204e9800998ecf8427e")

}

func testHash(t *testing.T, input string, expected string) {
	t.Helper()

	result, err := Md5(toBitArray(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != expected {
		t.Fatalf(
			"wrong hash for input %q: got %q, expected %q",
			input, result, expected,
		)
	}
}
