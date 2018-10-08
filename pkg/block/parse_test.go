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

		foo, _ := json.Marshal(blockMsg)
		/*
			for _, tx := range blockMsg.Transactions {
				if tx.Header.RecipientID == 0 {
					panic(string(foo))
				}
			}*/
		var comperands [2]string
		comperands[0] = string(foo)
		comperands[1] = parseTest.JSON
		for i, comperand := range comperands {
			jsonParsed, _ := gabs.ParseJSON([]byte(comperand))
			comperands[i] = string(jsonParsed.StringIndent("", "  "))
		}
		assert.EqualValues(t, comperands[1], comperands[0])

		//result, _ := json.Marshal(&blockMsg)
		//assert.JSONEq(t, parseTest.JSON, string(foo), "json valid")

		// fmt.Println(parseTest.JSON)
		//fmt.Printf("%+v\n", blockMsg.Transactions)
		//		panic("fuck")
	}
}
