package p2p

import (
	"sync"
	"time"

	"github.com/lukechampine/randmap/safe"
)

type Manager interface {
	BlockPeer(p *Peer, reason int)
	RandomPeer() *Peer
}

type manager struct {
	allPeers   map[string]*Peer
	allPeersMu sync.RWMutex

	blockedPeers   map[string]*Peer
	blockedPeersMu sync.Mutex

	scanForNewPeersInterval time.Duration

	blacklisted time.Time

	internetProtocols []string
}

func NewManager(c *client, peers, internetProtocols []string, scanForNewPeersInterval time.Duration) Manager {
	m := &manager{
		scanForNewPeersInterval: scanForNewPeersInterval,
		allPeers:                make(map[string]*Peer),
		blockedPeers:            make(map[string]*Peer),
		internetProtocols:       internetProtocols,
	}

	m.initPeers(c, peers)

	go m.jobs(c)

	return m
}

func (m *manager) jobs(c *client) {
	unblockTicker := time.NewTicker(5 * time.Minute)
	scanForPeersTicker := time.NewTicker(m.scanForNewPeersInterval)
	for {
		select {
		case <-unblockTicker.C:
			m.unblockPeers()
		case <-scanForPeersTicker.C:
			m.addPeersOf(c, m.RandomPeer())
		}
	}
}

func (m *manager) unblockPeers() {
	var unblockedPeers []*Peer

	m.blockedPeersMu.Lock()
	for _, p := range m.blockedPeers {
		unblocked := p.TryUnblock()
		if unblocked {
			unblockedPeers = append(unblockedPeers, p)
		}
	}
	m.blockedPeersMu.Unlock()

	if len(unblockedPeers) == 0 {
		return
	}

	m.allPeersMu.Lock()
	for _, p := range unblockedPeers {
		m.allPeers[p.baseUrl] = p
	}
	m.allPeersMu.Unlock()
}

func (m *manager) initPeers(c *client, peerBaseUrls []string) {
	for _, url := range peerBaseUrls {
		peer, err := NewPeer(url, m.internetProtocols)
		if err != nil {
			continue
		}
		m.addPeersOf(c, peer)
	}
}

func (m *manager) addPeersOf(c Client, peer *Peer) {
	getPeersMsg, err := c.GetPeersOf(peer.apiUrl)
	if err != nil {
		return
	}

	m.allPeersMu.Lock()
	m.allPeers[peer.baseUrl] = peer
	for _, baseUrl := range getPeersMsg.Peers {
		p, err := NewPeer(baseUrl, m.internetProtocols)
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

func (m *manager) BlockPeer(pToBlock *Peer, reason int) {
	pToBlock.Block(reason)

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
