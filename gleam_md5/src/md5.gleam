import algorithm
import argv
import gleam/bit_array
import gleam/io
import simplifile

pub fn main() {
  let args = argv.load().arguments
  case args {
    ["-f", file] -> {
      let file = simplifile.read_bits(file)
      case file {
        Ok(bits) -> {
          io.println(algorithm.md5(bits))
        }
        Error(error) -> {
          io.println(
            "Could not read file: " <> simplifile.describe_error(error),
          )
        }
      }
    }
    ["-s", string] -> io.println(algorithm.md5(bit_array.from_string(string)))
    _ ->
      io.println(
        "Usages:\n  gleam run -- -f <file>\n  gleam run -- -s <string>",
      )
  }
}
