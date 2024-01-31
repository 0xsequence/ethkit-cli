package internal

import (
	"math/big"
)

func HexToBigInt(str *string) *big.Int {
	if str == nil {
		return nil
	}
	num := new(big.Int)
	num.SetString((*str)[2:], 16)

	return num
}
