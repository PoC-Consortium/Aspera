package p2p

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
)

func TestParseBlocks(t *testing.T) {
	unmarshaler := &jsonpb.Unmarshaler{AllowUnknownFields: true}
	blockMsg := new(Block)

	for _, parseTest := range ParseTests {
		err := unmarshaler.Unmarshal(bytes.NewReader([]byte(parseTest.JSON)), blockMsg)
		assert.Nil(t, err, "parse block with id", parseTest.BlockID)
	}
}
