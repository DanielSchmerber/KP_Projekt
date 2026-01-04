package bitutil

import "testing"

func toBitArray(s string) *BitArray {
	return NewBitArrayFromBytes([]byte(s))
}

func TestNewBitArrayFromBytes(t *testing.T) {

	testArray := toBitArray("Hello world!")

	if testArray.Length() != 12*8 {
		t.Error("Wrong length")
	}
	t.Log("Test finished")
}

func TestBitArray_Split(t *testing.T) {

	testArray := toBitArray("Hello world!")

	split, err := testArray.Split(8)

	if err != nil {
		t.Error(err)
	}

	if len(split) != 12 {
		t.Errorf("Wrong amount of split arrays! got %d expected 12", len(split))
	}

	if split[0].Length() != 8 {
		t.Errorf("Wrong length of split array! got %d expected 8", split[0].Length())
	}

}

func TestBitArray_ToUint32Array(t *testing.T) {
	testArray := toBitArray("Hello world!")

	intArray, err := testArray.ToUint32Array()

	if err != nil {
		t.Error(err)
	}

	if len(intArray) != 3 {
		t.Error("Wrong amount of integers array")
	}

}
