package appendix

import (
	pb "github.com/PoC-Consortium/aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/aspera/pkg/encoding"
)

type EncrypedMessage struct {
	*pb.Message
}

func (a *EncryptedMessgae) WriteBytes(e encoding.Encoder) {
	e.WriteBytesWithInt32Len(a.IsText, []byte(a.Data))
	e.WriteBytesWithInt32Len(a.IsText, []byte(a.Nonce))
}

func (a *EncryptedMessage) SizeInBytes() int {
	return 4 + len(a.Data) + len(a.Nonce)
}
