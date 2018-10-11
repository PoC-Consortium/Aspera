package crypto

import (
	"crypto/sha256"

	"github.com/ac0v/aspera/pkg/crypto/curve25519"
)

func secretPhraseToPrivateKey(secretPhrase string) []byte {
	bs := sha256.Sum256([]byte(secretPhrase))
	bs[31] &= 0x7F
	bs[31] |= 0x40
	bs[0] &= 0xF8
	return bs[:]
}

func secretPhraseToPublicKey(secretPhrase string) []byte {
	pubKey := make([]byte, 32)
	encryptedSecretPhrase := sha256.Sum256([]byte(secretPhrase))
	curve25519.Keygen(pubKey[:], nil, encryptedSecretPhrase[:])
	return pubKey
}

func sign(msg []byte, secretPhrase string) []byte {
	P := make([]byte, 32)
	s := make([]byte, 32)

	digest := sha256.New()
	_, _ = digest.Write([]byte(secretPhrase))
	encryptedSecretPhrase := digest.Sum(nil)[:]
	curve25519.Keygen(P, s, encryptedSecretPhrase)

	digest.Reset()
	_, _ = digest.Write(msg)
	m := digest.Sum(nil)

	digest.Reset()
	_, _ = digest.Write(m)
	_, _ = digest.Write(s)

	x := digest.Sum(nil)

	Y := make([]byte, 32)

	curve25519.Keygen(Y, nil, x)

	digest.Reset()
	_, _ = digest.Write(m)
	_, _ = digest.Write(Y)
	h := digest.Sum(nil)

	v := make([]byte, 32)
	curve25519.Sign(v, h, x, s)

	sig := make([]byte, 64)
	copy(sig[:32], v)
	copy(sig[32:], h)

	return sig
}

func Verify(sig, msg, pubKey []byte, canonical bool) bool {
	if canonical {
		if !curve25519.IsCanonicalSignature(sig) {
			return false
		}
		if !curve25519.IsCanonicalPublicKey(pubKey) {
			return false
		}
	}

	Y := make([]byte, 32)
	v := make([]byte, 32)
	copy(v, sig[:32])

	h := make([]byte, 32)
	copy(h, sig[32:])

	curve25519.Verify(Y, v, h, pubKey)

	digest := sha256.New()
	_, _ = digest.Write(msg)
	m := digest.Sum(nil)

	digest.Reset()
	_, _ = digest.Write(m)
	_, _ = digest.Write(Y)
	h2 := digest.Sum(nil)

	for i := range h {
		if h[i] != h2[i] {
			return false
		}
	}

	return true
}
