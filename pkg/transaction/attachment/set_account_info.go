package attachment

import (
	"encoding/binary"

	"gopkg.in/restruct.v1"
)

type SetAccountInfo struct {
	NumName        uint8  `struct:"uint8,sizeof=Name" json:"-"`
	Name           string `json:"name"`
	NumDescription uint16 `struct:"uint16,sizeof=Description" json:"-"`
	Description    string `json:"description"`
	Version        int8   `struct:"-" json:"version.AccountInfo,omitempty"`
}

func (attachment *SetAccountInfo) FromBytes(bs []byte, version uint8) (int, error) {
	err := restruct.Unpack(bs, binary.LittleEndian, attachment)
	return 1 + len(attachment.Name) + 1 + len(attachment.Description), err
}

func (attachment *SetAccountInfo) ToBytes(version uint8) ([]byte, error) {
	return restruct.Pack(binary.LittleEndian, attachment)
}

func (attachment *SetAccountInfo) GetFlag() uint32 {
	return StandardAttachmentFlag
}
