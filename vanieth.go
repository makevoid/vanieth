package main

import (
	"encoding/hex"
	"fmt"
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
		fmt.Println("pvt:", keyStr)
		fmt.Println("addr:", addrStr)
		println("")
	}
}

func main() {
	runtime.GOMAXPROCS(8)

	toMatch := "123"

	count := 10000
	for ixx := 0; ixx < count; ixx++ {
		go genAddr(toMatch)
	}

	time.Sleep(20000 * time.Millisecond)
	// crypto.PubkeyToAddress(pubKey)
}

//
// pvt: 3bfd853e59b6d38cb125fb027dc0e4d9354729c6954e0d933de11f4c2e63e012
// addr: abc5e63e9d1165c592c7f47847a60d38f93bc4bf

// https://github.com/ethereum/go-ethereum
