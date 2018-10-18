package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SetAlias struct {
	NumAliasName uint8  `struct:"uint8,sizeof=AliasName" json:"-"`
	AliasName    string `json:"alias"`
	NumAliasURI  uint16 `struct:"uint16,sizeof=AliasURI" json:"-"`
	AliasURI     string `json:"uri"`
	Version      int8   `struct:"-" json:"version.AliasAssignment,omitempty"`
}

func (attachment *SetAlias) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 1 + len(attachment.AliasName) + 2 + len(attachment.AliasURI), err
}

func (attachment *SetAlias) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (attachment *SetAlias) GetFlag() uint32 {
	return StandardAttachmentFlag
}
