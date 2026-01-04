package main

import (
	"flag"
	"go_md5/bitutil"
	"os"
)

func main() {

	name := flag.String("s", "", "Convert string")
	filePath := flag.String("f", "", "Convert file path")

	flag.Parse()

	if (*name != "" && *filePath != "") || *name == *filePath {
		println("Usage: -s=\"string\" or -f=\"file path\"")
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
