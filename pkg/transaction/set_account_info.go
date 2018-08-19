package transaction

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SetAccountInfoTransaction struct {
	NumName        uint8 `struct:"uint8,sizeof=Name"`
	Name           []byte
	NumDescription uint8 `struct:"uint8,sizeof=Description"`
	Description    []byte
}

func SetAccountInfoTransactionFromBytes(bs []byte) (Attachment, int, error) {
	var tx SetAccountInfoTransaction
	err := restruct.Unpack(bs, binary.LittleEndian, &tx)
	return &tx, 1 + len(tx.Name) + 1 + len(tx.Description), err
}
