package compat

import (
	api "github.com/ac0v/aspera/pkg/api/p2p"
	"github.com/ac0v/aspera/pkg/api/p2p/compat/template"

	"github.com/valyala/fastjson"
)

func Upgrade(srcBs []byte) ([]byte, error) {
	var err error
	var src *fastjson.Value

	if src, err = new(fastjson.Parser).ParseBytes(srcBs); err != nil {
		return []byte{}, err
	}

	return []byte(template.Upgrade(src)), err
}

func Downgrade(pb *api.GetNextBlocksResponse) []byte {
	return []byte(template.Downgrade(pb))
}
