# Vanieth

A comprehensive and fast Ethereum vanity address "generator" written in golang.

### Prerequisites:

You have to have go (golang) installed

Go get this repo:

    go get github.com/makevoid/vanieth

Run it:

    $GOPATH/bin/vanieth

Copy it to your path, or add $GOPATH/bin to your path

### Example run:

```
$ vanieth 42
{"address":"0x42f32B004Da1093d51AE40a58F38E33BA4f46397","private":"4774628228852ee570d188f92cd10df3282bb5d895fc701733f43fca6bfb9852","public":"04d811caac49ba458fda498e5bc385bc9cc6e67aa6b19ba754c6cd75953ef06310e8607798ce5810a0b32fbd41fe8915de52fd511e7660038ff7067a0e94fc9481"}
```

The returned address and private key are in hex format. As you can see the ethereum address starts with the mythical `42`.

Here's a more complex vanity address, this will take a significally longer time to do a complete run.

```
$ vanieth 1234
{"address":"0x12341b4c716B8FCFA8E13A83CA3dFd2c6051E60D","private":"ee50661eb0080cd36ce380f3dad5511c91f97ccee67bd14dc7a91335a34720d1","public":"04e0526fbc5552e4ff117a5c065ad3ce6f8211e160e12bdd3dded3dab2bfc268916489ed2c8d4af6c624406085c5e9a6946bdfbe0d74de26384a7c9baaf6f2de64"}
```

The more chars you add, the longer the time will be, exponentially!

### Advanced features

Apart from searching for just prefixes you can also search for contract addresses, regular expressions and dump details for an existing or previously found private key:

### Usage

```
Usage:
  vanieth [-acilqs] [-n num] [-d dist] (-p key | search)

  -a, --address
    	Search for results in the main address (can specify with -c to search both at once)
  -c, --contract
    	Search through first "distance" number of contract addresses (or 10 if unspecified)
  -n, --count results
    	Keep searching until this many results have been found (default 1)
  -d, --distance depth
    	Specify depth of contract addresses to search (only if -c or -l specified)
  -i, --ignore-case
    	Search in case-insensitive fashion
  -l, --list
    	List all contract addresses within given "distance" number along with output
  -s, --no-sum
    	Don't convert the address to a checksum address
  -p, --private key
    	Specify a single private key to display
  -q, --quiet
    	Don't print out speed progress updates, just the found addresses (forced if not TTY)
```

#### Examples:

```vanieth -n 3 'ABC'```

Will find 3 addresses that have `ABC` at the beginning.

```vanieth -c 'ABC'```

Will find any address that has `ABC` at the beginning of any of the first 10 contract addresses.

```vanieth -cd1 '00+AB'```

Will find any address that has `AB` after 2 or more `0` chars in the first contract address.

```vanieth '.*ABC'```

Will find a single address that contains `ABC` anywhere.

```vanieth '.*DEF$'```

Will find a single address that contains `DEF` at the end.

```vanieth -i 'A.*A$'```

Will find a single address that contains either `A` or `a` at both the start and end.

```vanieth -ld1 '.*ABC'```

Will find a single address that contains `ABC` anywhere, and also list the first contract address.

```vanieth -ld5 -p '349fbc254ff918305ae51967acc1e17cfbd1b7c7e84ef8fa670b26f3be6146ba'```

Will list the details and first five contract address for the supplied private key.


Enjoy,

[@makevoid](https://twitter.com/makevoid) & @norganna
