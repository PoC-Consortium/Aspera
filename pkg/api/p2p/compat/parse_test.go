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
		if assert.NoError(t, unmarshaler.Unmarshal(bytes.NewReader(data), msg)) {
			unmarshaler.Unmarshal(bytes.NewReader(data), msg)
			dst := Downgrade(msg)
			//fmt.Println(string(dst))
			compareJSON(t, string(dst), parseTest.JSON)
			//panic(string(dst))
			//fmt.Println(string(dst))
		} else {
			fmt.Println(string(data))
		}
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
