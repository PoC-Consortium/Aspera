package appendix

import (
	pb "github.com/PoC-Consortium/aspera/pkg/api/p2p"
	"github.com/PoC-Consortium/aspera/pkg/encoding"
)

type PublicKeyAnnouncement struct {
	*pb.PublicKeyAnnouncement
}

func (a *EncryptedToSelfMessage) WriteBytes(e encoding.Encoder) {
	e.Writebytes(a.PublicKey)
}

func (a *EncryptedToSelfMessage) SizeInBytes() int {
	return len(a.PublicKey)
}
