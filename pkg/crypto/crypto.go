package crypto

import (
	"crypto/sha256"

	"github.com/ac0v/aspera/pkg/crypto/curve25519"
)

func secretPhraseToPrivateKey(secretPhrase string) *[32]byte {
	bs := sha256.Sum256([]byte(secretPhrase))
	bs[31] &= 0x7F
	bs[31] |= 0x40
	bs[0] &= 0xF8
	return &bs
}

func secretPhraseToPublicKey(secretPhrase string) []byte {
	pubKey := make([]byte, 32)
	encryptedSecretPhrase := sha256.Sum256([]byte(secretPhrase))
	curve25519.Keygen(pubKey[:], nil, encryptedSecretPhrase[:])
	return pubKey
}
