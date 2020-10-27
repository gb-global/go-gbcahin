package main

import (
	"fmt"
	"io/ioutil"

	"gbchain-org/go-gbchain/common/hexutil"
)

func main() {
	data, err := ioutil.ReadFile("crossDemo.abi")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hexutil.Encode(data))
}
