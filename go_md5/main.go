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
	interactiveMode := flag.Bool("i", false, "Start interactive mode")

	flag.Parse()

	if *interactiveMode {
		startInteractiveMode()
		return
	}

	if (*name != "" && *filePath != "") || *name == *filePath {
		println("Usage: -s=\"string\" , -f=\"file path\", -i ")
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
