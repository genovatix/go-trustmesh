package main

import (
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"trustmesh/common" // Adjust the import path to where your common package is located
)

func main() {
	// Step 1: Define Latitude and Longitude
	// Example coordinates for the Golden Gate Bridge
	lat := 37.8199
	lon := -122.4783

	// Step 2: Initialize NetworkAddress
	address, _ := common.GenerateAddress(lat, lon, 256)

	fmt.Println("Successfully generated a valid NetworkAddress with ZKP.")

	spew.Dump(address)
	log.Println(hex.EncodeToString([]byte(address.PublicKey)))

}
