package encryption

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"github.com/PoC-Consortium/Aspera/pkg/crypto/curve25519"
)

var (
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = errors.New("invalid blocksize")

	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")

	// ErrInvalidPKCS7Padding indicates PKCS7 unpad fails to bad input.
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")

	ErrInvalidCipherData = errors.New("invalid cipherdata")
)

func encrypt(data, encryptorPrivateKey, decryptorPublicKey []byte) ([]byte, []byte, error) {
	if len(data) == 0 {
		return []byte{}, []byte{}, nil
	}

	data, err := pkcs7Pad(data, aes.BlockSize)
	if err != nil {
		return nil, nil, err
	}

	encrypted := make([]byte, len(data)+aes.BlockSize)

	nonce := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, nil, err
	}

	buf := bytes.NewBuffer(nil)
	w := gzip.NewWriter(buf)

	w.Write(data)
	if err := w.Close(); err != nil {
		return nil, nil, err
	}

	sharedSecret := make([]byte, 32)
	curve25519.Curve(sharedSecret, encryptorPrivateKey, decryptorPublicKey)
	for i := range nonce {
		sharedSecret[i] ^= nonce[i]
	}

	key := sha256.Sum256(sharedSecret)

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, nil, err
	}

	stream := cipher.NewCBCEncrypter(block, iv)

	stream.CryptBlocks(encrypted[aes.BlockSize:], data)

	return encrypted, nonce, nil
}

func decrypt(encrypted, nonce, decryptorPrivateKey, encryptorPublicKey []byte) ([]byte, error) {
	if len(encrypted) == 0 {
		return []byte{}, nil
	}

	if len(encrypted)%aes.BlockSize != 0 {
		return nil, ErrInvalidCipherData
	}

	decrypted := make([]byte, len(encrypted)+aes.BlockSize)

	iv := encrypted[:aes.BlockSize]
	sharedSecret := make([]byte, 32)
	curve25519.Curve(sharedSecret, decryptorPrivateKey, encryptorPublicKey)
	for i := range nonce {
		sharedSecret[i] ^= nonce[i]
	}

	key := sha256.Sum256(sharedSecret)

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCBCDecrypter(block, iv)
	stream.CryptBlocks(decrypted[aes.BlockSize:], encrypted)

	decrypted, err = pkcs7Unpad(decrypted, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	return decrypted[32:], nil
}

// pkcs7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// pkcs7Unpad validates and unpads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}
