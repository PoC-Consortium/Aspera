package block

import (
	"encoding/hex"
	"github.com/ac0v/aspera/pkg/burstmath"
	t "github.com/ac0v/aspera/pkg/transaction"
)

type Block struct {
	PayloadLength       int64            `json:"payloadLength,omitempty"`
	TotalAmountNQT      int64            `json:"totalAmountNQT"`
	GenerationSignature string           `json:"generationSignature,omitempty"`
	GeneratorPublicKey  string           `json:"generatorPublicKey,omitempty"`
	PayloadHash         string           `json:"payloadHash,omitempty"`
	BlockSignature      string           `json:"blockSignature,omitempty"`
	Transactions        []*t.Transaction `json:"transactions,omitempty"`
	Version             int32            `json:"version,omitempty"`
	Nonce               string           `json:"nonce,omitempty"`
	TotalFeeNQT         int64            `json:"totalFeeNQT,omitempty"`
	BlockATs            *string          `json:"blockATs"`
	PreviousBlock       uint64           `json:"previousBlock,omitempty,string"`
	Timestamp           int64            `json:"timestamp,omitempty"`
	Block               uint64           `json:"block,omitempty,string"`
	Height              int32            `json:"height,omitempty"`
	PreviousBlockHash   string           `json:"previousBlockHash,omitempty"` // if version > 1
}

func (b *Block) CalcScoop() uint32 {
	genSig, err := hex.DecodeString(b.GenerationSignature)
	if err != nil {
		panic("generation signature wrong format")
	}
	return burstmath.CalcScoop(b.Height, genSig)
}
