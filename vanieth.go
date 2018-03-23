package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
	"golang.org/x/crypto/ssh/terminal"
	"./lib"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ogier/pflag"
)

var mu sync.Mutex

func checksumAddr(addr string) string {
	a := []byte(addr)
	hash := hex.EncodeToString(crypto.Keccak256(a))
	for i := 0; i < len(a); i++ {
		c := a[i]
		if c >= 'a' {
			if hash[i] >= '8' {
				a[i] = c - 0x20
			}
		}
	}
	return string(a)
}

func contractGen(addr []byte, pos int) string {
	b, _ := rlp.EncodeToBytes([]interface{}{addr, uint(pos)})
	e := crypto.Keccak256(b)
	return hex.EncodeToString(e[12:])
}

// generates a public key,  address
func addrGen(toMatch *regexp.Regexp) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	addrStr := hex.EncodeToString(addr[:])
	addrMatch(addrStr, toMatch, key)
}

// tries to match the address with the string provided by the user, exits if successful
func addrMatch(addrStr string, toMatch *regexp.Regexp, key *ecdsa.PrivateKey) {
	var found bool

	f := map[string]interface{}{}

	if contractAddress || allAddresses {
		var contracts []string
		addr, _ := hex.DecodeString(addrStr)
		for i := 0; i < contractDistance; i++ {
			searchAddr := contractGen(addr, i)
			sumAddr := searchAddr
			if !noChecksum {
				sumAddr = checksumAddr(searchAddr)
				if !ignoreCase {
					searchAddr = sumAddr
				}
			}
			searchCount++
			if allAddresses {
				contracts = append(contracts, "0x"+sumAddr)
			}

			if contractAddress && !found && toMatch != nil && toMatch.MatchString(searchAddr) {
				found = true
			}
			if found && !allAddresses {
				f[fmt.Sprintf("contract-%d", i+1)] = "0x" + sumAddr
				break
			}
		}
		if allAddresses {
			f["contracts"] = contracts
		}
	}

	if mainAddress {
		searchAddr := addrStr
		if !noChecksum {
			addrStr = checksumAddr(addrStr)
			if !ignoreCase {
				searchAddr = addrStr
			}
		}
		if !found && toMatch != nil && toMatch.MatchString(searchAddr) {
			found = true
		}
		searchCount++
	}

	if found || toMatch == nil {
		mu.Lock()
		defer mu.Unlock()

		f["address"] = "0x" + addrStr
		f["public"] = hex.EncodeToString(crypto.FromECDSAPub(&key.PublicKey))
		f["private"] = hex.EncodeToString(crypto.FromECDSA(key))

		foundChan <- f
		return
	}
}

var (
	foundChan = make(chan map[string]interface{}, 10000)

	searchCount                  int64
	foundCount, totalCount       int
	contractDistance             int
	mainAddress, contractAddress bool
	noChecksum, allAddresses     bool
	ignoreCase, quietMode        bool
	privateKey                   string
)

// main, executes addrGen ad-infinitum, until the required matches are found
func main() {
	runtime.GOMAXPROCS(8)

	flag := pflag.NewFlagSet("vanieth", pflag.ExitOnError)

	flag.Usage = func () {
		println("Usage:")
		println("  vanieth [-acilqs] [-n num] [-d dist] (-p key | search)")
		println()
		flag.PrintDefaults()
		println()
		lib.PrintUsageExamples()
	}

	flag.BoolVarP(&mainAddress, "address", "a", false, "Search for results in the main address (can specify with -c to search both at once)")
	flag.BoolVarP(&contractAddress, "contract", "c", false, "Search through first \"distance\" number of contract addresses (or 10 if unspecified)")
	flag.BoolVarP(&allAddresses, "list", "l", false, "List all contract addresses within given \"distance\" number along with output")
	flag.BoolVarP(&noChecksum, "no-sum", "s", false, "Don't convert the address to a checksum address")
	flag.BoolVarP(&ignoreCase, "ignore-case", "i", false, "Search in case-insensitive fashion")
	flag.BoolVarP(&quietMode, "quiet", "q", false, "Don't print out speed progress updates, just the found addresses (forced if not TTY)")
	flag.IntVarP(&contractDistance, "distance", "d", 0, "Specify `depth` of contract addresses to search (only if -c or -l specified)")
	flag.IntVarP(&totalCount, "count", "n", 1, "Keep searching until this many `results` have been found")
	flag.StringVarP(&privateKey, "private", "p", "", "Specify a single private `key` to display")
	flag.Parse(os.Args[1:])

	if !mainAddress && !contractAddress {
		mainAddress = true
	}

	if contractAddress && contractDistance == 0 {
		contractDistance = 10
	}

	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		quietMode = true
	}

	if privateKey != "" {
		keyBytes, err := hex.DecodeString(privateKey)
		if err != nil {
			println("Cannot convert private key from hex")
		}
		key, err := crypto.ToECDSA(keyBytes)
		if err != nil {
			println("Cannot parse private key", err)
		}

		addr := crypto.PubkeyToAddress(key.PublicKey)
		addrStr := hex.EncodeToString(addr[:])

		addrMatch(addrStr, nil, key)

		f := <-foundChan
		j, _ := json.Marshal(f)
		println(string(j))
		return
	}

	var match string

	args := flag.Args()
	if len(args) != 1 {
		println("Cannot search, no search string provided")
		println()
		flag.Usage()
		os.Exit(1)
	} else {
		match = args[0]
	}

	if ignoreCase {
		match = strings.ToLower(match)
	}

	toMatch := regexp.MustCompile("^" + match)

	go func() {
		tock := time.NewTicker(time.Second)
		if quietMode {
			tock.Stop()
		}
		for {
			select {
			case <-tock.C:
				var n int64
				n, searchCount = searchCount, 0
				fmt.Printf("\rRate: %s/sec   \b\b", lib.FormatRate(n))
			case f := <-foundChan:
				foundCount++
				j, _ := json.Marshal(f)
				if quietMode {
					println(string(j))
				} else {
					fmt.Printf("\r%s\n", string(j))
				}
				if foundCount >= totalCount {
					os.Exit(0)
				}
			}
		}
	}()

	for {
		go addrGen(toMatch)
		time.Sleep(1 * time.Microsecond)
	}
}
