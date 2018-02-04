package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto" // you need to `go get` this, as is not in the stdlib
)

// Example run:
//
// $ go run vanieth.go 1234
//
// Address found:
// addr: 123411cc4a2e2e3238ee8e22d0d7b3cf2c8add9c
// pvt: 208439bf49edbc236bcffaa831e32006b91e6251150992fe5e704a3c3870415d
//
// https://github.com/ethereum/go-ethereum
//

// "main" method, generates a public key,  address
//
func addrGen(toMatch string) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	addrStr := hex.EncodeToString(addr[:])
	addrMatch(addrStr, toMatch, key)
}

// tries to match the address with the string provided by the user, exits if successful
//
func addrMatch(addrStr string, toMatch string, key *ecdsa.PrivateKey) {
	toMatch = strings.ToLower(toMatch)
	addrStrMatch := strings.TrimPrefix(addrStr, toMatch)
	found := addrStrMatch != addrStr
	if found {
		// fmt.Println("pub:", hex.EncodeToString(crypto.FromECDSAPub(&key.PublicKey))) // uncomment if you want the public key
		keyStr := hex.EncodeToString(crypto.FromECDSA(key))
		addrFound(addrStr, keyStr)
		os.Exit(0) // here the program exits when it found a match
	}
}

// main, executes addrGen ad-infinitum, until a match is found
//
func main() {
	runtime.GOMAXPROCS(8)

	var toMatch string
	if len(os.Args) == 1 {
		errNoArg()
		os.Exit(1)
	} else {
		toMatch = os.Args[1]
	}

	for true {
		go addrGen(toMatch)
		time.Sleep(1 * time.Microsecond)
	}
}

// non-interesting functions follow...

func addrFound(addrStr string, keyStr string) {
	println("Address found:")
	fmt.Printf("addr: 0x%s\n", addrStr)
	fmt.Printf("pvt: 0x%s\n", keyStr)
	println("\nexiting...")
}

func errNoArg() {
	println("You need to pass a vanity match, retry with an extra agrument like: 42")
	println("\nexample: go run vanieth.go 42")
	println("\nexiting...")
}
