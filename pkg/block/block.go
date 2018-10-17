package block

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"

	"github.com/ac0v/aspera/pkg/burstmath"
	"github.com/ac0v/aspera/pkg/crypto"
	"github.com/ac0v/aspera/pkg/crypto/shabal256"
	jutils "github.com/ac0v/aspera/pkg/json"
	t "github.com/ac0v/aspera/pkg/transaction"

	"github.com/json-iterator/go"
)

var (
	ErrBlockUnexpectedLen          = errors.New("block unexpected length in byte serialisation")
	ErrPreviousBlockMismatch       = errors.New("previous block id doesn't match current block's")
	ErrTimestampTooLate            = errors.New("timestamp to late")
	ErrTimestampSmallerPrevious    = errors.New("timestamp smaller than previous block's")
	ErrGenerationSignatureMismatch = errors.New("generation signature mismatch")
)

const (
	generationSignatureLen = 64
	// TODO: move constants
	oneBurst               = 100000000
	maxTimestampDifference = 15
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Block struct {
	PayloadLength       uint32           `json:"payloadLength"`
	TotalAmountNQT      int64            `json:"totalAmountNQT"`
	GenerationSignature jutils.HexSlice  `json:"generationSignature,omitempty"`
	GeneratorPublicKey  jutils.HexSlice  `json:"generatorPublicKey,omitempty"`
	PayloadHash         jutils.HexSlice  `json:"payloadHash,omitempty"`
	BlockSignature      jutils.HexSlice  `json:"blockSignature,omitempty"`
	Transactions        []*t.Transaction `json:"transactions"`
	Version             int32            `json:"version,omitempty"`
	Nonce               uint64           `json:"nonce,omitempty,string"`
	TotalFeeNQT         int64            `json:"totalFeeNQT,omitempty"`
	BlockATs            *jutils.HexSlice `json:"blockATs"`
	PreviousBlock       uint64           `json:"previousBlock,omitempty,string"`
	Timestamp           uint32           `json:"timestamp,omitempty"`
	Block               uint64           `json:"block,omitempty,string"`
	Height              int32            `json:"height,omitempty"`
	PreviousBlockHash   jutils.HexSlice  `json:"previousBlockHash,omitempty"` // if version > 1
}

func (b *Block) CalcScoop() uint32 {
	return burstmath.CalcScoop(b.Height, b.GenerationSignature)
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

	if err := binary.Write(w, binary.LittleEndian, b.PayloadHash); err != nil {
		return nil, err
	}

	if err := binary.Write(w, binary.LittleEndian, b.GeneratorPublicKey); err != nil {
		return nil, err
	}

	if err := binary.Write(w, binary.LittleEndian, b.GenerationSignature); err != nil {
		return nil, err
	}

	if b.Version > 1 {
		if err := binary.Write(w, binary.LittleEndian, b.PreviousBlockHash); err != nil {
			return nil, err
		}
	}

	if err := binary.Write(w, binary.LittleEndian, b.Nonce); err != nil {
		return nil, err
	}

	if b.BlockATs != nil {
		if err := binary.Write(w, binary.LittleEndian, *b.BlockATs); err != nil {
			return nil, err
		}
	}

	if err := binary.Write(w, binary.LittleEndian, b.BlockSignature); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func (b *Block) CalculateHashAndID() ([32]byte, uint64, error) {
	if bs, err := b.ToBytes(); err == nil {
		hash, id := crypto.BytesToHashAndID(bs)
		return hash, id, nil
	} else {
		return [32]byte{}, 0, err
	}
}

func (b *Block) toError(message string) error {
	if v, err := json.Marshal(b); err == nil {
		return errors.New(message + " << " + string(v))
	} else {
		return nil
	}
}

func (b *Block) Validate(previous *Block) error {
	now := burstmath.DateToTimestamp(time.Now())

	if b.Version != 3 {
		return b.toError("invalid block version")
	} else if b.Timestamp <= previous.Timestamp {
		return ErrTimestampSmallerPrevious
	} else if b.Timestamp > now+maxTimestampDifference {
		return ErrTimestampTooLate
	} else if previousHash, previousID, err := previous.CalculateHashAndID(); err != nil {
		return err
	} else if previousID != b.PreviousBlock {
		return ErrPreviousBlockMismatch
	} else if !bytes.Equal(previousHash[:], b.PreviousBlockHash) {
		return ErrPreviousBlockMismatch
	}

	// ToDo: check for duplicte blocks - may this should go to the raw storage stuff
	// throw new BlockNotAcceptedException("Duplicate block or invalid id for block " + block.getHeight());

	for _, t := range b.Transactions {
		if err := t.VerifySignature(); err != nil {
			return err
		}
	}

	generationSignatureExp := CalculateGenerationSignature(previous)
	for i := range b.GenerationSignature {
		if generationSignatureExp[i] != b.GenerationSignature[i] {
			return ErrGenerationSignatureMismatch
		}
	}

	return nil
}

func CalculateGenerationSignature(previous *Block) []byte {
	bs := make([]byte, 8)
	_, id := crypto.BytesToHashAndID(previous.GeneratorPublicKey)
	binary.BigEndian.PutUint64(bs, id)
	hash := shabal256.Sum256(append(previous.GenerationSignature, bs...))
	return hash[:]
}
