package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"regexp"
	"strings"

	"github.com/makevoid/vanieth/lib"
	"github.com/ogier/pflag"
)

// main, executes addrGen ad-infinitum, until the required matches are found
func main() {
	var (
		runSeconds    int64
		maxProcesses  int
		foundCount    int
		findCount     int
		quietMode     bool
		privateKey    string
		sourceAddress string
	)

	results := make(chan *lib.Match, 10000)
	matcher := &lib.Matcher{
		Results: results,
	}

	flag := pflag.NewFlagSet("vanieth", pflag.ExitOnError)

	flag.Usage = func() {
		println("Usage:")
		println("  vanieth [-acilqs] [-n num] [-d dist] (-p key | search)")
		println()
		flag.PrintDefaults()
		println()
		lib.PrintUsageExamples()
	}

	flag.BoolVarP(&matcher.FindInMain, "address", "a", false, "Search for results in the main address (can specify with -c to search both at once)")
	flag.BoolVarP(&matcher.FindInContract, "contract", "c", false, "Search through first \"distance\" number of contract addresses (or 10 if unspecified)")
	flag.BoolVarP(&matcher.ShowContractAddresses, "list", "l", false, "List all contract addresses within given \"distance\" number along with output")
	flag.BoolVarP(&matcher.DoNotChecksum, "no-sum", "s", false, "Don't convert the address to a checksum address")
	flag.BoolVarP(&matcher.IgnoreCase, "ignore-case", "i", false, "Search in case-insensitive fashion")
	flag.BoolVarP(&quietMode, "quiet", "q", false, "Don't print out speed progress updates, just the found addresses (forced if not TTY)")
	flag.IntVarP(&matcher.ContractDepth, "distance", "d", 0, "Specify `depth` of contract addresses to search (only if -c or -l specified)")
	flag.IntVarP(&findCount, "count", "n", 0, "Keep searching until this many `results` have been found")
	flag.IntVarP(&maxProcesses, "max-procs", "", 0, "Set number of simultaneous processes (default = numCPUs)")
	flag.Int64VarP(&runSeconds, "timed", "t", 0, "Allow to run for given number of `seconds`")
	flag.StringVarP(&privateKey, "key", "", "", "Specify a single private `key` to display")
	flag.StringVarP(&sourceAddress, "scan", "", "", "Scan a specified source address (only useful for searching contract addresses)")
	flag.Parse(os.Args[1:])

	if !matcher.FindInMain && !matcher.FindInContract {
		matcher.FindInMain = true
	}

	if matcher.FindInContract && matcher.ContractDepth == 0 {
		matcher.ContractDepth = 10
	}

	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		quietMode = true
	}

	if maxProcesses == 0 {
		maxProcesses = runtime.NumCPU()
	}

	args := flag.Args()
	if len(args) < 1 && privateKey == "" && sourceAddress == "" {
		println("Cannot search, no search string provided")
		println()
		flag.Usage()
		os.Exit(1)
	} else if len(args) > 0 {
		matchArg := args[0]

		// Strip off the 0x
		matchArg = strings.TrimPrefix(matchArg, "0x")

		if matcher.IgnoreCase {
			matchArg = strings.ToLower(matchArg)
		}

		if isPlain, _ := regexp.MatchString(`^[0-9a-fA-F]$`, matchArg); isPlain {
			// This is a plain prefix matchArg, we don't need a regular expression
			matcher.Prefix = "0x" + matchArg
		} else {
			matcher.Regex = regexp.MustCompile("^0x" + matchArg)
		}
	}

	if privateKey != "" {
		account, err := lib.PrivateKeyAccount(privateKey)
		if err != nil {
			println("Error creating account from private key", err)
			os.Exit(1)
		}
		match := matcher.Match(account)
		j, _ := json.Marshal(match)
		println(string(j))
		return
	}

	if sourceAddress != "" {
		account, err := lib.AddressAccount(sourceAddress)
		if err != nil {
			println("Error creating account from address", err)
			os.Exit(1)
		}
		matcher.FindInContract = true // Only makes sense for searching contract addresses
		if matcher.ContractDepth == 0 {
			matcher.ContractDepth = 10
		}
		match := matcher.Match(account)
		j, _ := json.Marshal(match)
		println(string(j))
		return
	}

	var ctx = context.Background()
	var cancel context.CancelFunc
	if runSeconds > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(runSeconds)*time.Second)
	} else {
		ctx, cancel = context.WithCancel(ctx)
		if findCount == 0 {
			// We default to count = 1 if no run time specified
			findCount = 1
		}
	}

	go func() {
		tock := time.NewTicker(time.Second)
		if quietMode {
			tock.Stop()
		}

		var rated bool
		defer func() {
			if rated {
				// Don't leave the new-line hanging
				println()
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return

			case <-tock.C:
				fmt.Printf("\rRate: %s/sec   \b\b", lib.FormatRate(lib.SearchRate()))
				rated = true

			case f := <-results:
				foundCount++
				j, _ := json.Marshal(f)
				if quietMode {
					println(string(j))
				} else {
					fmt.Printf("\r%s\n", string(j))
					rated = false
				}

				if findCount > 0 && foundCount >= findCount {
					// We have met our requirements, time to leave
					cancel()
					return
				}
			}
		}
	}()

	// Create a semaphore that will count up to maxProcesses
	semaphore := make(chan bool, maxProcesses)

	// While running
	for {
		select {
		case <-ctx.Done():
			// Not running anymore
			return

		default:
			semaphore <- true // Keep going until semaphore fills up
			if ctx.Err() == nil {
				// We're not shutting down, so run another
				matcher.Run(ctx, semaphore)
			}
		}
	}
}
