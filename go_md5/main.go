package main

import "go_md5/bitutil"

func main() {

	bits := bitutil.NewBitArrayFromBytes([]byte("Frank jagt im komplett verwahrlosten Taxi quer durch Bayern"))

	result, err := Md5(bits)

	if err != nil {
		println(err.Error())
	}

	println(result)
}
