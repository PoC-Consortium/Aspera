package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/json-iterator/go"
	"gopkg.in/resty.v1"

	b "github.com/ac0v/aspera/pkg/block"
	"github.com/ac0v/aspera/pkg/config"
	p2p "github.com/ac0v/aspera/pkg/p2p"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var fileHeader = `package block

import "github.com/ac0v/aspera/pkg/json"

type BlockTest struct {
    Block    Block
    BlockATs string
    TXLen    int
}

var BlockTests = []BlockTest{
`

var fileFooter = "}"

var blockTestTmpl = `BlockTest{
    Block: Block{
	PayloadLength:        {{.Block.PayloadLength}},
	TotalAmountNQT:       {{.Block.TotalAmountNQT}},
	GenerationSignature:  []byte{ {{.GenerationSignature}} },
	GeneratorPublicKey:   []byte{ {{.GeneratorPublicKey}} },
	PayloadHash:          []byte{ {{.PayloadHash}} },
	BlockSignature:       []byte{ {{.BlockSignature}} },
	Version:              {{.Block.Version}},
	Nonce:                {{.Block.Nonce}},
	TotalFeeNQT:          {{.Block.TotalFeeNQT}},
	PreviousBlock:        {{.Block.PreviousBlock}},
	Timestamp:            {{.Block.Timestamp}},
	Block:                {{.Block.Block}},
	Height:               {{.Block.Height}},
	PreviousBlockHash:    []byte{ {{.PreviousBlockHash}} },
        BlockATs:             []byte{ {{.BlockATs}} },
    },
    TXLen:     {{.TXLen}},
},
`

type BlockTest struct {
	Block *b.Block

	TXLen int

	BlockATs            string
	GenerationSignature string
	GeneratorPublicKey  string
	PayloadHash         string
	BlockSignature      string
	PreviousBlockHash   string
}

type BlockMsg8125 struct {
	Height int32  `json:"height"`
	Block  uint64 `json:"block,string"`
}

func main() {
	p2pConfig := &config.P2P{
		Timeout: 5 * time.Second,
		Debug:   false,
		Peers:   []string{"wallet.burst.cryptoguru.org"},
	}
	client := p2p.NewClient(p2pConfig, []string{"v4"})

	blockIDs := []uint64{
		16917752638128180357,
		17169998969130562818,
		10851012679396814781,
		8868708821622932189,
		9278508053345228779,
		4947518625215221655,
		13789368535761104494,
		7396048386025791037,
		8038733917809622647,
		16219776992541504875,
		3013573467081493371,
		6218551705245429261,
	}

	var allBlocks []*b.Block
	for _, id := range blockIDs {
	AGAIN:
		if getNextBlocksMsg, _, err := client.GetNextBlocks(id); err == nil {
			allBlocks = append(allBlocks, getNextBlocksMsg.NextBlocks...)
		} else {
			goto AGAIN
		}
	}

	for _, b := range allBlocks {
	RETRY:
		res, err := resty.R().Post("https://wallet.burst.cryptoguru.org:8125/burst?requestType=getBlock&timestamp=" + strconv.FormatUint(uint64(b.Timestamp), 10))
		if err != nil {
			log.Println(err)
			goto RETRY
		}

		var blockMsg8125 BlockMsg8125
		if err := json.Unmarshal(res.Body(), &blockMsg8125); err != nil {
			log.Println(err)
			goto RETRY
		}

		b.Block = blockMsg8125.Block
		b.Height = blockMsg8125.Height
	}

	t, err := template.New("parse test template").Parse(blockTestTmpl)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("pkg/block/block_test_structs.go")
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer([]byte(fileHeader))
	for _, b := range allBlocks {
		blockTest := BlockTest{Block: b}
		blockTest.TXLen = len(b.Transactions)

		if blockTest.Block.BlockATs != nil {
			blockTest.BlockATs = toByteStr(blockTest.Block.BlockATs.Bs)
		}

		blockTest.GenerationSignature = toByteStr(blockTest.Block.GenerationSignature.Bs)
		blockTest.GeneratorPublicKey = toByteStr(blockTest.Block.GeneratorPublicKey.Bs)
		blockTest.PayloadHash = toByteStr(blockTest.Block.PayloadHash.Bs)
		blockTest.BlockSignature = toByteStr(blockTest.Block.BlockSignature.Bs)
		blockTest.PreviousBlockHash = toByteStr(blockTest.Block.PreviousBlockHash.Bs)

		t.Execute(buf, &blockTest)
	}

	buf.Write(([]byte(fileFooter)))
	log.Println(string(buf.Bytes()))

	f.Write(buf.Bytes())
}

func toByteStr(bs []byte) string {
	s := ""
	for i, b := range bs {
		if i == len(bs)-1 {
			s += fmt.Sprintf("%d", b)
		} else {
			s += fmt.Sprintf("%d, ", b)
		}
	}
	return s
}
