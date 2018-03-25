package lib

import (
	"strconv"
)

// FormatRate will output the number as a fixed string with commas.
func FormatRate(n int64) string {
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

// PrintUsageExamples will print out the various usage examples.
func PrintUsageExamples() {
	println("Examples:")
	println()
	println("  vanieth -n 3 'ABC'")
	println("     Find 3 addresses that have `ABC` at the beginning.")
	println()
	println("  vanieth -t 5 'ABC'")
	println("     Find as many address that have `ABC` at the beginning as possible within 5 seconds.")
	println()
	println("  vanieth -c 'ABC'")
	println("     Find any address that has `ABC` at the beginning of any of the first 10 contract addresses.")
	println()
	println("  vanieth -cd1 '00+AB'")
	println("     Find any address that has `AB` after 2 or more `0` chars in the first contract address.")
	println()
	println("  vanieth '.*ABC'")
	println("     Find a single address that contains `ABC` anywhere.")
	println()
	println("  vanieth '.*DEF$'")
	println("     Find a single address that contains `DEF` at the end.")
	println()
	println("  vanieth -i 'A.*A$'")
	println("     Find a single address that contains either `A` or `a` at both the start and end.")
	println()
	println("  vanieth -ld1 '.*ABC'")
	println("     Find a single address that contains `ABC` anywhere, and also list the first contract address.")
	println()
	println("  vanieth -ld5 --key=0x349fbc254ff918305ae51967acc1e17cfbd1b7c7e84ef8fa670b26f3be6146ba")
	println("     List the details and first five contract address for the supplied private key.")
	println()
	println("  vanieth -l --scan=0x950024ae4d9934c65c9fd04249e0f383910d27f2")
	println("     Show the first 10 contract addresses of the supplied address.")
	println()
}
