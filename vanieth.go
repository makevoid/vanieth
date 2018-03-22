package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/ethereum/go-ethereum/crypto" // you need to `go get` this, as is not in the stdlib
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

	flag.Usage = func() {
		println("Usage:")
		println("  vanieth [-acilqs] [-n num] [-d dist] (-p key | search)")
		println()
		flag.PrintDefaults()
		println()
		println("Examples:")
		println()
		println("  vanieth -n 3 'ABC'")
		println("     Will find 3 addresses that have `ABC` at the beginning.")
		println()
		println("  vanieth -c 'ABC'")
		println("     Will find any address that has `ABC` at the beginning of any of the first 10 contract addresses.")
		println()
		println("  vanieth -cd1 '00+AB'")
		println("     Will find any address that has `AB` after 2 or more `0` chars in the first contract address.")
		println()
		println("  vanieth '.*ABC'")
		println("     Will find a single address that contains `ABC` anywhere.")
		println()
		println("  vanieth '.*DEF$'")
		println("     Will find a single address that contains `DEF` at the end.")
		println()
		println("  vanieth -i 'A.*A$'")
		println("     Will find a single address that contains either `A` or `a` at both the start and end.")
		println()
		println("  vanieth -ld1 '.*ABC'")
		println("     Will find a single address that contains `ABC` anywhere, and also list the first contract address.")
		println()
		println("  vanieth -ld5 -p '349fbc254ff918305ae51967acc1e17cfbd1b7c7e84ef8fa670b26f3be6146ba'")
		println("     Will list the details and first five contract address for the supplied private key.")
		println()
	}
	flag.BoolVarP(&mainAddress, "address", "a", false, "Search in the main address")
	flag.BoolVarP(&contractAddress, "contract", "c", false, "Search through first `distance` contract addresses (or 10 if unspecified)")
	flag.IntVarP(&contractDistance, "distance", "d", 0, "Specify distance into contract to search")
	flag.BoolVarP(&allAddresses, "list", "l", false, "List all `distance` contract addresses with result")
	flag.BoolVarP(&noChecksum, "no-sum", "s", false, "Don't convert to checksum address")
	flag.BoolVarP(&ignoreCase, "ignore-case", "i", false, "Search in case-insensitive fashion")
	flag.BoolVarP(&quietMode, "quiet", "q", false, "Don't print out speed progress updates, just the found addresses (forced if not TTY)")
	flag.IntVarP(&totalCount, "count", "n", 1, "How many results to find")
	flag.StringVarP(&privateKey, "private", "p", "", "Specify a single private key to display")
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
		fmt.Printf("%s\n", string(j))
		return
	}

	var match string

	args := flag.Args()
	if len(args) != 1 {
		println("Cannot search, no search string provided", flag.Args())
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
				fmt.Printf("\rRate: %s/sec   \b\b", Format(n))
			case f := <-foundChan:
				foundCount++
				j, _ := json.Marshal(f)
				if quietMode {
					fmt.Printf("%s\n", string(j))
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

func Format(n int64) string {
	in := strconv.FormatInt(n, 10)
	out := make([]byte, len(in)+(len(in)-2+int(in[0]/'0'))/3)
	if in[0] == '-' {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}
