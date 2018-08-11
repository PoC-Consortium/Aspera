package crypto

import (
	"crypto/sha256"
)

func secretPhraseToPrivateKey(secretPhrase string) *[32]byte {
	bs := sha256.Sum256([]byte(secretPhrase))
	bs[31] &= 0x7F
	bs[31] |= 0x40
	bs[0] &= 0xF8
	return &bs
}
