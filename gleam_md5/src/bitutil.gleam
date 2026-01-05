import gleam/bit_array
import gleam/int
import gleam/list
import gleam_community/maths
import gleam/float

pub type Bit512 {
  Bit512(array: BitArray)
}

pub fn pad_bits_to_modulu(bitarray: BitArray, modulo: Int) {
  let bit_size = bit_array.bit_size(bitarray)
  case bit_size % modulo {
    0 -> bitarray
    x -> {
      let pad_amount = modulo - x

      let pad = <<0:size(pad_amount)>>

      bit_array.append(bitarray, pad)
    }
  }
}

pub fn chunk_messages(bitarray: BitArray, size: Int) -> List(BitArray) {
  let padded = pad_bits_to_modulu(bitarray, size)
  chunk_bitarray(padded, size / 8)
}

fn chunk_bitarray(bitarray: BitArray, size: Int) -> List(BitArray) {
  let bitarray_size = bit_array.byte_size(bitarray)

  case bitarray_size {
    0 -> []
    x if x <= size -> [bitarray]
    _ ->
      case
        bit_array.slice(bitarray, 0, size),
        bit_array.slice(bitarray, size, bitarray_size - size)
      {
        Ok(first), Ok(rest) -> [first, ..chunk_bitarray(rest, size)]
        _, _ -> {
          []
        }
      }
  }
}

pub fn from_bitarray(bitarray: BitArray) {
  case bitarray |> bit_array.bit_size {
    512 -> Ok(Bit512(bitarray))
    x -> Error("The size is " <> x |> int.to_string <> " instead of 512")
  }
}

pub fn to_intlist(bits: Bit512) {
  split_32bit_ints(bits.array, [])
}

fn split_32bit_ints(bits: BitArray, acc: List(Int)) -> List(Int) {
  case bits {
    <<header:little-size(32), rest:bits>> ->
      split_32bit_ints(rest, [header, ..acc])
    _ -> list.reverse(acc)
  }
}
const bitmask = 0xFFFFFFFF
pub fn u32(x: Int) -> Int {
  int.bitwise_and(x, bitmask)
}

pub fn leftrotate_32bit(x, c) {
  let x = int.bitwise_and(x, bitmask)

  let left = int.bitwise_shift_left(x, c)
  let right =
    int.bitwise_shift_right(x, 32 - c)
    |> int.bitwise_and(bitmask)

  int.bitwise_or(left, right)
  |> int.bitwise_and(bitmask)
}

pub fn add32(a: Int, b: Int) -> Int {
  u32(a + b)
}

