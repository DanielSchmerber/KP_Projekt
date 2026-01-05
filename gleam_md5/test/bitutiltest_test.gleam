import bitutil
import gleam/bit_array
import gleam/list
import gleam/result
import gleam/string

pub fn split_test() {
  let split = bitutil.chunk_messages(<<"Hello World">>, 8)

  assert list.length(split) == { "Hello World" |> string.length }

  case split {
    [first, ..] -> {
      assert { first |> bit_array.to_string } == Ok("H")
    }
    _ -> panic
  }
}

pub fn to_intarray_test() {
  let intarray =
    <<"Hello World!">>
    |> bitutil.pad_bits_to_modulu(512)
    |> bitutil.from_bitarray
    |> result.map(bitutil.to_intlist)

  case intarray {
    Ok(val) -> {
      assert list.length(val) == 512 / 32
    }
    _ -> {
      panic as "could not convert to int list"
    }
  }
}

pub fn u32_basic_test() {
  assert bitutil.u32(0xFFFFFFFF) == 0xFFFFFFFF
}

pub fn u32_truncates_overflow_test() {
  assert bitutil.u32(0x1FFFFFFFF) == 0xFFFFFFFF
}

pub fn u32_negative_test() {
  assert bitutil.u32(-1) == 0xFFFFFFFF
}

pub fn leftrotate_simple_test() {
  assert bitutil.leftrotate_32bit(1, 1) == 2
}

pub fn leftrotate_wraparound_test() {
  assert bitutil.leftrotate_32bit(0x80000000, 1) == 1
}

pub fn leftrotate_zero_test() {
  assert bitutil.leftrotate_32bit(0x12345678, 0) == 0x12345678
}

pub fn leftrotate_full_cycle_test() {
  assert bitutil.leftrotate_32bit(0xDEADBEEF, 32) == 0xDEADBEEF
}

pub fn add32_no_overflow_test() {
  assert bitutil.add32(1, 2) == 3
}

pub fn add32_overflow_test() {
  assert bitutil.add32(0xFFFFFFFF, 1) == 0
}

pub fn add32_negative_test() {
  assert bitutil.add32(-1, 1) == 0
}
