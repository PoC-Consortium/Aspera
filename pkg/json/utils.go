package json

import (
	"encoding/hex"
	"errors"

	"github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	ErrWrongQuotingFormat = errors.New("wrong quoting format")
)

type HexSlice []byte

func (hexSlice *HexSlice) MarshalJSON() ([]byte, error) {
	encoded := make([]byte, hex.EncodedLen(len(*hexSlice)))
	_ = hex.Encode(encoded, *hexSlice)
	QuoteBytes(&encoded)
	return encoded, nil
}

func (hexSlice *HexSlice) UnmarshalJSON(bs []byte) error {
	if err := UnquoteBytes(&bs); err != nil {
		return err
	}

	decoded := make([]byte, hex.DecodedLen(len(bs)))
	if _, err := hex.Decode(decoded, bs); err == nil {
		*hexSlice = decoded
		return nil
	} else {
		return err
	}
}

func UnquoteBytes(bs *[]byte) error {
	tmp := *bs
	if len(tmp) < 2 || tmp[0] != '"' || tmp[len(tmp)-1] != '"' {
		return ErrWrongQuotingFormat
	}
	tmp = tmp[1 : len(tmp)-1]
	*bs = tmp
	return nil
}

func QuoteBytes(bs *[]byte) {
	tmp := make([]byte, len(*bs)+2)
	copy(tmp[1:len(tmp)-1], *bs)
	tmp[0] = '"'
	tmp[len(tmp)-1] = '"'
	*bs = tmp
}
