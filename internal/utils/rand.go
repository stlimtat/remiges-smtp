package utils

import (
	"crypto/rand"
	"math/big"
)

func RandInt(maxValue int64) (int64, error) {
	value, err := rand.Int(rand.Reader, big.NewInt(maxValue))
	if err != nil {
		return -1, err
	}
	return value.Int64(), nil
}
