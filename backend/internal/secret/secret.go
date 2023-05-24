package secret

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// New generates and returns a new secret.
func New() string {
	secretWords := make([]string, 4)
	for i := range secretWords {
		bigIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
		secretWords[i] = words[bigIndex.Int64()]
	}
	return strings.Join(secretWords, " ")
}
