import bitutil
import gleam/bit_array
import gleam/float
import gleam/int
import gleam/list
import gleam/result
import gleam/string
import gleam_community/maths

const bitmask = 0xFFFFFFFF

fn u32(x: Int) -> Int {
  int.bitwise_and(x, bitmask)
}

fn leftrotate_32bit(x, c) {
  let x = int.bitwise_and(x, bitmask)

  let left = int.bitwise_shift_left(x, c)
  let right =
    int.bitwise_shift_right(x, 32 - c)
    |> int.bitwise_and(bitmask)

  int.bitwise_or(left, right)
  |> int.bitwise_and(bitmask)
}

fn add32(a: Int, b: Int) -> Int {
  u32(a + b)
}

pub fn compute_table() {
  let two32 = 4_294_967_296.0
  list.range(0, 63)
  |> list.map(fn(x) {
    maths.sin(int.to_float(x + 1))
    |> float.absolute_value
    |> fn(v) { two32 *. v }
    |> float.floor()
    |> float.truncate()
  })
}

fn preprocess_bits_and_chunk(bits: BitArray) {
  let original_length = bit_array.bit_size(bits)
  let bits = bit_array.append(bits, <<1:1>>)
  let length_mod_512 = bit_array.bit_size(bits) % 512
  let pad_bits = { { 448 - length_mod_512 } + 512 } % 512
  let bits = bit_array.append(bits, <<0:size(pad_bits)>>)
  let length_padding = <<original_length:little-size(64)>>
  let bits = bit_array.append(bits, length_padding)
  let blocks =
    bitutil.chunk_messages(bits, 512)
    |> list.map(bitutil.from_bitarray)
    |> list.map(fn(x) {
      case x {
        Ok(val) -> val
        _ -> panic as "Could not split bit array (this should not happen)"
      }
    })
}

pub fn md5(bits: BitArray) -> String {
  let blocks = preprocess_bits_and_chunk(bits)

  let ktable = compute_table()

  // Initial values defined in RFC 1321
  let a0 = 0x67452301
  let b0 = 0xEFCDAB89
  let c0 = 0x98BADCFE
  let d0 = 0x10325476

  let #(a, b, c, d) =
    blocks
    |> list.fold(#(a0, b0, c0, d0), fn(state, block) {
      let #(a, b, c, d) = state
      step(ktable, block, a, b, c, d)
    })

  <<a:little-size(32), b:little-size(32), c:little-size(32), d:little-size(32)>>
  |> bit_array.base16_encode
  |> string.lowercase
}

pub fn s(i: Int) -> Int {
  case i {
    i if i >= 0 && i < 16 ->
      case i % 4 {
        0 -> 7
        1 -> 12
        2 -> 17
        _ -> 22
      }
    i if i >= 16 && i < 32 ->
      case i % 4 {
        0 -> 5
        1 -> 9
        2 -> 14
        _ -> 20
      }
    i if i >= 32 && i < 48 ->
      case i % 4 {
        0 -> 4
        1 -> 11
        2 -> 16
        _ -> 23
      }
    i if i >= 48 && i < 64 ->
      case i % 4 {
        0 -> 6
        1 -> 10
        2 -> 15
        _ -> 21
      }
    _ -> panic("s(i): index out of range (0..63)")
  }
}

fn inner_step(i, k_i, messages, a, b, c, d) {
  let #(f, g) = case i {
    _ if i <= 15 -> {
      let f =
        int.bitwise_and(b, c)
        |> int.bitwise_or(int.bitwise_not(b) |> int.bitwise_and(d))
      let g = i
      #(f, g)
    }

    _ if i <= 31 -> {
      let f =
        int.bitwise_and(d, b)
        |> int.bitwise_or(int.bitwise_not(d) |> int.bitwise_and(c))
      let g = { 5 * i + 1 } % 16
      #(f, g)
    }

    _ if i <= 47 -> {
      let f = int.bitwise_exclusive_or(b, c) |> int.bitwise_exclusive_or(d)
      let g = { 3 * i + 5 } % 16
      #(f, g)
    }

    _ -> {
      let f =
        int.bitwise_not(d) |> int.bitwise_or(b) |> int.bitwise_exclusive_or(c)
      let g = { 7 * i } % 16
      #(f, g)
    }
  }

  let f = u32(f)
  let g = u32(g)

  let temp = d
  let d = c
  let c = b

  let m_g = get(messages, g) |> result.unwrap(0)
  let sum = add32(add32(add32(a, f), k_i), m_g)
  let rotated = leftrotate_32bit(sum, s(i))
  let b = add32(b, rotated)

  let a = temp
  #(a, b, c, d)
}

//Gleam uses Linked lists for arrays
//Gleam usually doesnt support indexed lookups, because the are O(n)
//external libraries like gleam_array that use plattform specific implementations could be used, to speed this up
fn get(arr, index) {
  case index, arr {
    0, [el, ..] -> Ok(el)
    _, [_, ..rest] -> get(rest, index - 1)
    _, _ -> Error("Index out of bound")
  }
}

pub fn set(messages: List(a), index: Int, value: a) -> List(a) {
  case messages, index {
    _, i if i < 0 -> messages
    [], 0 -> [value]
    [], _ -> []
    [_, ..rest], 0 -> [value, ..rest]
    [head, ..rest], i -> [head, ..set(rest, i - 1, value)]
  }
}

fn step(ktable: List(Int), messageblock: bitutil.Bit512, a, b, c, d) {
  let message_chunks = messageblock |> bitutil.to_intlist()
  let zipped_list = list.range(0, 63) |> list.zip(ktable)

  let #(aa, bb, cc, dd) =
    zipped_list
    |> list.fold(#(a, b, c, d), fn(acc, pair) {
      let #(i, k_i) = pair
      let #(a, b, c, d) = acc
      inner_step(i, k_i, message_chunks, a, b, c, d)
    })

  #(add32(a, aa), add32(b, bb), add32(c, cc), add32(d, dd))
}
