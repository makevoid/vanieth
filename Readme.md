# Vanieth

<img width="653" height="653" alt="vanieth gopher" src="https://github.com/user-attachments/assets/2501fd0b-4196-4f74-86fe-6f908498a5fe" />

> ⚡ An Ethereum vanity address generator written in Go many years ago by [@makevoid](https://twitter.com/makevoid) with @norganna contributing to it.
> The project was originally aimed to be didactical and easy to read and it's only `~400` lines of Go.


[![Go Version](https://img.shields.io/badge/go-%3E%3D1.17-blue.svg)](https://golang.org/)
[![Docker](https://img.shields.io/badge/docker-ready-brightgreen.svg)](https://hub.docker.com/r/makevoid/vanieth)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**Note:** This is a stable release that works well but is not currently under active development.

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
docker run makevoid/vanieth ./vanieth abc
```

### Using Go

```bash
# Install
go get github.com/makevoid/vanieth

# Run
$GOPATH/bin/vanieth 42

# Or add to PATH for convenience
export PATH=$PATH:$GOPATH/bin
vanieth 42
```

## 📖 Overview

Vanieth generates Ethereum vanity addresses with custom patterns. It leverages parallel processing to efficiently search for addresses matching your criteria, whether you're looking for simple prefixes, complex regex patterns, or specific contract addresses.

### Example Output

```json
{
  "address": "0x42f32B004Da1093d51AE40a58F38E33BA4f46397",
  "private": "4774628228852ee570d188f92cd10df3282bb5d895fc701733f43fca6bfb9852",
  "public": "04d811caac49ba458fda498e5bc385bc9cc6e67aa6b19ba754c6cd75953ef06310e8607798ce5810a0b32fbd41fe8915de52fd511e7660038ff7067a0e94fc9481"
}
```

> ⚠️ **Security Note**: The generation time increases exponentially with pattern length. A 4-character pattern takes significantly longer than a 2-character pattern.

## 🛠️ Installation

### Prerequisites

- **Go**: Version 1.17 or higher (uses go-ethereum crypto libraries)
- **Docker**: (Optional) For containerized execution

### Build from Source

```bash
# Clone and enter the repository
git clone https://github.com/makevoid/vanieth.git
cd vanieth

# Build the binary
./build.sh

# Run
./vanieth abc
```

### Docker Build

```bash
# Build the container
docker-compose build

# Test run
docker-compose run vanieth ./vanieth abc
```

## 💡 Usage

### Basic Syntax

```
vanieth [-acilqs] [-n num] [-d dist] (-key=key | -scan=address | search)
```

### Options

| Flag | Long Form | Description |
|------|-----------|-------------|
| `-a` | `--address` | Search in the main address (combine with `-c` to search both) |
| `-c` | `--contract` | Search through contract addresses |
| `-n` | `--count` | Number of results to find before stopping |
| `-d` | `--distance` | Depth of contract addresses to search |
| `-i` | `--ignore-case` | Case-insensitive search |
| `-l` | `--list` | List all contract addresses within distance |
| `-s` | `--no-sum` | Skip checksum address conversion |
| `-q` | `--quiet` | Suppress progress updates |
| `-t` | `--timed` | Run for specified number of seconds |
| | `--key` | Display details for a specific private key |
| | `--scan` | Scan a specified source address |
| | `--max-procs` | Set number of parallel processes (default: CPU count) |

## 📚 Examples

### Find Simple Patterns

```bash
# Find address starting with "ABC"
vanieth 'ABC'

# Find 3 addresses with "ABC" prefix
vanieth -n 3 'ABC'

# Search for 5 seconds
vanieth -t 5 'ABC'
```

### Regular Expression Patterns

```bash
# Address containing "ABC" anywhere
vanieth '.*ABC'

# Address ending with "DEF"
vanieth '.*DEF$'

# Case-insensitive: starts and ends with 'A'
vanieth -i 'A.*A$'

# Address with "AB" after 2+ zeros
vanieth '00+AB'
```

### Contract Address Search

```bash
# Search in first 10 contract addresses
vanieth -c 'ABC'

# Search in first contract address only
vanieth -cd1 '00+AB'

# List first 5 contract addresses
vanieth -ld5 --key=0x349fbc254ff918305ae51967acc1e17cfbd1b7c7e84ef8fa670b26f3be6146ba
```

### Address Analysis

```bash
# Show contract addresses for existing address
vanieth -l --scan=0x950024ae4d9934c65c9fd04249e0f383910d27f2
```

## ⚡ Performance Tips

1. **Pattern Length**: Each additional character exponentially increases search time
2. **Parallel Processing**: Use `--max-procs` to optimize for your CPU
3. **Regex Complexity**: Simple prefixes are faster than complex regex patterns
4. **Contract Search**: Limiting distance (`-d`) improves performance

## 🔒 Security Considerations

- **Private Keys**: Never share or expose generated private keys
- **Randomness**: Uses cryptographically secure random generation
- **Verification**: Always verify addresses on a testnet first
- **Storage**: Store private keys securely and encrypted

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- [@makevoid](https://twitter.com/makevoid)
- @norganna

## 🙏 Acknowledgments

- Built with Go using [go-ethereum](https://github.com/ethereum/go-ethereum) crypto libraries
- Inspired by the Ethereum community's need for memorable addresses
- Thanks to all contributors and users

---

**NOTES:** Original Logo of the project: https://github.com/makevoid/vanieth/blob/master/screenshots/readme_banner.png

---

**⭐ Star this repository if you find it useful!**

**🐛 Found a bug?** [Open an issue](https://github.com/makevoid/vanieth/issues)



**💬 Questions?** [Start a discussion](https://github.com/makevoid/vanieth/discussions)
