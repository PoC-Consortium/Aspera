package compat

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/valyala/fastjson"
)

func Upgrade(srcBs []byte) ([]byte, error) {
	var err error
	var src *fastjson.Value

	if src, err = new(fastjson.Parser).ParseBytes(srcBs); err != nil {
		return []byte{}, err
	}
	res := string(srcBs)

	fromTo := map[string]string{
		"totalAmountNQT": "totalAmount",
		"totalFeeNQT":    "totalFee",
		"amountNQT":      "amount",
		"feeNQT":         "fee",
		"quantityQNT":    "quantity",
	}
	for from, to := range fromTo {
		res = strings.Replace(res, fmt.Sprintf(`"%s":`, from), fmt.Sprintf(`"%s":`, to), -1)
	}

	if nextBlocks := src.GetArray("nextBlocks"); nextBlocks != nil {
		for _, block := range nextBlocks {
			if transactions := block.GetArray("transactions"); transactions != nil {
				for _, transaction := range transactions {
					transaction.GetObject().Visit(func(k []byte, v *fastjson.Value) {
						if bytes.Equal(k, []byte("attachment")) {
						} else {
							//x, _ := v.StringBytes()
							//fmt.Println(string(x))
						}
					})
				}
			}
		}
	}
	return []byte(res), err
}
