package math

import (
	"math/big"
	"strconv"
)

// Various big integer limit values.
var (
	tt255     = BigPow(2, 255)
	tt256     = BigPow(2, 256)
	tt256m1   = new(big.Int).Sub(tt256, big.NewInt(1))
	tt63      = BigPow(2, 63)
	tt64      = BigPow(2, 64)
	MaxBig256 = new(big.Int).Set(tt256m1)
	MaxBig63  = new(big.Int).Sub(tt63, big.NewInt(1))
	MaxBig64  = new(big.Int).Sub(tt64, big.NewInt(1))
	BigZero   = new(big.Int)
)

// BigPow returns a ** b as a big integer.
func BigPow(a, b int64) *big.Int {
	r := big.NewInt(a)
	return r.Exp(r, big.NewInt(b), nil)
}

// S256 interprets x as a two's complement number.
// x must not exceed 256 bits (the result is undefined if it does) and is not modified.
//
//   S256(0)        = 0
//   S256(1)        = 1
//   S256(2**255)   = -2**255
//   S256(2**256-1) = -1
func S256(x *big.Int) *big.Int {
	if x.Cmp(tt255) < 0 {
		return x
	}
	return new(big.Int).Sub(x, tt256)
}

// U256 encodes as a 256 bit two's complement number. This operation is destructive.
func U256(x *big.Int) *big.Int {
	return x.And(x, tt256m1)
}

// BigFromUint64 creates a big int from a 64 bit unsigned integer
func BigFromUint64(i uint64) *big.Int {
	var b big.Int
	b.SetString(strconv.FormatUint(i, 10), 10)
	return &b
}

func StringFromBigBytes(bs []byte) string {
	var b big.Int
	b.SetBytes(bs)
	return b.String()
}
