package p2p

import (
	"bytes"
	"sync"
	"time"

	pb "github.com/ac0v/aspera/internal/api/protobuf-spec"
	r "github.com/ac0v/aspera/pkg/registry"
	"github.com/lukechampine/randmap/safe"
)

var (
	syncPeerCount = 10
)

type Manager interface {
	BlockPeer(p *Peer)
	RandomPeer() *Peer
}

type manager struct {
	allPeers   map[string]*Peer
	allPeersMu sync.RWMutex

	blockedPeers   map[string]*Peer
	blockedPeersMu sync.Mutex

	scanForNewPeersInterval time.Duration

	blacklisted time.Time
}

func NewManager(client *Client, registry *r.Registry, scanForNewPeersInterval time.Duration) Manager {
	m := &manager{
		scanForNewPeersInterval: scanForNewPeersInterval,
		allPeers:                make(map[string]*Peer),
		blockedPeers:            make(map[string]*Peer),
	}

	m.initPeers(client, registry.Config.Peers)

	go m.scanForNewPeers(client)

	return m
}

func (m *manager) initPeers(client *Client, peerBaseURLs []string) {
	for _, url := range peerBaseURLs {
		peer, err := NewPeer(url)
		if err != nil {
			continue
		}
		m.addPeersOf(client, peer)
	}
}

func (m *manager) addPeersOf(client *Client, peer *Peer) {
	getPeersMsg := new(pb.GetPeers)
	res, err := client.buildRequest("getPeers").Post(peer.apiURL)
	if err != nil {
		return
	}

	client.unmarshaler.Unmarshal(bytes.NewReader(res.Body()), getPeersMsg)

	m.allPeersMu.Lock()
	m.allPeers[peer.baseURL] = peer
	for _, baseURL := range getPeersMsg.Peers {
		p, err := NewPeer(baseURL)
		if err != nil {
			continue
		}
		m.allPeers[baseURL] = p
	}
	m.allPeersMu.Unlock()
}

func (m *manager) RandomPeer() *Peer {
	m.allPeersMu.RLock()
	p := randmap.Val(m.allPeers).(*Peer)
	m.allPeersMu.RUnlock()
	return p
}

func (m *manager) BlockPeer(pToBlock *Peer) {
	m.blockedPeersMu.Lock()
	_, blocked := m.blockedPeers[pToBlock.baseURL]
	if blocked {
		m.blockedPeersMu.Unlock()
		return
	}
	m.blockedPeers[pToBlock.baseURL] = pToBlock
	m.blockedPeersMu.Unlock()

	m.allPeersMu.Lock()
	delete(m.allPeers, pToBlock.baseURL)
	m.allPeersMu.Unlock()
}

func (m *manager) scanForNewPeers(client *Client) {
	for range time.NewTicker(m.scanForNewPeersInterval).C {
		m.addPeersOf(client, m.RandomPeer())
	}
}
