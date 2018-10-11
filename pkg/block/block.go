package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ac0v/aspera/pkg/burstmath"
	t "github.com/ac0v/aspera/pkg/transaction"
)

var (
	ErrBlockUnexpectedLen = errors.New("block unexpected length in byte serialisation")
)

const (
	// TODO: move constants
	oneBurst = 100000000
)

type Block struct {
	PayloadLength       uint32           `json:"payloadLength"`
	TotalAmountNQT      int64            `json:"totalAmountNQT"`
	GenerationSignature string           `json:"generationSignature,omitempty"`
	GeneratorPublicKey  string           `json:"generatorPublicKey,omitempty"`
	PayloadHash         string           `json:"payloadHash,omitempty"`
	BlockSignature      string           `json:"blockSignature,omitempty"`
	Transactions        []*t.Transaction `json:"transactions"`
	Version             int32            `json:"version,omitempty"`
	Nonce               uint64           `json:"nonce,omitempty,string"`
	TotalFeeNQT         int64            `json:"totalFeeNQT,omitempty"`
	BlockATs            *string          `json:"blockATs"`
	PreviousBlock       uint64           `json:"previousBlock,omitempty,string"`
	Timestamp           uint32           `json:"timestamp,omitempty"`
	Block               uint64           `json:"block,omitempty,string"`
	Height              int32            `json:"height,omitempty"`
	PreviousBlockHash   string           `json:"previousBlockHash,omitempty"` // if version > 1
}

func (b *Block) CalcScoop() (uint32, error) {
	if genSig, err := hex.DecodeString(b.GenerationSignature); err == nil {
		return burstmath.CalcScoop(b.Height, genSig), nil
	} else {
		return 0, err
	}
}

func (b *Block) ToBytes() ([]byte, error) {
	bsCap := 4 + 4 + 8 + 4 + 4 + 32 + 32 + (32 + 32) + 8 + 64
	if b.Version < 3 {
		bsCap += 4 + 4
	} else {
		bsCap += 8 + 8
	}
	if b.BlockATs != nil {
		bsCap += len(*b.BlockATs)
	}

	w := bytes.NewBuffer(nil)

	if err := binary.Write(w, binary.LittleEndian, b.Version); err != nil {
		return nil, err
	}

	if err := binary.Write(w, binary.LittleEndian, b.Timestamp); err != nil {
		return nil, err
	}

	if err := binary.Write(w, binary.LittleEndian, b.PreviousBlock); err != nil {
		return nil, err
	}

	if err := binary.Write(w, binary.LittleEndian, uint32(len(b.Transactions))); err != nil {
		return nil, err
	}

	if b.Version < 3 {
		totalAmountQNT := int32(b.TotalAmountNQT / oneBurst)
		if err := binary.Write(w, binary.LittleEndian, totalAmountQNT); err != nil {
			return nil, err
		}

		totalFeeNQT := int32(b.TotalFeeNQT / oneBurst)
		if err := binary.Write(w, binary.LittleEndian, totalFeeNQT); err != nil {
			return nil, err
		}
	} else {
		if err := binary.Write(w, binary.LittleEndian, b.TotalAmountNQT); err != nil {
			return nil, err
		}

		if err := binary.Write(w, binary.LittleEndian, b.TotalFeeNQT); err != nil {
			return nil, err
		}
	}

	if err := binary.Write(w, binary.LittleEndian, b.PayloadLength); err != nil {
		return nil, err
	}

	payloadHash, err := hex.DecodeString(b.PayloadHash)
	if err != nil {
		return nil, err
	}
	if err := binary.Write(w, binary.LittleEndian, payloadHash); err != nil {
		return nil, err
	}

	generatorPublicKey, err := hex.DecodeString(b.GeneratorPublicKey)
	if err != nil {
		return nil, err
	}
	if err := binary.Write(w, binary.LittleEndian, generatorPublicKey); err != nil {
		return nil, err
	}

	generationSignature, err := hex.DecodeString(b.GenerationSignature)
	if err != nil {
		return nil, err
	}
	if err := binary.Write(w, binary.LittleEndian, generationSignature); err != nil {
		return nil, err
	}

	if b.Version > 1 {
		previousBlockHash, err := hex.DecodeString(b.PreviousBlockHash)
		if err != nil {
			return nil, err
		}
		if err := binary.Write(w, binary.LittleEndian, previousBlockHash); err != nil {
			return nil, err
		}
	}

	if err := binary.Write(w, binary.LittleEndian, b.Nonce); err != nil {
		return nil, err
	}

	if b.BlockATs != nil {
		blockATs, err := hex.DecodeString(*b.BlockATs)
		if err != nil {
			return nil, err
		}
		if err := binary.Write(w, binary.LittleEndian, blockATs); err != nil {
			return nil, err
		}
	}

	blockSignature, err := hex.DecodeString(b.BlockSignature)
	if err != nil {
		return nil, err
	}
	if err := binary.Write(w, binary.LittleEndian, blockSignature); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (b *Block) CalculateHash() (*[32]byte, error) {
	if bs, err := b.ToBytes(); err == nil {
		// hash := sha256.New()
		// hash.Write(bs)
		bs := sha256.Sum256(bs)
		// bs := hash.Sum(nil)
		return &bs, nil
	} else {
		return nil, err
	}
}

func (b *Block) CalculateID() (uint64, error) {
	if hash, err := b.CalculateHash(); err == nil {
		bigInt := big.NewInt(0)
		bigInt.SetBytes([]byte{
			hash[7], hash[6], hash[5], hash[4],
			hash[3], hash[2], hash[1], hash[0]})
		return bigInt.Uint64(), nil
	} else {
		return 0, err
	}
}
