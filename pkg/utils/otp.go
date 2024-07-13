package utils

import (
	"crypto/rand"
	"math"
	"math/big"
)

func GenerateOTP(maxDigits uint32) int {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, float64(maxDigits)))),
	)
	if err != nil {
		panic(err)
	}
	return int(bi.Int64())
}