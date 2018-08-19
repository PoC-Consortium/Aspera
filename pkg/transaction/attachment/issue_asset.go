package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type IssueAssetTransaction struct {
	NumName        uint8 `struct:"uint8,sizeof=Name"`
	Name           []byte
	NumDescription uint16 `struct:"uint16,sizeof=Description"`
	Description    []byte
	Quantity       uint64
	Decimals       uint8
}

func IssueAssetTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx IssueAssetTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 1 + len(tx.Name) + 2 + len(tx.Description) + 8 + 1, err
}
