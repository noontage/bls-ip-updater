package main

import (
	"bls-ip-updater/src/bls"
	"bls-ip-updater/src/conoha"
	"fmt"
)

func main() {
	bc := bls.NewClient()
	addr, err := bc.GetGlobalIP()
	if err != nil {
		panic(err)
	}
	fmt.Println(addr.String())

	cc := conoha.NewClient()
	if err := cc.Test(); err != nil {
		panic(err)
	}
}
