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

	registry *r.Registry
}

func NewManager(client *Client, registry *r.Registry, scanForNewPeersInterval time.Duration) Manager {
	m := &manager{
		scanForNewPeersInterval: scanForNewPeersInterval,
		allPeers:                make(map[string]*Peer),
		blockedPeers:            make(map[string]*Peer),
		registry:                registry,
	}

	m.initPeers(client, registry.Config.Network.P2P.Peers)

	go m.scanForNewPeers(client)

	return m
}

func (m *manager) initPeers(client *Client, peerBaseUrls []string) {
	for _, url := range peerBaseUrls {
		peer, err := NewPeer(m.registry, url)
		if err != nil {
			continue
		}
		m.addPeersOf(client, peer)
	}
}

func (m *manager) addPeersOf(client *Client, peer *Peer) {
	getPeersMsg := new(pb.GetPeers)
	res, err := client.buildRequest("getPeers").Post(peer.apiUrl)
	if err != nil {
		return
	}

	client.unmarshaler.Unmarshal(bytes.NewReader(res.Body()), getPeersMsg)

	m.allPeersMu.Lock()
	m.allPeers[peer.baseUrl] = peer
	for _, baseUrl := range getPeersMsg.Peers {
		p, err := NewPeer(m.registry, baseUrl)
		if err != nil {
			continue
		}
		m.allPeers[baseUrl] = p
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
	_, blocked := m.blockedPeers[pToBlock.baseUrl]
	if blocked {
		m.blockedPeersMu.Unlock()
		return
	}
	m.blockedPeers[pToBlock.baseUrl] = pToBlock
	m.blockedPeersMu.Unlock()

	m.allPeersMu.Lock()
	delete(m.allPeers, pToBlock.baseUrl)
	m.allPeersMu.Unlock()
}

func (m *manager) scanForNewPeers(client *Client) {
	for range time.NewTicker(m.scanForNewPeersInterval).C {
		m.addPeersOf(client, m.RandomPeer())
	}
}
