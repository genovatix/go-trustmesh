package main

import (
	"log"
	"trustmesh/common"
)

func main() {
	netAddr := common.NewNetworkAddress(-45, 75)
	addr, err := netAddr.EncodeToString([]byte("12345678910111213141516111111111"))
	if err != nil {
		panic(err)
	}
	pub := netAddr.PublicKeyBase64()
	log.Println(addr)
	log.Println(pub)
}
