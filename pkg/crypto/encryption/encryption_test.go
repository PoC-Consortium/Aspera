package encryption

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

type encryptionTest struct {
	encryptorPublicKey  string
	encryptorPrivateKey string

	decryptorPublicKey  string
	decryptorPrivateKey string

	msg string
}

var encryptionTests = []encryptionTest{
	encryptionTest{
		encryptorPublicKey:  "a9fc9b42e3918c3e109fdc819aa092fd06ce7c18d653bf003b31d48d35699104",
		encryptorPrivateKey: "e04c16cfd3d1cbf11a51fa0c75c09d43d307a10ae5149b1ec0ddba661b9d2f5e",

		decryptorPublicKey:  "000782b931373ece1d45be84605783581b1d4b71203338aaf38c35ea6cb09a7a",
		decryptorPrivateKey: "5894745467270a9ecf103b882961a419b76215dda9633b1644680a251bf05456",

		msg: "per aspera ad astra",
	},
	encryptionTest{
		encryptorPublicKey:  "a9fc9b42e3918c3e109fdc819aa092fd06ce7c18d653bf003b31d48d35699104",
		encryptorPrivateKey: "e04c16cfd3d1cbf11a51fa0c75c09d43d307a10ae5149b1ec0ddba661b9d2f5e",

		decryptorPublicKey:  "000782b931373ece1d45be84605783581b1d4b71203338aaf38c35ea6cb09a7a",
		decryptorPrivateKey: "5894745467270a9ecf103b882961a419b76215dda9633b1644680a251bf05456",

		msg: "",
	},
}

func TestEncryption(t *testing.T) {
	for _, test := range encryptionTests {
		encryptorPublicKey, _ := hex.DecodeString(test.encryptorPublicKey)
		encryptorPrivateKey, _ := hex.DecodeString(test.encryptorPrivateKey)

		decryptorPublicKey, _ := hex.DecodeString(test.decryptorPublicKey)
		decryptorPrivateKey, _ := hex.DecodeString(test.decryptorPrivateKey)

		msg := []byte(test.msg)

		encrypted, nonce, err := encrypt(msg, encryptorPrivateKey, decryptorPublicKey)
		if assert.Nil(t, err) {
			decrypted, err := decrypt(encrypted, nonce, decryptorPrivateKey, encryptorPublicKey)
			if assert.Nil(t, err) {
				assert.Equal(t, msg, decrypted)
			}
		}
	}
}

func BenchmarkEncryption(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test := encryptionTests[i%len(encryptionTests)]

		encryptorPublicKey, _ := hex.DecodeString(test.encryptorPublicKey)
		encryptorPrivateKey, _ := hex.DecodeString(test.encryptorPrivateKey)

		decryptorPublicKey, _ := hex.DecodeString(test.decryptorPublicKey)
		decryptorPrivateKey, _ := hex.DecodeString(test.decryptorPrivateKey)

		msg := []byte(test.msg)

		encrypted, nonce, _ := encrypt(msg, encryptorPrivateKey, decryptorPublicKey)
		_, _ = decrypt(encrypted, nonce, decryptorPrivateKey, encryptorPublicKey)
	}
}
