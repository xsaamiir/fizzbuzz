package fizzbuzz

import (
	"strconv"
)

// FizzBuzz returns a list of strings with numbers from 1 to limit,
// where:
// 	- all multiples of int1 are replaced by str1
// 	- all multiples of int2 are replaced by str2
// 	- all multiples of int1 and int2 are replaced by str1str2.
func FizzBuzz(int1, int2, limit int, str1, str2 string) []string {
	res := make([]string, limit)

	for i := 1; i <= limit; i++ {
		var s string

		switch {
		case i%int1 == 0 && i%int2 == 0:
			s = str1 + str2
		case i%int1 == 0:
			s = str1
		case i%int2 == 0:
			s = str2
		default:
			s = strconv.Itoa(i)
		}

		res[i-1] = s
	}

	return res
}
