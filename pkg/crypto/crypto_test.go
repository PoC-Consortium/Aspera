package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretPhraseToPrivateKey(t *testing.T) {
	privateKey := secretPhraseToPrivateKey("glad suffer red during single glow shut slam hill death lust although")
	assert.Equal(t, "e04c16cfd3d1cbf11a51fa0c75c09d43d307a10ae5149b1ec0ddba661b9d2f5e",
		hex.EncodeToString(privateKey[:]))
}
