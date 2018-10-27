package compat

import (
	"bytes"
	"testing"

	"fmt"
	//"github.com/json-iterator/go"
	//"encoding/json"
	"github.com/Jeffail/gabs"
	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/jsonpb"

	api "github.com/ac0v/aspera/pkg/api/p2p"
)

func TestParseBlocks(t *testing.T) {
	unmarshaler := &jsonpb.Unmarshaler{AllowUnknownFields: false}
	//var json = jsoniter.ConfigCompatibleWithStandardLibrary
	for _, parseTest := range ParseTests {
		data := []byte(parseTest.JSON)
		data, _ = Upgrade(data)

		msg := new(api.GetNextBlocksResponse)
		if !assert.NoError(t, unmarshaler.Unmarshal(bytes.NewReader(data), msg)) {
			fmt.Println(string(data))
			panic(1)
		}
		//panic(string(data))
		/*
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
		*/
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
