package tools

import (
	"math/big"
)

// input two string and get the sum of them, big output
func Sum(a, b string) string {
	if a == "" {
		a = "0"
	}
	if b == "" {
		b = "0"
	}
	//convert
	var bigA, bigB, sum *big.Int
	var ok bool
	if bigA, ok = big.NewInt(0).SetString(a, 10); !ok {
		return ""
	}
	if bigB, ok = big.NewInt(0).SetString(b, 10); !ok {
		return ""
	}
	return sum.Add(bigA, bigB).String()
}
