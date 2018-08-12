package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

const secretPhrase = "glad suffer red during single glow shut slam hill death lust although"

func TestSecretPhraseToPrivateKey(t *testing.T) {
	privateKey := secretPhraseToPrivateKey(secretPhrase)
	assert.Equal(t, "e04c16cfd3d1cbf11a51fa0c75c09d43d307a10ae5149b1ec0ddba661b9d2f5e",
		hex.EncodeToString(privateKey[:]))
}

func TestSecretPhraseToPublicKey(t *testing.T) {
	publicKey := secretPhraseToPublicKey(secretPhrase)
	assert.Equal(t, "a9fc9b42e3918c3e109fdc819aa092fd06ce7c18d653bf003b31d48d35699104",
		hex.EncodeToString(publicKey[:]))
}

func TestSign(t *testing.T) {
	msg := make([]byte, 128)
	for i := 0; i < 128; i++ {
		msg[i] = byte(i)
	}
	sig := sign(msg, secretPhrase)
	assert.Equal(t, "cea74dd7177ddd8bd88e4a0d8807da533e95195d97cc4bcb0f946f355a8c85035ebe57fe9628f93405cf9f7e83bfc64968bc992d3508434a3023fba6b6cab491", (hex.EncodeToString(sig)))
}

func BenchmarkSecretPhraseToPrivateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		secretPhraseToPrivateKey(secretPhrase)
	}
}

func BenchmarkSecretPhraseToPublicKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		secretPhraseToPublicKey(secretPhrase)
	}
}

func BenchmarkSign(b *testing.B) {
	msg := make([]byte, 128)
	for i := 0; i < 128; i++ {
		msg[i] = byte(i)
	}
	for i := 0; i < b.N; i++ {
		sign(msg, secretPhrase)
		for j := 0; j < 128; j++ {
			msg[j] = byte(i + j)
		}
	}
}
