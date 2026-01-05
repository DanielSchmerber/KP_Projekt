package main

import (
	"bufio"
	"flag"
	"go_md5/bitutil"
	"os"
	"strings"
)

func main() {

	name := flag.String("s", "", "Convert string")
	filePath := flag.String("f", "", "Convert file path")
	help := flag.Bool("h", false, "Prints help message")
	interactiveMode := flag.Bool("i", false, "Start interactive mode")

	flag.Parse()

	if *interactiveMode {
		startInteractiveMode()
		return
	}

	if *help {
		printHelpText()
		return
	}

	if (*name != "" && *filePath != "") || *name == *filePath {
		printHelpText()
		return
	}

	var bits *bitutil.BitArray

	if *name != "" {
		bits = bitutil.NewBitArrayFromBytes([]byte(*name))
	}

	if *filePath != "" {
		bytes, err := os.ReadFile(*filePath)
		if err != nil {
			println("read file error:", err.Error())
			return
		}

		bits = bitutil.NewBitArrayFromBytes(bytes)
	}

	result, err := Md5(bits)

	if err != nil {
		println(err.Error())
	}

	println(result)
}

func printHelpText() {
	println("Usages:")
	println("  -s=\"string\"       Convert string to md5")
	println("  -f=\"file path\"    Convert file path to md5")
	println("  -h                  Prints help message")
	println("  -i                  Start interactive mode")
}

func startInteractiveMode() {
	println("Write quit or exit to exit interactive mode")

	reader := bufio.NewReader(os.Stdin)

	for {
		print("Enter your string: ")

		line, err := reader.ReadString('\n')
		if err != nil {
			println("Input closed.")
			break
		}
		//remove any trailing newlines
		line = strings.TrimRight(line, "\r\n")

		if line == "quit" || line == "exit" {
			println("Goodbye!")
			break
		}

		result, err := Md5(bitutil.NewBitArrayFromBytes([]byte(line)))
		if err != nil {
			println(err.Error())
			continue
		}

		println(result)
	}
}
