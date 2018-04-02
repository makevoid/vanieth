package lib

import (
	"context"
	"regexp"

	"fmt"
	"os"

	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
)

var currentSearches, lastSearches int64

// SearchRate returns the number of searches performed since the previous call.
func SearchRate() int64 {
	lastSearches, currentSearches = currentSearches, 0
	return lastSearches
}

// Match represents a found item.
type Match struct {
	Account   *Account       `json:"-"`
	Address   string         `json:"address"`
	Public    string         `json:"public,omitempty"`
	Private   string         `json:"private,omitempty"`
	Contracts map[int]string `json:"contracts,omitempty"`
}

// Matcher contains the matching configuration.
type Matcher struct {
	FindInMain            bool
	FindInContract        bool
	IgnoreCase            bool
	DoNotChecksum         bool
	ShowContractAddresses bool
	ContractDepth         int
	Regex                 *regexp.Regexp
	Prefix                string
	Results               chan *Match
}

// Run an instance of the searcher
func (m *Matcher) Run(ctx context.Context, semaphore chan bool) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Fprintln(os.Stderr, "\nRecovered from unexpected error:", err)
			}

			// Clear the semaphore so another process can run
			<-semaphore
		}()

		for {
			select {
			case <-ctx.Done():
				// Application is ending
				return

			default:
				if match := m.Match(CreateAccount()); match != nil {
					// Found a match!
					m.Results <- match
				}
			}
		}
	}()
}

// Match calculates if `account` matches the conditions of the `Matcher`.
func (m *Matcher) Match(account *Account) *Match {
	found, contracts := m.find(account)
	if !found {
		return nil
	}

	return &Match{
		Account:   account,
		Address:   account.Addr.Hex(),
		Public:    account.PublicKey(),
		Private:   account.PrivateKey(),
		Contracts: m.stringContracts(contracts),
	}
}

func (m *Matcher) stringContracts(c map[int]common.Address) map[int]string {
	if c == nil {
		return nil
	}

	r := map[int]string{}
	for k, v := range c {
		r[k] = m.addressString(v)
	}
	return r
}

func (m *Matcher) addressString(a common.Address) string {
	if m.DoNotChecksum {
		return "0x" + hex.EncodeToString(a.Bytes())
	}
	return a.Hex()
}

// Internal function to perform the search.
func (m *Matcher) find(account *Account) (found bool, contracts map[int]common.Address) {
	if m.FindInMain {
		if m.investigate(account.Addr) {
			found = true
			if m.ShowContractAddresses {
				contracts = account.GetContracts(m.ContractDepth)
			}
		}
	}

	if m.FindInContract {
		// If we are showing all contract addresses then either way we need to get all the addresses.
		if m.ShowContractAddresses {
			contracts = account.GetContracts(m.ContractDepth)
			for i := 0; i < m.ContractDepth; i++ {
				if m.investigate(contracts[i]) {
					found = true
					return
				}
			}
			return
		}

		for i := 0; i < m.ContractDepth; i++ {
			contract := account.Contract(i)
			if m.investigate(contract) {
				found = true
				contracts = map[int]common.Address{}
				contracts[i+1] = contract
				return
			}
		}
	}

	return
}

// Internal function to perform the matching on the address.
func (m *Matcher) investigate(addr common.Address) bool {
	var search string

	// If we're ignoring case
	if m.IgnoreCase {
		// Simply get the lowercased address string
		search = "0x" + hex.EncodeToString(addr.Bytes())
	} else {
		// Get the full EIP55 checksummed address string
		search = addr.Hex()
	}

	// Increment the search count
	currentSearches++

	// Check the prefix first if specified
	if n := len(m.Prefix); n > 0 && search[:n] != m.Prefix {
		// Prefix doesn't match
		return false
	}

	// Check the regex next if present
	if m.Regex != nil && !m.Regex.MatchString(search) {
		// Regex failed
		return false
	}

	// If the prefix didn't fail, and the regex didn't fail, then is a match.
	return true
}
