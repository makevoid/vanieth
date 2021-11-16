# Vanieth

![](https://github.com/makevoid/vanieth/blob/master/screenshots/readme_banner.png?raw=true)

A comprehensive and fast Ethereum vanity address "generator" written in golang.

### Docker Run

If you are just interested in running the program I recently added a way to run the CLI app via Docker:

    docker run makevoid/vanieth ./vanieth abc


### Golang Run - Prerequisites:

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
  vanieth [-acilqs] [-n num] [-d dist] (-key=key | -scan=address | search)

  -a, --address
    	Search for results in the main address (can specify with -c to search both at once)
  -c, --contract
    	Search through first "distance" number of contract addresses (or 10 if unspecified)
  -n, --count results
    	Keep searching until this many results have been found
  -d, --distance depth
    	Specify depth of contract addresses to search (only if -c or -l specified)
  -i, --ignore-case
    	Search in case-insensitive fashion
  --key key
    	Specify a single private key to display
  -l, --list
    	List all contract addresses within given "distance" number along with output
  --max-procs int
    	Set number of simultaneous processes (default = numCPUs)
  -s, --no-sum
    	Don't convert the address to a checksum address
  -q, --quiet
    	Don't print out speed progress updates, just the found addresses (forced if not TTY)
  --scan string
    	Scan a specified source address (only useful for searching contract addresses)
  -t, --timed seconds
    	Allow to run for given number of seconds
```

#### Examples:

```
vanieth -n 3 'ABC'
```

Find 3 addresses that have `ABC` at the beginning.

```
vanieth -t 5 'ABC'
```

Find as many address that have `ABC` at the beginning as possible within 5 seconds.


```
vanieth -c 'ABC'
```

Find any address that has `ABC` at the beginning of any of the first 10 contract addresses.

```
vanieth -cd1 '00+AB'
```

Find any address that has `AB` after 2 or more `0` chars in the first contract address.

```
vanieth '.*ABC'
```

Find a single address that contains `ABC` anywhere.

```
vanieth '.*DEF$'
```

Find a single address that contains `DEF` at the end.

```
vanieth -i 'A.*A$'
```

Find a single address that contains either `A` or `a` at both the start and end.

```
vanieth -ld1 '.*ABC'
```

Find a single address that contains `ABC` anywhere, and also list the first contract address.

```
vanieth -ld5 --key=0x349fbc254ff918305ae51967acc1e17cfbd1b7c7e84ef8fa670b26f3be6146ba
```

List the details and first five contract address for the supplied private key.

```
vanieth -l --scan=0x950024ae4d9934c65c9fd04249e0f383910d27f2
```

Show the first 10 contract addresses of the supplied address.

### Go Build

`go get` the project following the instructions at the top

cd into your $GOROOT where this project is located

Build the executable with

    ./build.sh

### Docker Build

Now that you have built the executable you can package it up as a docker container via docker-compose

Just run:

    docker-compose build

this will build the docker container with your changes.

Test run via:

    docker-compose run vanieth ./vanieth abc

---


Enjoy,

[@makevoid](https://twitter.com/makevoid) & @norganna
