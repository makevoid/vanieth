package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
)

func genAddr(match string) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	addrStr := hex.EncodeToString(addr[:])
	addrStrMatch := strings.TrimPrefix(addrStr, match)
	found := addrStrMatch != addrStr
	if found {
		// fmt.Println("pub:", key.PublicKey)
		keyStr := hex.EncodeToString(crypto.FromECDSA(key))
		println("Address found:")
		fmt.Println("addr:", addrStr)
		fmt.Println("pvt:", keyStr)
		println("\nexiting...")
		os.Exit(0)
	}
}

func main() {
	runtime.GOMAXPROCS(8)

	var toMatch string
	if len(os.Args) == 1 {
		println("You need to pass a vanity match, retry with an extra agrument like: 42")
		println("\nexample: go run vanieth.go 42")
		println("\nexiting...")
		os.Exit(1)
	} else {
		toMatch = os.Args[1]
	}

	for true {
		go genAddr(toMatch)
		time.Sleep(1 * time.Millisecond)
	}
}

// Example run:
//
// $ go run vanieth.go 1234
//
// Address found:
// addr: 123411cc4a2e2e3238ee8e22d0d7b3cf2c8add9c
// pvt: 208439bf49edbc236bcffaa831e32006b91e6251150992fe5e704a3c3870415d
//
// https://github.com/ethereum/go-ethereum
