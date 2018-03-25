package lib

import (
	"crypto/ecdsa"
	"encoding/hex"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Account represents an account's details, such as key and address.
type Account struct {
	Key  *ecdsa.PrivateKey
	Addr common.Address
}

// CreateAccount will return a new randomly generate account, including key.
func CreateAccount() *Account {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	return &Account{
		Key:  key,
		Addr: addr,
	}
}

// AddressAccount will return a new account without any keys, based around a public address.
func AddressAccount(addr string) (*Account, error) {
	a, err := NewAddress(addr)
	if err != nil {
		return nil, err
	}

	return &Account{
		Addr: a,
	}, nil
}

// PrivateKeyAccount will return a new account based on an input private key.
func PrivateKeyAccount(private string) (*Account, error) {
	k, err := DecodeHex(private)
	if err != nil {
		return nil, err
	}

	key, err := crypto.ToECDSA(k)
	if err != nil {
		return nil, err
	}

	addr := crypto.PubkeyToAddress(key.PublicKey)
	return &Account{
		Key:  key,
		Addr: addr,
	}, nil
}

// Contract gets the contract address at nonce `pos`.
//
// Contract addresses are deterministic based on the order of generation. For a given account address,
// the first contract created by that address will always have a certain predictable address. This
// method works out what that address will be at the given ordinal position.
func (a *Account) Contract(pos int) common.Address {
	return crypto.CreateAddress(a.Addr, uint64(pos))
}

// GetContracts returns the first `n` contracts in a map.
func (a *Account) GetContracts(n int) (contracts map[int]common.Address) {
	contracts = map[int]common.Address{}
	for i := 0; i < n; i++ {
		contracts[i+1] = crypto.CreateAddress(a.Addr, uint64(i))
	}
	return
}

// PublicKey returns the public key as a data string.
func (a *Account) PublicKey() string {
	if a.Key == nil {
		return ""
	}
	return "0x" + hex.EncodeToString(crypto.FromECDSAPub(&a.Key.PublicKey))
}

// PrivateKey returns the private key as a data string.
func (a *Account) PrivateKey() string {
	if a.Key == nil {
		return ""
	}
	return "0x" + hex.EncodeToString(crypto.FromECDSA(a.Key))
}

// NewAddress returns a parsed address from an address string.
func NewAddress(s string) (common.Address, error) {
	b, err := DecodeHex(s)
	return common.BytesToAddress(b), err
}

// DecodeHex will decode a hex string, and strip off any 0x prefix.
func DecodeHex(s string) ([]byte, error) {
	return hex.DecodeString(strings.TrimPrefix(s, "0x"))
}
