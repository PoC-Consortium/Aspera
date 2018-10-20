package block

import (
	//"bytes"
	"testing"

	//"fmt"
	"github.com/json-iterator/go"
	//"github.com/golang/protobuf/jsonpb"
	//"encoding/json"
	"github.com/Jeffail/gabs"
	"github.com/stretchr/testify/assert"

	"github.com/ac0v/aspera/pkg/transaction"
)

func TestParseBlocks(t *testing.T) {
	//unmarshaler := &jsonpb.Unmarshaler{AllowUnknownFields: false}
	blockMsg := new(Block)

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for _, parseTest := range ParseTests {
		// convert from JSON -> go type
		err := json.Unmarshal([]byte(parseTest.JSON), &blockMsg)
		assert.Nil(t, err, "parse block with id", parseTest.BlockID)
		if err != nil {
			panic(err)
		}

		marshalled, _ := json.Marshal(blockMsg)

		compareJSON(t, string(marshalled), parseTest.JSON)

		for i, tx := range blockMsg.Transactions {
			bs, _ := tx.ToBytes()
			blockMsg.Transactions[i], _ = transaction.FromBytes(bs)
		}

		marshalled, _ = json.Marshal(blockMsg)
		compareJSON(t, string(marshalled), parseTest.JSON)

	}
}

func compareJSON(t *testing.T, left string, right string) {
	comperands := []string{left, right}
	for i, comperand := range comperands {
		jsonParsed, _ := gabs.ParseJSON([]byte(comperand))
		comperands[i] = string(jsonParsed.StringIndent("", "  "))
	}
	assert.EqualValues(t, comperands[1], comperands[0])
}
