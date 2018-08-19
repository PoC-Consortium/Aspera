package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SetAliasTransaction struct {
	NumAliasName uint8 `struct:"uint8,sizeof=AliasName"`
	AliasName    []byte
	NumAliasURI  uint16 `struct:"uint16,sizeof=AliasURI"`
	AliasURI     []byte
}

func SetAliasTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx SetAliasTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 1 + len(tx.AliasName) + 2 + len(tx.AliasURI), err
}
