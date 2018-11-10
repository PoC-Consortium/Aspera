package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var p Peer
var m *manager

func TestMain(tm *testing.M) {
	p = NewPeer("http://wallet.burst.cryptoguru.org:8123/burst")
	m = NewManager(
		[]string{"wallet.burst.cryptoguru.org", "burst-a.bayernport.com"}, []string{"v4"}).(*manager)
	tm.Run()
}

func TestSetAndGetHeight(t *testing.T) {
	assert.Nil(t, p.SetHeight())
	assert.NotEmpty(t, p.GetHeight())
}

func TestGetNextBlockIDsBody(t *testing.T) {
	body, err := p.GetNextBlockIDsBody(17169998969130562818)
	if assert.Nil(t, err) {
		assert.NotEmpty(t, body)
	}
}

func TestGetNextBlocksBody(t *testing.T) {
	body, err := p.GetNextBlocksBody(17169998969130562818)
	if assert.Nil(t, err) {
		assert.NotEmpty(t, body)
	}
}

func TestGetPeerUrls(t *testing.T) {
	peers, err := p.GetPeerUrls()
	if assert.Nil(t, err) {
		assert.NotEmpty(t, peers)
	}
}

func TestNewManager(t *testing.T) {
	assert.True(t, len(m.peers) > 2, "added no extra peers")
}
