import algorithm
import argv
import gleam/bit_array
import gleam/io
import gleam/string
import gleamy/bench
import input
import simplifile

pub fn main() {
  let args = argv.load().arguments
  case args {
    ["-f", file] | ["--file", file] -> {
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
    ["-s", string] | ["--string", string] ->
      io.println(algorithm.md5(bit_array.from_string(string)))

    ["-b", ..] | ["--benchmark", ..] -> {
      bench.run(
        [
          bench.Input("Empty", ""),
          bench.Input("Small", "test"),
          bench.Input("Block-64B", string.repeat("a", 64)),
          bench.Input("Medium-1KB", string.repeat("a", 1024)),
          bench.Input("Large-1MB", string.repeat("a", 1024 * 1024)),
        ],
        [
          bench.SetupFunction("md5", fn(item) {
            fn(items) { algorithm.md5(bit_array.from_string(items)) }
          }),
        ],
        [bench.Duration(1000), bench.Warmup(100)],
      )
      |> bench.table([bench.IPS, bench.Min, bench.Max, bench.Mean, bench.P(99)])
      |> io.println()
      io.print("Measured in ms")
    }

    ["-i", ..] | ["--interactive"] -> {
      start_interactive_mode()
    }

    _ ->
      io.println(
        "Usages:\n  gleam run -- -f/--file <file>\n  gleam run -- -s/--string <string> \n  gleam run -- -b/--benchmark\n   gleam run -- -i/--interactive",
      )
  }
}

fn start_interactive_mode() {
  io.println("Write quit or exit to exit interactive mode")
  let assert Ok(line) = input.input(prompt: "Enter your string: ")

  case line {
    "quit" | "exit" -> {
      io.println("Goodbye!")
    }

    s -> {
      io.println(algorithm.md5(bit_array.from_string(s)))
      start_interactive_mode()
    }
  }
}
