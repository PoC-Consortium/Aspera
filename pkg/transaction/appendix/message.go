package appendix

import (
	pb "github.com/PoC-Consortium/Aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/Aspera/pkg/encoding"
)

type Message struct {
	*pb.Message
}

func (a *Message) WriteBytes(e encoding.Encoder) {
	e.WriteBytesWithInt32Len(a.IsText, []byte(a.Content))
}

func (a *Message) SizeInBytes() int {
	return 4 + len(a.Content)
}
