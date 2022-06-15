package chains

import (
	"fmt"
	"math/big"
	"testing"
)

func TestTokenAmount_Float64(t *testing.T) {
	b := big.NewInt(12345600000000005)
	token := TokenAmount{
		Decimals: 18,
		Amount:   b,
	}
	tFloat, err := token.Float64()
	if err != nil {
		panic(err)
	}
	fmt.Println(tFloat)
}
