package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const base62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var max = big.NewInt(int64(len(base62)))

func GenerateShortID() string {
	var b strings.Builder
	b.Grow(6)

	for i := 0; i < 6; i++ {
		n, _ := rand.Int(rand.Reader, max)
		b.WriteByte(base62[n.Int64()])
	}

	return b.String()
}
