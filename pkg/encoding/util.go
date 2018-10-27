package encoding

import (
	"encoding/base64"
	"encoding/hex"
)

func HexStringBytesToBase64Bytes(bs []byte) []byte {
	hexBs := make([]byte, hex.DecodedLen(len(bs)))
	hex.Decode(hexBs, bs)

	base64Bs := make([]byte, base64.StdEncoding.EncodedLen(len(hexBs)))
	base64.StdEncoding.Encode(base64Bs, hexBs)

	return base64Bs
}
