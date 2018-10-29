package compat

import (
	"bytes"
	"io/ioutil"
	"testing"

	_ "fmt"
	"github.com/Jeffail/gabs"
	"github.com/stretchr/testify/assert"

	"github.com/golang/protobuf/jsonpb"

	api "github.com/ac0v/aspera/pkg/api/p2p"
)

func TestParseBlocks(t *testing.T) {
	files, err := ioutil.ReadDir("test_files")
	if err != nil {
		t.Fatal("failed to open test files")
	}

	var failed int
	unmarshaler := &jsonpb.Unmarshaler{AllowUnknownFields: false}
	for _, f := range files {
		t.Log("parsing ", f.Name())
		javaWalletBs, err := ioutil.ReadFile("test_files/" + f.Name())
		if err != nil {
			t.Fatalf("failed to open file %s", "test_files/"+f.Name())
		}

		protoJSONBs, _ := Upgrade(javaWalletBs)
		msg := new(api.GetNextBlocksResponse)
		if assert.NoError(t, unmarshaler.Unmarshal(bytes.NewReader(protoJSONBs), msg)) {
			javaWalletBsRebuilt := Downgrade(msg)
			//fmt.Println(string(javaWalletBsRebuilt))
			if !compareJSON(t, string(javaWalletBsRebuilt), string(javaWalletBs)) {
				comperands := []string{string(javaWalletBsRebuilt), string(javaWalletBs)}
				for i, comperand := range comperands {
					jsonParsed, _ := gabs.ParseJSON([]byte(comperand))
					comperands[i] = string(jsonParsed.StringIndent("", "  "))
				}

				//ioutil.WriteFile("/tmp/a."+f.Name(), []byte(comperands[0]), 0644)
				//ioutil.WriteFile("/tmp/b."+f.Name(), []byte(comperands[1]), 0644)
				failed++
			}
			//panic(string(dst))
			// fmt.Println(string(dst))
		} else {
			t.Log(string(protoJSONBs))
		}
	}
}

func compareJSON(t *testing.T, left string, right string) bool {
	comperands := []string{left, right}
	for i, comperand := range comperands {
		jsonParsed, _ := gabs.ParseJSON([]byte(comperand))
		comperands[i] = string(jsonParsed.StringIndent("", "  "))
	}
	return assert.EqualValues(t, comperands[1], comperands[0])
}
