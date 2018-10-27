package compat

import (
	"github.com/valyala/fastjson"

	"github.com/ac0v/aspera/pkg/api/p2p/compat/template"
)

func Upgrade(srcBs []byte) ([]byte, error) {
	var err error
	var src *fastjson.Value

	if src, err = new(fastjson.Parser).ParseBytes(srcBs); err != nil {
		return []byte{}, err
	}

	return []byte(template.Upgrade(src)), err
}
