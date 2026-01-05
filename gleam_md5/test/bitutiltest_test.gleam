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
