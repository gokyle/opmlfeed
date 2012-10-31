package shorten

/*
 * the URL shortening code
 */

import (
	"crypto/rand"
	"math/big"
)

// length of generated short URL ID
const ShortLen = 6

type short_index uint

func Shorten() (shortid string) {
	const symbols = "abcdefghijklmnopqrstuvwxyz0123456789"
	for i := 0; i < ShortLen; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(symbols))))
		if err != nil {
			shortid = ""
			return
		} else {
			index := n.Int64()
			shortid += string(symbols[index])
		}
	}
	return
}
