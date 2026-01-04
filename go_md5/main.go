package main

import "fmt"

func main() {
	table := generateKtable()

	for i, v := range table {
		fmt.Printf("K[%2d] = 0x%08x\n", i, v)
	}
}
