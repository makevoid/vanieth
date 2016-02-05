# Vanieth

A simple Ethereum vanity address "generator" written in golang.

### Prerequisites:

You have to have go (golang) installed

Then download this repo into your GOPATH, cd into it and run:

    go get

to get the repo only dependency (`github.com/ethereum/go-ethereum/crypto`)

### Example run:

    go run vanieth.go 42

    Address found:
    addr: 0x429e6a85ed72fddf6c5679da1ac033ab65ad68a7
    pvt: 0x0ce2f8e425121d5b7f078b6bce4c9bf23937ee4fd9b62ff2e81d84b724eb5e1b

The returned address and private key are in hex format. As you can see the ethereum address starts with the mythical `42`.

Here's a more complex vanity address, this will take a significally longer time to do a complete run.

    go run vanieth.go 1234

    Address found:
    addr: 0x123411cc4a2e2e3238ee8e22d0d7b3cf2c8add9c
    pvt: 0x208439bf49edbc236bcffaa831e32006b91e6251150992fe5e704a3c3870415d

The more chars you add, the longer the time will be, exponentially!

### More efficient run:

Compile and run:

go build vanieth.go; ./vanieth 1234

This will be a bit more efficient than the examples shown above.


Enjoy,

[@makevoid](https://twitter.com/makevoid)
